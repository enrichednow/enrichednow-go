package enrichednow

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// SavedFormulaService manages saved formulas: a formula with its params and caps, run again by hand or on a
// schedule. A saved formula's runs deliver only rows it did not deliver before.
type SavedFormulaService struct{ client *Client }

// ScheduleParams is when a saved formula runs, as [SavedFormulaService.Create], [SavedFormulaService.Update] and
// [SavedFormulaService.Duplicate] take it.
type ScheduleParams struct {
	// Cadence is [CadenceHourly], [CadenceDaily], [CadenceWeekdays], [CadenceWeekly] or [CadenceMonthly]: one the
	// formula allows ([Formula.Schedules]), and the plan does ([PlanLimits.MinIntervalMinutes]).
	Cadence string `json:"cadence"`
	// LocalTime is HH:MM, 24-hour, in Timezone. For hourly, only the minute counts.
	LocalTime string `json:"local_time"`
	// Timezone is an IANA time zone, such as Europe/London, America/New_York or UTC.
	Timezone string `json:"timezone"`
	// DayOfWeek is the day of a weekly schedule, 0 for Monday to 6 for Sunday; nil: Monday.
	DayOfWeek *int `json:"day_of_week,omitempty"`
	// DayOfMonth is the day of a monthly schedule, 1 to 31 (past a short month's end, its last day); nil: the 1st.
	DayOfMonth *int `json:"day_of_month,omitempty"`
}

// SavedFormulaCreateParams are the options of [SavedFormulaService.Create]. A field left at its zero value is not
// sent, and the API's default applies.
type SavedFormulaCreateParams struct {
	// Name is its name; empty: one made from the params, which follows them when they change.
	Name string
	// Schedule is when it runs; nil: only by hand.
	Schedule *ScheduleParams
	// UnitCap is the most one run may spend; zero: the formula's default. Above that, owners and admins only.
	UnitCap int
	// RowCap is the most rows one run may produce; zero: the formula's default.
	RowCap int
	// MonthlyUnitCap is the most its runs may spend in a calendar month (UTC); zero: no monthly cap.
	MonthlyUnitCap int
	// Notify is when scheduled runs send mail about new rows: [NotifyEachRun] (the default), [NotifyDigest] or
	// [NotifyNone].
	Notify string
}

// Create saves a formula with its params. params is a parameter type of this package, a struct of yours or a
// map, as [FormulaService.Run] takes it; a parameter type made for another formula is refused before anything is
// sent.
//
// It is not sent again after a failure that may have reached the API, so that one call never makes two saved
// formulas. A schedule the plan does not allow fails with [ErrPlanLimitReached]; one for a formula whose terms are
// not accepted, with [ErrAttestationRequired].
func (s *SavedFormulaService) Create(ctx context.Context, formulaID string, params any, opts *SavedFormulaCreateParams) (*SavedFormula, error) {
	if _, err := pathEscape("formula id", formulaID); err != nil {
		return nil, err
	}
	sent, err := formulaParams(formulaID, params)
	if err != nil {
		return nil, err
	}
	body := map[string]any{"formula_id": formulaID, "params": sent}
	if opts != nil {
		if opts.Name != "" {
			body["name"] = opts.Name
		}
		if opts.Schedule != nil {
			body["schedule"] = opts.Schedule
		}
		setCaps(body, false, opts.UnitCap, opts.RowCap)
		if opts.MonthlyUnitCap > 0 {
			body["monthly_unit_cap"] = opts.MonthlyUnitCap
		}
		if opts.Notify != "" {
			body["notify"] = opts.Notify
		}
	}
	var out SavedFormula
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/saved-formulas", body: body}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SavedFormulaListParams filter and page [SavedFormulaService.List].
type SavedFormulaListParams struct {
	// FormulaID keeps the saved formulas of one formula.
	FormulaID string
	// Status keeps those with this status: active, paused or archived. Without it, archived ones are left out.
	Status string
	// Limit is the saved formulas per page: 1 to 100. Zero leaves it to the API (50).
	Limit int
	// Cursor is the NextCursor of the page before.
	Cursor string
}

// List returns one page of the account's saved formulas.
func (s *SavedFormulaService) List(ctx context.Context, params *SavedFormulaListParams) (*SavedFormulaPage, error) {
	q := url.Values{}
	if params != nil {
		setString(q, "formula_id", params.FormulaID)
		setString(q, "status", params.Status)
		setInt(q, "limit", params.Limit)
		setString(q, "cursor", params.Cursor)
	}
	var out SavedFormulaPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/saved-formulas", query: q, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAll iterates over the account's saved formulas, fetching page after page.
func (s *SavedFormulaService) ListAll(ctx context.Context, params *SavedFormulaListParams) *Iter[SavedFormula] {
	p := SavedFormulaListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]SavedFormula, string, error) {
		p.Cursor = cursor
		page, err := s.List(ctx, &p)
		if err != nil {
			return nil, "", err
		}
		return page.SavedFormulas, page.NextCursor, nil
	})
}

