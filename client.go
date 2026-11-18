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
	return c, nil
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
}

const (
	backoffBase   = 500 * time.Millisecond
	backoffMax    = 8 * time.Second
	retryAfterMax = 60 * time.Second
)

func (c *Client) do(ctx context.Context, r request, out any) error {
	if c.apiKey == "" && !r.public {
		return ErrMissingAPIKey
	}
	var payload []byte
	if r.body != nil {
		var err error
		if payload, err = json.Marshal(r.body); err != nil {
			return fmt.Errorf("enrichednow: encoding the request body: %w", err)
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
				return err
			}
		}
		lastErr = c.attempt(ctx, r, target, payload, out)
		if lastErr == nil || !retryable(ctx, lastErr) {
			return lastErr
		}
	}
	return lastErr
}

// attempt sends the request once and decodes the answer.
func (c *Client) attempt(ctx context.Context, r request, target string, payload []byte, out any) error {
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
		return fmt.Errorf("enrichednow: %s %s: %w", r.method, r.path, err)
	}
	req.Header.Set("Accept", "application/json")
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
		return &ConnectionError{Method: r.method, Path: r.path, Err: err}
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ConnectionError{Method: r.method, Path: r.path, Err: err}
	}
	if resp.StatusCode >= 400 {
		return apiError(r, resp, data)
	}
	if out == nil || len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("enrichednow: %s %s: decoding the response: %w", r.method, r.path, err)
	}
	if keeper, ok := out.(interface{ setRaw([]byte) }); ok {
		keeper.setRaw(data)
	}
	return nil
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
