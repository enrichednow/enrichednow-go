package enrichednow

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Sentinel errors. An [*APIError] matches them with [errors.Is]:
//
//	if errors.Is(err, enrichednow.ErrNotFound) { ... }
var (
	// ErrMissingAPIKey is returned by [NewClient] when no key was given and ENRICHEDNOW_API_KEY is not set,
	// and by a call that needs a key on a client built [WithoutAPIKey].
	ErrMissingAPIKey = errors.New("enrichednow: no API key: pass enrichednow.WithAPIKey(...) or set " + EnvAPIKey)

	// ErrAuthentication: the API key is missing, mistyped or revoked (401 unauthorized).
	ErrAuthentication = errors.New("enrichednow: authentication failed")
	// ErrPermission: the account may not do this (403 forbidden), e.g. a data request without Datasets access.
	ErrPermission = errors.New("enrichednow: not allowed")
	// ErrNotFound: there is no such scraper, run, company or posting (404 not_found).
	ErrNotFound = errors.New("enrichednow: not found")
	// ErrValidation: the request does not fit (422 invalid_input). [APIError.Fields] lists the problems.
	ErrValidation = errors.New("enrichednow: invalid input")
	// ErrQuotaExceeded: this month's allowance is used up (402 quota_exceeded).
	ErrQuotaExceeded = errors.New("enrichednow: quota exceeded")
	// ErrRateLimit: any 429. It covers both too_many_runs and rate_limited.
	ErrRateLimit = errors.New("enrichednow: rate limited")
	// ErrTooManyRuns: the account is at its limit of runs going at the same time (429 too_many_runs).
	// No run was created; send the request again once one of your runs has ended.
	ErrTooManyRuns = errors.New("enrichednow: too many runs")
	// ErrServer: something broke on the API's side (any 5xx).
	ErrServer = errors.New("enrichednow: server error")
	// ErrDataUnavailable: the dataset cannot be reached at the moment (503 data_unavailable).
	ErrDataUnavailable = errors.New("enrichednow: data unavailable")

	// ErrRunFailed matches every [*RunFailedError].
	ErrRunFailed = errors.New("enrichednow: run did not succeed")
	// ErrRecordsExpired is reported by a records iterator when the records are gone: they are kept for 7 days
	// after the run ended.
	ErrRecordsExpired = errors.New("enrichednow: the records of this run have expired")
)

// The error codes of the API. Match on these, not on the message.
const (
	CodeUnauthorized    = "unauthorized"
	CodeForbidden       = "forbidden"
	CodeNotFound        = "not_found"
	CodeInvalidInput    = "invalid_input"
	CodeQuotaExceeded   = "quota_exceeded"
	CodeTooManyRuns     = "too_many_runs"
	CodeRateLimited     = "rate_limited"
	CodeInternal        = "internal"
	CodeDataUnavailable = "data_unavailable"
)

// APIError is an error answered by the API: an HTTP status of 4xx or 5xx with the body
// {"error": {"code", "message", "details"}}.
type APIError struct {
	// StatusCode is the HTTP status.
	StatusCode int
	// Code is the API's stable snake_case code, e.g. "not_found". Empty when the answer was not the API's own
	// error body.
	Code string
	// Message is a sentence for a person. Its wording can change; do not parse it.
	Message string
	// Details is always non-nil. For invalid_input it holds "fields"; for quota_exceeded "used_units" and
	// "allowance_units"; for too_many_runs "active_runs" and "max_concurrent_runs"; for a forbidden answered
	// because the account must confirm its email first, "reason" = "email_not_verified".
	Details map[string]any
	// RequestID is the X-Request-Id header of the response, when one was sent.
	RequestID string
	// RetryAfter is the wait the Retry-After header asked for; zero when there was none.
	RetryAfter time.Duration
	// Method and Path say which request failed.
	Method string
	Path   string
}

// FieldError is one problem of an invalid_input error.
type FieldError struct {
	// Path is where the problem is, as the caller wrote it: "input.board_url", "wait", "limit".
	Path    string
	Message string
}

func (e *APIError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "enrichednow: %s %s: %d", e.Method, e.Path, e.StatusCode)
	if e.Code != "" {
		b.WriteString(" " + e.Code)
	}
	if e.Message != "" {
		b.WriteString(": " + e.Message)
	}
	for _, f := range e.Fields() {
		fmt.Fprintf(&b, " [%s: %s]", f.Path, f.Message)
	}
	if e.RequestID != "" {
		fmt.Fprintf(&b, " (request %s)", e.RequestID)
	}
	return b.String()
}

// Is reports whether the error is of the kind a sentinel of this package names.
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrAuthentication:
		return e.StatusCode == http.StatusUnauthorized || e.Code == CodeUnauthorized
	case ErrPermission:
		return e.StatusCode == http.StatusForbidden || e.Code == CodeForbidden
	case ErrNotFound:
		return e.StatusCode == http.StatusNotFound || e.Code == CodeNotFound
	case ErrValidation:
		return e.StatusCode == http.StatusUnprocessableEntity || e.Code == CodeInvalidInput
	case ErrQuotaExceeded:
		return e.StatusCode == http.StatusPaymentRequired || e.Code == CodeQuotaExceeded
	case ErrRateLimit:
		return e.StatusCode == http.StatusTooManyRequests
	case ErrTooManyRuns:
		return e.Code == CodeTooManyRuns
	case ErrServer:
		return e.StatusCode >= 500
	case ErrDataUnavailable:
		return e.Code == CodeDataUnavailable
	}
	return false
}

// Fields returns the problems of an invalid_input error, one per field; nil for any other error.
func (e *APIError) Fields() []FieldError {
	raw, _ := e.Details["fields"].([]any)
	var out []FieldError
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		path, _ := m["path"].(string)
		message, _ := m["message"].(string)
		out = append(out, FieldError{Path: path, Message: message})
	}
	return out
}

// RunFailedError says that a run ended, but not with records: its status is failed, unsupported or cancelled.
// The request itself worked; the scrape did not. It matches [ErrRunFailed] with [errors.Is].
type RunFailedError struct {
	// Run is the run as it ended.
	Run *Run
}

// Status is the status the run ended with: failed, unsupported or cancelled.
func (e *RunFailedError) Status() RunStatus { return e.Run.Status }

// Reason is the run's stable reason code, e.g. "render_required" or "fetch_failed"; empty when the run gave none.
func (e *RunFailedError) Reason() string { return e.Run.Reason }

func (e *RunFailedError) Error() string {
	if e.Run.Reason == "" {
		return fmt.Sprintf("enrichednow: run %s ended %s", e.Run.ID, e.Run.Status)
	}
	return fmt.Sprintf("enrichednow: run %s ended %s (reason: %s)", e.Run.ID, e.Run.Status, e.Run.Reason)
}

// Is reports whether target is [ErrRunFailed].
func (e *RunFailedError) Is(target error) bool { return target == ErrRunFailed }
