package enrichednow

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// The run endings a webhook endpoint can ask to be told about. EventRunSucceeded is the status ok.
const (
	EventRunSucceeded   = "run.succeeded"
	EventRunPartial     = "run.partial"
	EventRunFailed      = "run.failed"
	EventRunUnsupported = "run.unsupported"
	EventRunCancelled   = "run.cancelled"
)

// WebhookEvents is every event an endpoint can ask for.
var WebhookEvents = []string{EventRunSucceeded, EventRunPartial, EventRunFailed, EventRunUnsupported, EventRunCancelled}

// The headers a webhook request carries.
const (
	HeaderWebhookSignature = "X-Webhook-Signature"
	HeaderWebhookTimestamp = "X-Webhook-Timestamp"
)

// DefaultWebhookTolerance is how far from now a webhook request's timestamp may be, either way.
const DefaultWebhookTolerance = 300 * time.Second

// ErrWebhookVerification is returned by [VerifyWebhook] and [ParseWebhook] for a request that is not one this API
// signed with the secret, recently: answer it 400 and do nothing with it.
var ErrWebhookVerification = errors.New("enrichednow: the webhook request did not verify")

// WebhookEndpoint is a URL of yours that is sent a signed POST when a run ends with one of the events it asked for.
type WebhookEndpoint struct {
	rawBody
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Events      []string `json:"events"`
	// Enabled is false when it was switched off, by you or after repeated failures (DisabledReason says which).
	Enabled        bool       `json:"enabled"`
	DisabledReason *string    `json:"disabled_reason"`
	DisabledAt     *time.Time `json:"disabled_at"`
	LastSuccessAt  *time.Time `json:"last_success_at"`
	LastFailureAt  *time.Time `json:"last_failure_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

// WebhookEndpointWithSecret is an endpoint as it is created or its secret rotated: the only time the signing
// secret is shown. Keep it: it is what [VerifyWebhook] checks a request with.
type WebhookEndpointWithSecret struct {
	WebhookEndpoint
	Secret string `json:"secret"`
}

// WebhookEndpointList is the account's endpoints.
type WebhookEndpointList struct {
	rawBody
	Endpoints []WebhookEndpoint `json:"endpoints"`
	// RunSecretSet says whether the account has a secret for the webhook_url of a single run.
	RunSecretSet bool `json:"run_secret_set"`
}

// WebhookSecret is a signing secret, shown once.
type WebhookSecret struct {
	rawBody
	Secret string `json:"secret"`
}

// WebhookDelivery is one event sent to one URL, with every attempt made at it.
type WebhookDelivery struct {
	rawBody
	ID string `json:"id"`
	// EndpointID is empty for a run's own webhook_url.
	EndpointID *string `json:"endpoint_id"`
	RunID      string  `json:"run_id"`
	// EventID is the same on every attempt and when the delivery is sent again: use it to ignore a repeat.
	EventID string `json:"event_id"`
	Event   string `json:"event"`
	URL     string `json:"url"`
	// Status is pending, succeeded or failed.
	Status         string     `json:"status"`
	Attempts       int        `json:"attempts"`
	MaxAttempts    int        `json:"max_attempts"`
	NextAttemptAt  *time.Time `json:"next_attempt_at"`
	ResponseStatus *int       `json:"response_status"`
	Error          *string    `json:"error"`
	DurationMS     *int       `json:"duration_ms"`
	LastAttemptAt  *time.Time `json:"last_attempt_at"`
	CreatedAt      time.Time  `json:"created_at"`
	FinishedAt     *time.Time `json:"finished_at"`
	// ResendOf is the delivery this one was sent again from.
	ResendOf *string `json:"resend_of"`
}

// WebhookDeliveryPage is one page of the delivery log, newest first.
type WebhookDeliveryPage struct {
	rawBody
	Deliveries []WebhookDelivery `json:"deliveries"`
	// NextCursor is the cursor of the next page; empty on the last one.
	NextCursor string `json:"next_cursor"`
}

// WebhookEvent is what a webhook request carries. [ParseWebhook] returns it from a request that verified.
type WebhookEvent struct {
	rawBody
	// ID (evt_…) is the same on every attempt and when the delivery is sent again: use it to ignore a repeat.
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Data      struct {
		// Run is the run that ended, as [RunService.Get] returns it.
		Run Run `json:"run"`
	} `json:"data"`
}

// WebhookService manages the endpoints told when a run ends, their signing secrets, and the delivery log.
type WebhookService struct {
	client *Client
	// Deliveries is the log of what was sent, and sending one again.
	Deliveries *WebhookDeliveryService
}

// WebhookEndpointParams are what [WebhookService.Create] takes: an https URL and the events it is told about.
type WebhookEndpointParams struct {
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Description string   `json:"description,omitempty"`
}

// WebhookEndpointChange is what [WebhookService.Update] changes; a nil field is left as it is.
type WebhookEndpointChange struct {
	URL         *string  `json:"url,omitempty"`
	Events      []string `json:"events,omitempty"`
	Description *string  `json:"description,omitempty"`
	// Enabled true switches an endpoint that was switched off back on.
	Enabled *bool `json:"enabled,omitempty"`
}

// List returns the account's endpoints.
func (s *WebhookService) List(ctx context.Context) (*WebhookEndpointList, error) {
	var out WebhookEndpointList
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/webhooks", retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create adds an endpoint. The answer is the only time its signing secret is shown. It is not sent again after a
// failure that may have reached the API, so that one call never makes two endpoints.
func (s *WebhookService) Create(ctx context.Context, params WebhookEndpointParams) (*WebhookEndpointWithSecret, error) {
	var out WebhookEndpointWithSecret
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/webhooks", body: params}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one endpoint.
func (s *WebhookService) Get(ctx context.Context, endpointID string) (*WebhookEndpoint, error) {
	id, err := pathEscape("endpoint id", endpointID)
	if err != nil {
		return nil, err
	}
	var out WebhookEndpoint
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/webhooks/" + id, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update changes an endpoint's URL, events, description, or whether it is enabled.
func (s *WebhookService) Update(ctx context.Context, endpointID string, change WebhookEndpointChange) (*WebhookEndpoint, error) {
	id, err := pathEscape("endpoint id", endpointID)
	if err != nil {
		return nil, err
	}
	if change.URL == nil && change.Events == nil && change.Description == nil && change.Enabled == nil {
		return nil, errors.New("enrichednow: give at least one of URL, Events, Description, Enabled")
	}
	var out WebhookEndpoint
	if err := s.client.do(ctx, request{method: http.MethodPatch, path: "/v1/webhooks/" + id, body: change, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes an endpoint. Its deliveries stay in the log.
func (s *WebhookService) Delete(ctx context.Context, endpointID string) error {
	id, err := pathEscape("endpoint id", endpointID)
	if err != nil {
		return err
	}
	return s.client.do(ctx, request{method: http.MethodDelete, path: "/v1/webhooks/" + id, retry: true}, nil)
}

// RotateSecret gives an endpoint a new signing secret, shown in the answer, and ends the old one.
func (s *WebhookService) RotateSecret(ctx context.Context, endpointID string) (*WebhookEndpointWithSecret, error) {
	id, err := pathEscape("endpoint id", endpointID)
	if err != nil {
		return nil, err
	}
	var out WebhookEndpointWithSecret
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/webhooks/" + id + "/rotate-secret"}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRunSecret makes (or replaces) the account's secret for the webhook_url of a single run, shown once.
func (s *WebhookService) CreateRunSecret(ctx context.Context) (*WebhookSecret, error) {
	var out WebhookSecret
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/webhooks/run-secret"}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// WebhookDeliveryService is the delivery log.
type WebhookDeliveryService struct{ client *Client }

// WebhookDeliveryListParams filter and page [WebhookDeliveryService.List].
type WebhookDeliveryListParams struct {
	// Limit is the deliveries per page: 1 to 100. Zero leaves it to the API.
	Limit      int
	Cursor     string
	EndpointID string
	RunID      string
	// Status is pending, succeeded or failed.
	Status string
}

// List returns one page of the delivery log, newest first.
func (s *WebhookDeliveryService) List(ctx context.Context, params *WebhookDeliveryListParams) (*WebhookDeliveryPage, error) {
	q := url.Values{}
	if params != nil {
		setInt(q, "limit", params.Limit)
		setString(q, "cursor", params.Cursor)
		setString(q, "endpoint_id", params.EndpointID)
		setString(q, "run_id", params.RunID)
		setString(q, "status", params.Status)
	}
	var out WebhookDeliveryPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/webhooks/deliveries", query: q, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAll iterates over the delivery log, page after page.
func (s *WebhookDeliveryService) ListAll(ctx context.Context, params *WebhookDeliveryListParams) *Iter[WebhookDelivery] {
	p := WebhookDeliveryListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]WebhookDelivery, string, error) {
		p.Cursor = cursor
		page, err := s.List(ctx, &p)
		if err != nil {
			return nil, "", err
		}
		return page.Deliveries, page.NextCursor, nil
	})
}

// Resend sends a delivery's event again, as a new delivery with the same event id.
func (s *WebhookDeliveryService) Resend(ctx context.Context, deliveryID string) (*WebhookDelivery, error) {
	id, err := pathEscape("delivery id", deliveryID)
	if err != nil {
		return nil, err
	}
	var out WebhookDelivery
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/webhooks/deliveries/" + id + "/resend"}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// VerifyWebhookSignature checks that body was signed with secret at timestamp, and that timestamp is within
// tolerance of now (0: [DefaultWebhookTolerance]). The signature is "v1=" and the HMAC-SHA256, in hex, of the
// timestamp, a dot and the body exactly as it arrived; it is compared in constant time. It needs no client and
// makes no request. The error is [ErrWebhookVerification] for anything that does not verify.
func VerifyWebhookSignature(secret, timestamp string, body []byte, signature string, tolerance time.Duration) error {
	if tolerance <= 0 {
		tolerance = DefaultWebhookTolerance
	}
	if secret == "" || signature == "" || len(timestamp) == 0 || len(timestamp) > 16 {
		return ErrWebhookVerification
	}
	sent, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil || strings.TrimLeft(timestamp, "0123456789") != "" {
		return ErrWebhookVerification
	}
	if math.Abs(float64(time.Now().Unix()-sent)) > tolerance.Seconds() {
		return ErrWebhookVerification
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "."))
	mac.Write(body)
	expected := "v1=" + hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return ErrWebhookVerification
	}
	return nil
}

// VerifyWebhook checks a webhook request by its headers: see [VerifyWebhookSignature]. body must be the raw bytes
// that arrived, read before any JSON decoding.
//
//	body, _ := io.ReadAll(r.Body)
//	if err := enrichednow.VerifyWebhook(secret, r.Header, body); err != nil {
//		http.Error(w, "bad signature", http.StatusBadRequest)
//		return
//	}
func VerifyWebhook(secret string, header http.Header, body []byte) error {
	return VerifyWebhookSignature(secret, header.Get(HeaderWebhookTimestamp), body, header.Get(HeaderWebhookSignature), 0)
}

// ParseWebhook verifies a webhook request as [VerifyWebhook] does and returns the event it carries.
func ParseWebhook(secret string, header http.Header, body []byte) (*WebhookEvent, error) {
	if err := VerifyWebhook(secret, header, body); err != nil {
		return nil, err
	}
	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil || event.Type == "" {
		return nil, errors.Join(ErrWebhookVerification, errors.New("enrichednow: the webhook request verified, but its body is not an event"))
	}
	event.setRaw(body)
	return &event, nil
}