// Get returns one saved formula.
func (s *SavedFormulaService) Get(ctx context.Context, savedFormulaID string) (*SavedFormula, error) {
	return s.call(ctx, http.MethodGet, savedFormulaID, "", nil, true)
}

// call sends a request about one saved formula that answers with it.
func (s *SavedFormulaService) call(ctx context.Context, method, savedFormulaID, suffix string, body any, retry bool) (*SavedFormula, error) {
	id, err := pathEscape("saved formula id", savedFormulaID)
	if err != nil {
		return nil, err
	}
	var out SavedFormula
	if err := s.client.do(ctx, request{method: method, path: "/v1/saved-formulas/" + id + suffix, body: body, retry: retry}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SavedFormulaUpdateParams is what [SavedFormulaService.Update] changes. A field left at its zero value is left
// as it is; the Remove fields and AutoName send the API's null, which takes a setting away.
type SavedFormulaUpdateParams struct {
	// Name renames it. A rename stops the name from following the params.
	Name string
	// AutoName goes back to a name made from the params, which follows them when they change.
	AutoName bool
	// Params replace the params: the whole set, as [FormulaService.Run] takes them; what is left out takes the
	// formula's default. Rows already delivered stay delivered, unless ResetSeen.
	Params any
	// ResetSeen forgets what was delivered, so that rows already delivered come again (owners and admins).
	ResetSeen bool
	// UnitCap and RowCap change the caps of one run.
	UnitCap int
	RowCap  int
	// MonthlyUnitCap sets the most its runs may spend in a calendar month; RemoveMonthlyUnitCap removes the cap.
	MonthlyUnitCap       int
	RemoveMonthlyUnitCap bool
	// Schedule sets when it runs; RemoveSchedule makes it run only by hand.
	Schedule       *ScheduleParams
	RemoveSchedule bool
	// Notify changes when scheduled runs send mail: [NotifyEachRun], [NotifyDigest] or [NotifyNone].
	Notify string
}

// Update changes a saved formula, and returns it as it is now. Only what change sets is sent; an empty change is
// refused before anything is sent. It is safe to repeat. An archived saved formula cannot be changed:
// [ErrSavedFormulaArchived].
func (s *SavedFormulaService) Update(ctx context.Context, savedFormulaID string, change SavedFormulaUpdateParams) (*SavedFormula, error) {
	body := map[string]any{}
	switch {
	case change.Name != "" && change.AutoName:
		return nil, errors.New("enrichednow: give Name or AutoName, not both")
	case change.Name != "":
		body["name"] = change.Name
	case change.AutoName:
		body["name"] = nil
	}
	if change.Params != nil {
		body["params"] = change.Params
	}
	if change.ResetSeen {
		body["reset_seen"] = true
	}
	setCaps(body, false, change.UnitCap, change.RowCap)
	switch {
	case change.MonthlyUnitCap > 0 && change.RemoveMonthlyUnitCap:
		return nil, errors.New("enrichednow: give MonthlyUnitCap or RemoveMonthlyUnitCap, not both")
	case change.MonthlyUnitCap > 0:
		body["monthly_unit_cap"] = change.MonthlyUnitCap
	case change.RemoveMonthlyUnitCap:
		body["monthly_unit_cap"] = nil
	}
	switch {
	case change.Schedule != nil && change.RemoveSchedule:
		return nil, errors.New("enrichednow: give Schedule or RemoveSchedule, not both")
	case change.Schedule != nil:
		body["schedule"] = change.Schedule
	case change.RemoveSchedule:
		body["schedule"] = nil
	}
	if change.Notify != "" {
		body["notify"] = change.Notify
	}
	if len(body) == 0 {
		return nil, errors.New("enrichednow: nothing to change: give at least one field of SavedFormulaUpdateParams")
	}
	return s.call(ctx, http.MethodPatch, savedFormulaID, "", body, true)
}

// Archive archives a saved formula, and returns it archived: its schedule stops and it does not run again. Its
// runs and rows are kept as long as they would be. An archived saved formula can be duplicated.
func (s *SavedFormulaService) Archive(ctx context.Context, savedFormulaID string) (*SavedFormula, error) {
	return s.call(ctx, http.MethodDelete, savedFormulaID, "", nil, true)
}

// SavedFormulaDuplicateParams are the options of [SavedFormulaService.Duplicate].
type SavedFormulaDuplicateParams struct {
	// Name is the copy's name; empty: one made from the params.
	Name string
	// CopySeen starts the copy with what the original delivered, so that it does not deliver it again. Without
	// it the copy starts empty.
	CopySeen bool
	// Schedule is the copy's schedule; nil: none.
	Schedule *ScheduleParams
}

// Duplicate makes a copy of a saved formula, archived or not, and returns the copy. It is not sent again after a
// failure that may have reached the API.
func (s *SavedFormulaService) Duplicate(ctx context.Context, savedFormulaID string, params *SavedFormulaDuplicateParams) (*SavedFormula, error) {
	body := map[string]any{}
	if params != nil {
		if params.Name != "" {
			body["name"] = params.Name
		}
		if params.CopySeen {
			body["copy_seen"] = true
		}
		if params.Schedule != nil {
			body["schedule"] = params.Schedule
		}
	}
	return s.call(ctx, http.MethodPost, savedFormulaID, "/duplicate", body, false)
}

// SavedFormulaRunParams are the options of [SavedFormulaService.Run].
type SavedFormulaRunParams struct {
	// Preview runs at most 10 rows and 50 units; nothing it returns is remembered as delivered.
	Preview bool
	// IncludeSeen delivers rows already delivered too. They are still remembered.
	IncludeSeen bool
	// Wait asks the API to hold the request open until the run ends or this much time has passed: 0 to [MaxWait].
	Wait time.Duration
	// IdempotencyKey names this start; see [FormulaRunParams.IdempotencyKey].
	IdempotencyKey string
}

// Run starts a run of a saved formula now, with its params and caps (the monthly cap lowers the unit cap to what
// is left of it), and returns it. A paused saved formula can still be run by hand. The request carries an
// Idempotency-Key and is retried as [FormulaService.Run] is.
func (s *SavedFormulaService) Run(ctx context.Context, savedFormulaID string, params *SavedFormulaRunParams) (*FormulaRun, error) {
	id, err := pathEscape("saved formula id", savedFormulaID)
	if err != nil {
		return nil, err
	}
	body := map[string]any{}
	var wait time.Duration
	var key string
	if params != nil {
		if params.Preview {
			body["mode"] = ModePreview
		}
		if params.IncludeSeen {
			body["include_seen"] = true
		}
		wait, key = params.Wait, params.IdempotencyKey
	}
	return startFormulaRun(ctx, s.client, "/v1/saved-formulas/"+id+"/runs", body, wait, key)
}

// Pause pauses a saved formula's schedule. It can still be run by hand. It is safe to repeat.
func (s *SavedFormulaService) Pause(ctx context.Context, savedFormulaID string) (*SavedFormula, error) {
	return s.call(ctx, http.MethodPost, savedFormulaID, "/pause", nil, true)
}

// Resume starts a paused saved formula's schedule again. It is safe to repeat.
func (s *SavedFormulaService) Resume(ctx context.Context, savedFormulaID string) (*SavedFormula, error) {
	return s.call(ctx, http.MethodPost, savedFormulaID, "/resume", nil, true)
}

// ResetSeen forgets what a saved formula delivered, so that its next run delivers rows already delivered again.
// Owners and admins only.
func (s *SavedFormulaService) ResetSeen(ctx context.Context, savedFormulaID string) (*SeenReset, error) {
	id, err := pathEscape("saved formula id", savedFormulaID)
	if err != nil {
		return nil, err
	}
	var out SeenReset
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/saved-formulas/" + id + "/seen/reset"}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SavedFormulaRowListParams narrow and page [SavedFormulaService.Rows] and [SavedFormulaService.RowsPage].
type SavedFormulaRowListParams struct {
	// Since keeps the rows first seen at or after this time: what a reader that polls asks for.
	Since time.Time
	// Order is "asc" (the default: in the order they were written) or "desc" (newest first).
	Order string
	// Fields keeps only these columns in each row's Data.
	Fields []string
	// MinScore, when set, keeps the rows scored at least this.
	MinScore *float64
	// Limit is the rows per page: 1 to 1000. Zero leaves it to the API (100).
	Limit int
	// Cursor is the NextCursor of the page before.
	Cursor string
}

func (p *SavedFormulaRowListParams) query() url.Values {
	if p == nil {
		return url.Values{}
	}
	q := rowQuery(p.Fields, p.MinScore, p.Limit, p.Cursor)
	setSince(q, p.Since)
	setString(q, "order", p.Order)
	return q
}

// setSince sends a time as the API writes one, to the microsecond: 2026-11-25T09:02:55.250000+00:00.
func setSince(q url.Values, since time.Time) {
	if !since.IsZero() {
		q.Set("since", since.UTC().Format("2006-01-02T15:04:05.000000-07:00"))
	}
}

// RowsPage returns one page of a saved formula's rows, across its runs. Rows the suppression list now holds are
// left out, and so are a preview's.
func (s *SavedFormulaService) RowsPage(ctx context.Context, savedFormulaID string, params *SavedFormulaRowListParams) (*FormulaRowPage, error) {
	id, err := pathEscape("saved formula id", savedFormulaID)
	if err != nil {
		return nil, err
	}
	var out FormulaRowPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/saved-formulas/" + id + "/rows", query: params.query(), retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Rows iterates over a saved formula's rows across its runs, fetching page after page. A reader that polls keeps
// the time of its last look and asks with Since.
func (s *SavedFormulaService) Rows(ctx context.Context, savedFormulaID string, params *SavedFormulaRowListParams) *Iter[FormulaRow] {
	p := SavedFormulaRowListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]FormulaRow, string, error) {
		p.Cursor = cursor
		page, err := s.RowsPage(ctx, savedFormulaID, &p)
		if err != nil {
			return nil, "", err
		}
		return page.Rows, page.NextCursor, nil
	})
}

// SavedFormulaExportParams are the options of [SavedFormulaService.Export] and [SavedFormulaService.ExportAll].
type SavedFormulaExportParams struct {
	// Format is [FormatCSV] or [FormatNDJSON].
	Format string
	// Since, Order, Fields and MinScore are as in [SavedFormulaRowListParams].
	Since    time.Time
	Order    string
	Fields   []string
	MinScore *float64
	// Limit is the rows per page: 1 to 1000. Zero leaves it to the API (100), and ExportAll reads 1000 a page.
	Limit int
	// Cursor is the NextCursor of the page before; where ExportAll starts.
	Cursor string
}

// Export returns one page of a saved formula's rows written by the API, as CSV or NDJSON, as
// [FormulaRunService.Export] does for a run.
func (s *SavedFormulaService) Export(ctx context.Context, savedFormulaID string, params SavedFormulaExportParams) (*RowsExport, error) {
	id, err := pathEscape("saved formula id", savedFormulaID)
	if err != nil {
		return nil, err
	}
	q := (&SavedFormulaRowListParams{Since: params.Since, Order: params.Order, Fields: params.Fields, MinScore: params.MinScore,
		Limit: params.Limit, Cursor: params.Cursor}).query()
	return exportRows(ctx, s.client, "/v1/saved-formulas/"+id+"/rows", params.Format, q)
}

// ExportAll writes every row of a saved formula to w, as CSV or NDJSON, page after page: the result is one file
// (a CSV has its header row once).
func (s *SavedFormulaService) ExportAll(ctx context.Context, savedFormulaID string, params SavedFormulaExportParams, w io.Writer) error {
	if params.Limit == 0 {
		params.Limit = 1000
	}
	return exportAll(params.Format, params.Cursor, w, func(cursor string) (*RowsExport, error) {
		params.Cursor = cursor
		return s.Export(ctx, savedFormulaID, params)
	})
}

// SnapshotListParams narrow [SavedFormulaService.Snapshots].
type SnapshotListParams struct {
	// EntityKey keeps one entity: its history.
	EntityKey string
	// Limit is how many, 1 to 1000. Zero leaves it to the API (200).
	Limit int
}

// Snapshots returns what a saved formula that watches entities (prices, stock, sellers) observed: one entry each
// time an entity's tracked values changed, newest first. They are kept 13 months.
func (s *SavedFormulaService) Snapshots(ctx context.Context, savedFormulaID string, params *SnapshotListParams) ([]Snapshot, error) {
	id, err := pathEscape("saved formula id", savedFormulaID)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	if params != nil {
		setString(q, "entity_key", params.EntityKey)
		setInt(q, "limit", params.Limit)
	}
	var out struct {
		Snapshots []Snapshot `json:"snapshots"`
	}
	err = s.client.do(ctx, request{method: http.MethodGet, path: "/v1/saved-formulas/" + id + "/snapshots", query: q, retry: true}, &out)
	return out.Snapshots, err
}

// Limits returns what the account's plan allows of saved formulas: active schedules, the shortest interval, how
// long rows are kept.
func (s *SavedFormulaService) Limits(ctx context.Context) (*PlanLimits, error) {
	var out PlanLimits
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/saved-formulas/limits", retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SuppressionService is the account's suppression list: domains, email addresses and company names that formula
// runs leave out before any paid step, before rows are written and when rows are sent.
type SuppressionService struct{ client *Client }

// SuppressionListParams filter and page [SuppressionService.List].
type SuppressionListParams struct {
	// Kind keeps one kind: [SuppressDomain], [SuppressEmail] or [SuppressCompanyName].
	Kind string
	// Limit is the entries per page: 1 to 1000. Zero leaves it to the API (100).
	Limit int
	// Cursor is the NextCursor of the page before.
	Cursor string
}

// List returns one page of the suppression list.
func (s *SuppressionService) List(ctx context.Context, params *SuppressionListParams) (*SuppressionPage, error) {
	q := url.Values{}
	if params != nil {
		setString(q, "kind", params.Kind)
		setInt(q, "limit", params.Limit)
		setString(q, "cursor", params.Cursor)
	}
	var out SuppressionPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/suppressions", query: q, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAll iterates over the suppression list, fetching page after page.
func (s *SuppressionService) ListAll(ctx context.Context, params *SuppressionListParams) *Iter[Suppression] {
	p := SuppressionListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]Suppression, string, error) {
		p.Cursor = cursor
		page, err := s.List(ctx, &p)
		if err != nil {
			return nil, "", err
		}
		return page.Suppressions, page.NextCursor, nil
	})
}

