package enrichednow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// Client calls the EnrichedNow API. It is safe for use by several goroutines at once.
type Client struct {
	// Scrapers is the catalogue, and where runs are started.
	Scrapers *ScraperService
	// Runs reads, waits for and cancels runs, and reads their records.
	Runs *RunService
	// Usage reports the units used.
	Usage *UsageService
	// Data reads the datasets: companies and their job postings.
	Data *DataService
	// Webhooks manages the endpoints told when a run ends, their signing secrets, and the delivery log.
	Webhooks *WebhookService

	standing atomic.Pointer[Standing]

	apiKey     string
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	maxRetries int
	userAgent  string

	// Replaced in tests.
	sleep  func(ctx context.Context, d time.Duration) error
	jitter func() float64
}

// NewClient builds a client. The API key comes from [WithAPIKey] or, without it, from the environment variable
// ENRICHEDNOW_API_KEY; when neither is there the error is [ErrMissingAPIKey], unless [WithoutAPIKey] was given.
func NewClient(opts ...Option) (*Client, error) {
	cfg := config{baseURL: DefaultBaseURL, timeout: DefaultTimeout, maxRetries: DefaultMaxRetries}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.apiKey == "" {
		cfg.apiKey = strings.TrimSpace(os.Getenv(EnvAPIKey))
	}
	if cfg.apiKey == "" && !cfg.anonymous {
		return nil, ErrMissingAPIKey
	}
	base, err := url.Parse(cfg.baseURL)
	if err != nil || base.Host == "" || (base.Scheme != "http" && base.Scheme != "https") {
		return nil, fmt.Errorf("enrichednow: the base URL %q is not an http(s) address", cfg.baseURL)
	}
	if cfg.httpClient == nil {
		cfg.httpClient = &http.Client{}
	}
	ua := "enrichednow-go/" + Version
	if cfg.userAgent != "" {
		ua = cfg.userAgent + " " + ua
	}
	c := &Client{
		apiKey:     cfg.apiKey,
		baseURL:    strings.TrimRight(cfg.baseURL, "/"),
		httpClient: cfg.httpClient,
		timeout:    cfg.timeout,
		maxRetries: cfg.maxRetries,
		userAgent:  ua,
		sleep:      sleepContext,
		jitter:     rand.Float64,
	}
	c.Scrapers = &ScraperService{client: c}
	c.Runs = &RunService{client: c}
	c.Usage = &UsageService{client: c}
	c.Data = &DataService{
		Companies: &CompanyService{client: c},
		Postings:  &PostingService{client: c},
	}
	c.Webhooks = &WebhookService{client: c, Deliveries: &WebhookDeliveryService{client: c}}
	return c, nil
}

// Standing is where the account stood at the last answer that said so: every answer about runs carries the
// account's units and runs in progress in its headers. Nil until such an answer has come; [UsageService.Standing]
// asks for it.
func (c *Client) Standing() *Standing {
	return c.standing.Load()
}

// Health asks whether the service is up. It needs no API key.
func (c *Client) Health(ctx context.Context) (*Health, error) {
	var out Health
	if err := c.do(ctx, request{method: http.MethodGet, path: "/v1/health", public: true, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// request is one call of the API.
type request struct {
	method string
	path   string
	query  url.Values
	body   any
	// public: the endpoint answers without a key. The key is still sent when the client has one.
	public bool
	// retry: the request is safe to send again.
	retry bool
	// extraTimeout is added to the client's timeout, for a request the API holds open.
	extraTimeout time.Duration
	// idempotencyKey is sent as the Idempotency-Key header: what makes starting a run safe to send again.
	idempotencyKey string
	// accept is the Accept header; empty: application/json.
	accept string
}

// answer is what came back of a request that succeeded.
type answer struct {
	header http.Header
	data   []byte
}

const (
	backoffBase   = 500 * time.Millisecond
	backoffMax    = 8 * time.Second
	retryAfterMax = 60 * time.Second
)

func (c *Client) do(ctx context.Context, r request, out any) error {
	got, err := c.send(ctx, r)
	if err != nil {
		return err
	}
	if out == nil || len(got.data) == 0 {
		return nil
	}
	if err := json.Unmarshal(got.data, out); err != nil {
		return fmt.Errorf("enrichednow: %s %s: decoding the response: %w", r.method, r.path, err)
	}
	if keeper, ok := out.(interface{ setRaw([]byte) }); ok {
		keeper.setRaw(got.data)
	}
	return nil
}

// send sends the request, again while it is worth it, and returns the answer of the one that succeeded.
func (c *Client) send(ctx context.Context, r request) (*answer, error) {
	if c.apiKey == "" && !r.public {
		return nil, ErrMissingAPIKey
	}
	var payload []byte
	if r.body != nil {
		var err error
		if payload, err = json.Marshal(r.body); err != nil {
			return nil, fmt.Errorf("enrichednow: encoding the request body: %w", err)
		}
	}
	target := c.baseURL + r.path
	if len(r.query) > 0 {
		target += "?" + r.query.Encode()
	}
	attempts := 1
	if r.retry {
		attempts += c.maxRetries
	}
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			if err := c.sleep(ctx, c.delay(attempt, lastErr)); err != nil {
				return nil, err
			}
		}
		var got *answer
		got, lastErr = c.attempt(ctx, r, target, payload)
		if lastErr == nil {
			return got, nil
		}
		if !retryable(ctx, lastErr) {
			return nil, lastErr
		}
	}
	return nil, lastErr
}

// attempt sends the request once.
func (c *Client) attempt(ctx context.Context, r request, target string, payload []byte) (*answer, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout+r.extraTimeout)
		defer cancel()
	}
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, r.method, target, body)
	if err != nil {
		return nil, fmt.Errorf("enrichednow: %s %s: %w", r.method, r.path, err)
	}
	accept := r.accept
	if accept == "" {
		accept = "application/json"
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("User-Agent", c.userAgent)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	if r.idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", r.idempotencyKey)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &ConnectionError{Method: r.method, Path: r.path, Err: err}
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ConnectionError{Method: r.method, Path: r.path, Err: err}
	}
	if standing := standingFrom(resp.Header); standing != nil {
		c.standing.Store(standing)
	}
	if resp.StatusCode >= 400 {
		return nil, apiError(r, resp, data)
	}
	return &answer{header: resp.Header, data: data}, nil
}

