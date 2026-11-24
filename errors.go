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
	// ErrPermission: the account may not do this (any 403), e.g. a data request without Datasets access
	// (forbidden), or a schedule the plan does not allow ([ErrPlanLimitReached]).
	ErrPermission = errors.New("enrichednow: not allowed")
	// ErrNotFound: there is no such scraper, run, company, posting, formula or saved formula (any 404: not_found,
	// formula_not_found).
	ErrNotFound = errors.New("enrichednow: not found")
	// ErrConflict: any 409. The idempotency key was already used for a different request (idempotency_conflict),
	// or what was asked conflicts with how things stand: [ErrAttestationRequired], [ErrFormulaUnavailable],
	// [ErrRequirementMissing], [ErrSavedFormulaArchived].
	ErrConflict = errors.New("enrichednow: idempotency key already used")
	// ErrValidation: the request does not fit (any 422: invalid_input, unit_cap_too_low). [APIError.Fields] lists
	// the problems of an invalid_input; [APIError.MinUnitCap] says the least unit cap of a unit_cap_too_low.
	ErrValidation = errors.New("enrichednow: invalid input")
	// ErrQuotaExceeded: this month's allowance is used up (402 quota_exceeded).
	ErrQuotaExceeded = errors.New("enrichednow: quota exceeded")
	// ErrLimitReached: the run would go past a spend limit the account set itself (402 limit_reached): its
	// monthly limit (Details "scope" = "account"), the API key's ("scope" = "key", with "key_id"), a saved
	// formula's monthly unit cap ("scope" = "saved_formula", with "saved_formula_id"), or a formula run's own
	// ("scope" = "formula_run"). Details also has "limit_units", "used_units" and "units_needed". Units may be left
	// on the plan: raising the limit is what helps. Such an error matches [ErrQuotaExceeded] too, so a check for
	// that one still catches it.
	ErrLimitReached = errors.New("enrichednow: spend limit reached")
	// ErrRateLimit: any 429. It covers both too_many_runs and rate_limited.
	ErrRateLimit = errors.New("enrichednow: rate limited")
	// ErrTooManyRuns: the account is at its limit of runs going at the same time (429 too_many_runs).
	// No run was created; send the request again once one of your runs has ended.
	ErrTooManyRuns = errors.New("enrichednow: too many runs")
	// ErrServer: something broke on the API's side (any 5xx).
	ErrServer = errors.New("enrichednow: server error")
	// ErrDataUnavailable: the dataset cannot be reached at the moment (503 data_unavailable).
	ErrDataUnavailable = errors.New("enrichednow: data unavailable")

	// ErrAttestationRequired: the formula returns people, and an owner or an admin of the account must first
	// accept the terms for it, in the console (409 attestation_required). Details has "kind" and "text_version";
	// [AttestationService.List] shows the text. A key cannot accept them. It matches [ErrConflict] too.
	ErrAttestationRequired = errors.New("enrichednow: the terms for this formula are not accepted")
	// ErrFormulaUnavailable: the formula cannot be run now (409 formula_unavailable), e.g. it is deprecated
	// (Details "status"), or a saved formula's major version is no longer offered ("formula_id", "major"). It
	// matches [ErrConflict] too.
	ErrFormulaUnavailable = errors.New("enrichednow: formula unavailable")
	// ErrRequirementMissing: the formula needs a product or a source this account lacks (409
	// requirement_missing). [APIError.Missing] lists them, such as "product:data". It matches [ErrConflict] too.
	ErrRequirementMissing = errors.New("enrichednow: requirement missing")
	// ErrSavedFormulaArchived: the saved formula is archived, and cannot be changed or run (409
	// saved_formula_archived). Duplicate it ([SavedFormulaService.Duplicate]) to use it again. It matches
	// [ErrConflict] too.
	ErrSavedFormulaArchived = errors.New("enrichednow: saved formula archived")
	// ErrPlanLimitReached: the account's plan does not allow this (403 plan_limit_reached). Details "limit" says
	// which limit ("schedules", "min_interval_minutes", "destinations"), "allowed" what the plan allows, and
	// "current" or "asked" what was there or asked for. [SavedFormulaService.Limits] shows them all. It matches
	// [ErrPermission] too.
	ErrPlanLimitReached = errors.New("enrichednow: plan limit reached")

	// ErrRunFailed matches every [*RunFailedError] and [*FormulaRunFailedError].
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
	CodeLimitReached    = "limit_reached"
	CodeTooManyRuns     = "too_many_runs"
	CodeRateLimited     = "rate_limited"
	CodeInternal        = "internal"
	CodeDataUnavailable = "data_unavailable"

	CodeIdempotencyConflict  = "idempotency_conflict"
	CodeFormulaNotFound      = "formula_not_found"
	CodeFormulaUnavailable   = "formula_unavailable"
	CodeRequirementMissing   = "requirement_missing"
	CodeUnitCapTooLow        = "unit_cap_too_low"
	CodeAttestationRequired  = "attestation_required"
	CodeSavedFormulaArchived = "saved_formula_archived"
	CodePlanLimitReached     = "plan_limit_reached"
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
	// "allowance_units", and on starting a run "topup_units" and "units_needed" (the least the run costs: its
	// scraper's units a request, or a record for one record); for limit_reached "scope", "key_id" (a key's),
	// "saved_formula_id" (a saved formula's), "limit_units", "used_units" and "units_needed"; for too_many_runs
	// "active_runs" and "max_concurrent_runs"; for a forbidden answered because the account must confirm its email
	// first, "reason" = "email_not_verified".
	//
	// For the formulas: attestation_required has "kind" and "text_version"; formula_unavailable "status", or
	// "formula_id" and "major"; requirement_missing "missing" ([APIError.Missing]); unit_cap_too_low
	// "min_unit_cap" ([APIError.MinUnitCap]); saved_formula_archived "saved_formula_id"; plan_limit_reached
	// "limit", "allowed", and "current" or "asked".
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
		return e.StatusCode == http.StatusForbidden || e.Code == CodeForbidden || e.Code == CodePlanLimitReached
	case ErrNotFound:
		return e.StatusCode == http.StatusNotFound || e.Code == CodeNotFound || e.Code == CodeFormulaNotFound
	case ErrConflict:
		return e.StatusCode == http.StatusConflict
	case ErrValidation:
		return e.StatusCode == http.StatusUnprocessableEntity || e.Code == CodeInvalidInput || e.Code == CodeUnitCapTooLow
	case ErrQuotaExceeded:
		return e.StatusCode == http.StatusPaymentRequired || e.Code == CodeQuotaExceeded || e.Code == CodeLimitReached
	case ErrLimitReached:
		return e.Code == CodeLimitReached
	case ErrRateLimit:
		return e.StatusCode == http.StatusTooManyRequests
	case ErrTooManyRuns:
		return e.Code == CodeTooManyRuns
	case ErrServer:
		return e.StatusCode >= 500
	case ErrDataUnavailable:
		return e.Code == CodeDataUnavailable
	case ErrAttestationRequired:
		return e.Code == CodeAttestationRequired
	case ErrFormulaUnavailable:
		return e.Code == CodeFormulaUnavailable
	case ErrRequirementMissing:
		return e.Code == CodeRequirementMissing
	case ErrSavedFormulaArchived:
		return e.Code == CodeSavedFormulaArchived
	case ErrPlanLimitReached:
		return e.Code == CodePlanLimitReached
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