// MaxSuppressionEntries is how many entries [SuppressionService.Add] takes at once.
const MaxSuppressionEntries = 10000

// SuppressionAddParams are the options of [SuppressionService.Add].
type SuppressionAddParams struct {
	// Source is "manual" (the default) or "upload".
	Source string
}

// Add adds entries to the suppression list, and returns how many were new to it. Each value is kept in the form
// runs match it (a bare domain, a lower-case address, a company name without its legal form). It is safe to
// repeat: an entry the list holds is left as it is.
func (s *SuppressionService) Add(ctx context.Context, entries []SuppressionEntry, params *SuppressionAddParams) (int, error) {
	if len(entries) == 0 || len(entries) > MaxSuppressionEntries {
		return 0, fmt.Errorf("enrichednow: give 1 to %d entries", MaxSuppressionEntries)
	}
	body := map[string]any{"entries": entries}
	if params != nil && params.Source != "" {
		body["source"] = params.Source
	}
	var out struct {
		Added int `json:"added"`
	}
	err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/suppressions", body: body, retry: true}, &out)
	return out.Added, err
}

// SuppressionUploadParams are the options of [SuppressionService.Upload].
type SuppressionUploadParams struct {
	// Kind is the kind of every value of a file without a header row: [SuppressDomain], [SuppressEmail] or
	// [SuppressCompanyName]. Without it each value is read as an address when it holds @, a domain when it looks
	// like one, and a company name otherwise.
	Kind string
}