// standingFrom reads the account's standing off an answer's headers; nil when it carries none.
func standingFrom(h http.Header) *Standing {
	if h.Get("X-Units-Limit") == "" && h.Get("X-RateLimit-Limit") == "" {
		return nil
	}
	number := func(name string) *int64 {
		v, err := strconv.ParseInt(strings.TrimSpace(h.Get(name)), 10, 64)
		if err != nil {
			return nil
		}
		return &v
	}
	s := &Standing{UnitsLimit: number("X-Units-Limit"), UnitsUsed: number("X-Units-Used"), UnitsRemaining: number("X-Units-Remaining"),
		UnitsTopup: number("X-Units-Topup"), RunsLimit: number("X-RateLimit-Limit"), RunsRemaining: number("X-RateLimit-Remaining"),
		Unmetered: h.Get("X-Units-Unmetered") == "true"}
	if at, err := time.Parse(time.RFC3339, h.Get("X-Units-Reset")); err == nil {
		s.UnitsReset = &at
	}
	return s
}

// ConnectionError is a request that got no answer: the connection failed, or the time ran out. Unwrap gives the
// cause, so errors.Is(err, context.DeadlineExceeded) tells a timeout apart.
type ConnectionError struct {
	// Method and Path say which request failed.
	Method, Path string
	// Err is the cause.
	Err error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("enrichednow: %s %s: %v", e.Method, e.Path, e.Err)
}

// Unwrap returns the cause.
func (e *ConnectionError) Unwrap() error { return e.Err }

func apiError(r request, resp *http.Response, data []byte) *APIError {
	e := &APIError{StatusCode: resp.StatusCode, Method: r.method, Path: r.path, Details: map[string]any{},
		RequestID: resp.Header.Get("X-Request-Id"), RetryAfter: retryAfter(resp.Header.Get("Retry-After"))}
	var body struct {
		Error struct {
			Code    string         `json:"code"`
			Message string         `json:"message"`
			Details map[string]any `json:"details"`
		} `json:"error"`
	}
	if json.Unmarshal(data, &body) == nil && body.Error.Code != "" {
		e.Code, e.Message = body.Error.Code, body.Error.Message
		if body.Error.Details != nil {
			e.Details = body.Error.Details
		}
	} else {
		e.Message = http.StatusText(resp.StatusCode)
	}
	return e
}

// retryAfter reads a Retry-After header: a number of seconds, or an HTTP date.
func retryAfter(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if seconds, err := strconv.ParseFloat(value, 64); err == nil {
		if seconds <= 0 {
			return 0
		}
		return time.Duration(seconds * float64(time.Second))
	}
	if at, err := http.ParseTime(value); err == nil {
		if d := time.Until(at); d > 0 {
			return d
		}
	}
	return 0
}

// retryable: a request that got no answer, a 429 that is about the rate of requests, and a 5xx.
func retryable(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	var transport *ConnectionError
	if errors.As(err, &transport) {
		return true
	}
	var api *APIError
	if errors.As(err, &api) {
		if api.StatusCode == http.StatusTooManyRequests {
			return api.Code != CodeTooManyRuns // that one clears when a run ends, not after a pause
		}
		return api.StatusCode >= 500
	}
	return false
}

// delay is the pause before retry number `attempt` (1 is the first retry): what Retry-After asked for, or else
// 0.5s doubling up to 8s, of which a random half to all is waited.
func (c *Client) delay(attempt int, last error) time.Duration {
	var api *APIError
	if errors.As(last, &api) && api.RetryAfter > 0 {
		return min(api.RetryAfter, retryAfterMax)
	}
	d := backoffBase << (attempt - 1)
	if d > backoffMax || d <= 0 {
		d = backoffMax
	}
	return time.Duration(float64(d) * (0.5 + 0.5*c.jitter()))
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func pathEscape(kind, id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("enrichednow: the %s is empty", kind)
	}
	return url.PathEscape(id), nil
}