// Missing returns what a requirement_missing error says the account lacks, such as "product:data" or a scraper id;
// nil for any other error.
func (e *APIError) Missing() []string {
	raw, _ := e.Details["missing"].([]any)
	var out []string
	for _, item := range raw {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// MinUnitCap returns the least unit cap a unit_cap_too_low error asks for, and whether the error said.
func (e *APIError) MinUnitCap() (int, bool) {
	n, ok := e.Details["min_unit_cap"].(float64)
	return int(n), ok
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

// FormulaRunFailedError says that a formula run ended failed, unsupported or cancelled. The request itself worked;
// the run did not. It matches [ErrRunFailed] with [errors.Is]. A cancelled run keeps the rows it wrote before it
// stopped: read them with [FormulaRunService.Rows].
type FormulaRunFailedError struct {
	// Run is the run as it ended.
	Run *FormulaRun
}

// Status is the status the run ended with: failed, unsupported or cancelled.
func (e *FormulaRunFailedError) Status() RunStatus { return e.Run.Status }

// Reason is the run's stable reason code, such as "source_failed:<scraper>"; empty when the run gave none.
func (e *FormulaRunFailedError) Reason() string { return e.Run.Reason }

func (e *FormulaRunFailedError) Error() string {
	if e.Run.Reason == "" {
		return fmt.Sprintf("enrichednow: formula run %s ended %s", e.Run.ID, e.Run.Status)
	}
	return fmt.Sprintf("enrichednow: formula run %s ended %s (reason: %s)", e.Run.ID, e.Run.Status, e.Run.Reason)
}

// Is reports whether target is [ErrRunFailed].
func (e *FormulaRunFailedError) Is(target error) bool { return target == ErrRunFailed }