// Upload adds a CSV file to the suppression list: at most 2 MB and 50,000 lines. With a header row holding kind and
// value, each line says its kind; otherwise the first column is read. Lines that cannot be read come back in
// Invalid, and the rest are added.
//
// The file is read into memory before it is sent. It is sent again after a failure, as adding an entry twice
// changes nothing; the counts of a second attempt then say "already" for what the first one added.
func (s *SuppressionService) Upload(ctx context.Context, csv io.Reader, params *SuppressionUploadParams) (*SuppressionUpload, error) {
	data, err := io.ReadAll(csv)
	if err != nil {
		return nil, fmt.Errorf("enrichednow: reading the file: %w", err)
	}
	q := url.Values{}
	if params != nil {
		setString(q, "kind", params.Kind)
	}
	var out SuppressionUpload
	r := request{method: http.MethodPost, path: "/v1/suppressions/upload", query: q, payload: data, contentType: "text/csv", retry: true}
	if err := s.client.do(ctx, r, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes an entry from the suppression list. value is matched as it is kept (see [Suppression.Value]).
// An entry that is not on the list is [ErrNotFound].
func (s *SuppressionService) Delete(ctx context.Context, kind, value string) error {
	k, err := pathEscape("suppression kind", kind)
	if err != nil {
		return err
	}
	v, err := pathEscape("suppression value", value)
	if err != nil {
		return err
	}
	return s.client.do(ctx, request{method: http.MethodDelete, path: "/v1/suppressions/" + k + "/" + v, retry: true}, nil)
}

// AttestationService reads the terms some formulas need accepted. Accepting them is for an owner or an admin, in
// the console: a key cannot.
type AttestationService struct{ client *Client }

// List returns each kind of terms, its current text, and whether the account accepted that version. A formula whose
// [Formula.NeedsAttestation] is true runs only once its terms are accepted; until then it fails with
// [ErrAttestationRequired].
func (s *AttestationService) List(ctx context.Context) ([]Attestation, error) {
	var out struct {
		Attestations []Attestation `json:"attestations"`
	}
	err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/attestations", retry: true}, &out)
	return out.Attestations, err
}

// FormulaRowService takes feedback on the rows of formula runs.
type FormulaRowService struct{ client *Client }

// RowFeedbackParams are the options of [FormulaRowService.Feedback].
type RowFeedbackParams struct {
	// Note is a sentence of yours, up to 1000 characters.
	Note string
	// CorrectDomain is, with [VerdictWrongCompany], the company's real domain. Without it the company is left
	// unmatched from now on.
	CorrectDomain string
}

// Feedback says what is wrong with a row, or that it is good: [VerdictGood], [VerdictWrongCompany],
// [VerdictWrongContact] or [VerdictNotRelevant]. One answer a row; a later one replaces it, so it is safe to
// repeat. wrong_company also corrects the match for this account's runs from now on.
func (s *FormulaRowService) Feedback(ctx context.Context, rowID, verdict string, params *RowFeedbackParams) (*RowFeedback, error) {
	id, err := pathEscape("row id", rowID)
	if err != nil {
		return nil, err
	}
	body := map[string]any{"verdict": verdict}
	if params != nil {
		if params.Note != "" {
			body["note"] = params.Note
		}
		if params.CorrectDomain != "" {
			body["correct_domain"] = params.CorrectDomain
		}
	}
	var out RowFeedback
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/formula-rows/" + id + "/feedback", body: body, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
