package enrichednow

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DefaultFormulaTimeout is how long [FormulaRunService.Wait] and [Client.RunFormula] wait for a formula run to
// end, unless told otherwise.
const DefaultFormulaTimeout = 10 * time.Minute

// FormulaWait is how long each look at a run in [FormulaRunService.Wait] asks the API to hold the request, at most:
// the API answers the moment the run ends.
const FormulaWait = 60 * time.Second

// Ptr returns a pointer to a copy of v: what the optional fields of the parameter types take.
//
//	params := enrichednow.JobOrderHunterParams{JobTitles: []string{"Java Developer"}, Locations: []string{"London"},
//		PostedWithinDays: enrichednow.Ptr(2)}
func Ptr[T any](v T) *T { return &v }

// formulaParams checks the params given for a formula and returns what is sent: an empty object for none, and an
// error, before anything is sent, for a parameter type made for another formula.
func formulaParams(formulaID string, params any) (any, error) {
	if params == nil {
		return map[string]any{}, nil
	}
	if v := reflect.ValueOf(params); (v.Kind() == reflect.Pointer || v.Kind() == reflect.Map) && v.IsNil() {
		return map[string]any{}, nil
	}
	if typed, ok := params.(FormulaParams); ok && typed.FormulaID() != formulaID {
		return nil, fmt.Errorf("enrichednow: the parameters are those of the formula %q, not of %q", typed.FormulaID(), formulaID)
	}
	return params, nil
}

// FormulaService is the catalogue of formulas, and where their runs are started.
type FormulaService struct{ client *Client }

// FormulaListParams narrow [FormulaService.List]. The zero value lists every formula.
type FormulaListParams struct {
	// Query is what you want, in your own words ("companies hiring Java developers in London with HR contacts"):
	// the formulas that match come best first.
	Query string
	// Category keeps the formulas of one category, such as "recruiting".
	Category string
	// Persona keeps the formulas made for one persona, such as "recruiter".
	Persona string
	// Status keeps the formulas with this status: beta or live.
	Status string
}

// List returns the catalogue of formulas. It needs no API key, and sends the key when the client has one: then
// [Formula.Available] says whether each can run for this account.
func (s *FormulaService) List(ctx context.Context, params *FormulaListParams) ([]Formula, error) {
	q := url.Values{}
	if params != nil {
		setString(q, "query", params.Query)
		setString(q, "category", params.Category)
		setString(q, "persona", params.Persona)
		setString(q, "status", params.Status)
	}
	var out struct {
		Formulas []Formula `json:"formulas"`
	}
	err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/formulas", query: q, public: true, retry: true}, &out)
	return out.Formulas, err
}

// Get returns one formula: what it does, the JSON Schema of its params, its presets, steps and columns, what it
// needs, its limits and sample rows. It needs no API key.
func (s *FormulaService) Get(ctx context.Context, formulaID string) (*FormulaDetails, error) {
	id, err := pathEscape("formula id", formulaID)
	if err != nil {
		return nil, err
	}
	var out FormulaDetails
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/formulas/" + id, public: true, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FormulaEstimateParams are the options of [FormulaService.Estimate].
type FormulaEstimateParams struct {
	// Preview estimates a preview run: at most 10 rows and 50 units.
	Preview bool
	// UnitCap is the most the run would be allowed to spend; zero: the formula's default.
	UnitCap int
	// RowCap is the most rows it would be allowed to produce; zero: the formula's default.
	RowCap int
}

// Estimate says about what a run with these params would cost, and how many rows it would give. It spends
// nothing, and is safe to repeat.
//
// params is a parameter type of this package (such as [JobOrderHunterParams]), a struct of yours or a map; nil
// is no params. A parameter type made for another formula is refused before anything is sent.
//
// Params that cannot run are not an error: the answer's Valid is false and its Errors say what to change.
func (s *FormulaService) Estimate(ctx context.Context, formulaID string, params any, opts *FormulaEstimateParams) (*FormulaEstimate, error) {
	id, err := pathEscape("formula id", formulaID)
	if err != nil {
		return nil, err
	}
	sent, err := formulaParams(formulaID, params)
	if err != nil {
		return nil, err
	}
	body := map[string]any{"params": sent}
	if opts != nil {
		setCaps(body, opts.Preview, opts.UnitCap, opts.RowCap)
	}
	var out FormulaEstimate
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/formulas/" + id + "/estimate", body: body, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// setCaps adds what is set of a run's mode and caps to a request body: the API's defaults are not sent.
func setCaps(body map[string]any, preview bool, unitCap, rowCap int) {
	if preview {
		body["mode"] = ModePreview
	}
	if unitCap > 0 {
		body["unit_cap"] = unitCap
	}
	if rowCap > 0 {
		body["row_cap"] = rowCap
	}
}

// FormulaRunParams are the options of [FormulaService.Run].
type FormulaRunParams struct {
	// Preview runs at most 10 rows and 50 units, charged as usual; nothing it returns is remembered as delivered.
	Preview bool
	// UnitCap is the most the run may spend; zero: the formula's default. It cannot be overspent: the run stops at
	// the cap, keeps what it has and ends partial.
	UnitCap int
	// RowCap is the most rows it may produce; zero: its max_results param, or the formula's default.
	RowCap int
	// Wait asks the API to hold the request open until the run ends or this much time has passed, whichever
	// comes first: from 0 to [MaxWait]. A run that has not ended by then is returned queued or running; follow it
	// with [FormulaRunService.Wait].
	Wait time.Duration
	// IdempotencyKey names this call so that it can be repeated without starting a second run, as
	// [RunParams.IdempotencyKey] does for a scraper run. Empty: a new key for each call.
	IdempotencyKey string
}

// Run starts a run of a formula and returns it: usually queued, or already ended when opts.Wait was long enough.
//
// params is a parameter type of this package (such as [JobOrderHunterParams]), a struct of yours or a map; nil
// is no params. Only what is set is sent, and the formula's defaults fill the rest. A parameter type made for
// another formula is refused before anything is sent.
//
// The request carries an Idempotency-Key, so it is safe to send again: after a connection error or a 5xx it is
// retried with the same key, and the API answers with the run the first attempt started.
//
// A run that ended failed, unsupported or cancelled is returned without an error, as the API answers it; call
// [FormulaRun.Err] to turn its status into one.
func (s *FormulaService) Run(ctx context.Context, formulaID string, params any, opts *FormulaRunParams) (*FormulaRun, error) {
	id, err := pathEscape("formula id", formulaID)
	if err != nil {
		return nil, err
	}
	sent, err := formulaParams(formulaID, params)
	if err != nil {
		return nil, err
	}
	body := map[string]any{"params": sent}
	var wait time.Duration
	var key string
	if opts != nil {
		setCaps(body, opts.Preview, opts.UnitCap, opts.RowCap)
		wait, key = opts.Wait, opts.IdempotencyKey
	}
	return startFormulaRun(ctx, s.client, "/v1/formulas/"+id+"/runs", body, wait, key)
}

// startFormulaRun sends the request that starts a formula run, with an idempotency key and the wait asked for.
func startFormulaRun(ctx context.Context, c *Client, path string, body any, wait time.Duration, key string) (*FormulaRun, error) {
	if wait < 0 || wait > MaxWait {
		return nil, fmt.Errorf("enrichednow: Wait must be from 0 to %s", MaxWait)
	}
	if key == "" {
		key = newIdempotencyKey()
	}
	r := request{method: http.MethodPost, path: path, body: body, retry: true, idempotencyKey: key}
	if wait > 0 {
		r.query = url.Values{"wait": {seconds(wait)}}
		r.extraTimeout = wait
	}
	var out FormulaRun
	if err := c.do(ctx, r, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func seconds(d time.Duration) string { return strconv.FormatFloat(d.Seconds(), 'f', -1, 64) }

// FormulaRunService reads, waits for and cancels formula runs, and reads their rows.
type FormulaRunService struct{ client *Client }

// FormulaRunListParams filter and page [FormulaRunService.List].
type FormulaRunListParams struct {
	// Limit is the runs per page: 1 to 100. Zero leaves it to the API (20).
	Limit int
	// Cursor is the NextCursor of the page before.
	Cursor string
	// FormulaID keeps only the runs of this formula.
	FormulaID string
	// SavedFormulaID keeps only the runs of this saved formula.
	SavedFormulaID string
	// Status keeps only the runs with this status.
	Status RunStatus
}

// List returns one page of the account's formula runs, newest first. A run in a list has no Steps.
func (s *FormulaRunService) List(ctx context.Context, params *FormulaRunListParams) (*FormulaRunPage, error) {
	q := url.Values{}
	if params != nil {
		setInt(q, "limit", params.Limit)
		setString(q, "cursor", params.Cursor)
		setString(q, "formula_id", params.FormulaID)
		setString(q, "saved_formula_id", params.SavedFormulaID)
		setString(q, "status", string(params.Status))
	}
	var out FormulaRunPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/formula-runs", query: q, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAll iterates over the account's formula runs, newest first, fetching page after page.
func (s *FormulaRunService) ListAll(ctx context.Context, params *FormulaRunListParams) *Iter[FormulaRun] {
	p := FormulaRunListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]FormulaRun, string, error) {
		p.Cursor = cursor
		page, err := s.List(ctx, &p)
		if err != nil {
			return nil, "", err
		}
		return page.Runs, page.NextCursor, nil
	})
}

// Get returns a formula run as it is now, with what each step did.
func (s *FormulaRunService) Get(ctx context.Context, runID string) (*FormulaRun, error) {
	return s.GetWait(ctx, runID, 0)
}

// GetWait returns a formula run when it has ended, or as it is once wait (at most [MaxWait]) has passed: the API
// holds the request open in the meantime.
func (s *FormulaRunService) GetWait(ctx context.Context, runID string, wait time.Duration) (*FormulaRun, error) {
	id, err := pathEscape("formula run id", runID)
	if err != nil {
		return nil, err
	}
	if wait < 0 || wait > MaxWait {
		return nil, fmt.Errorf("enrichednow: wait must be from 0 to %s", MaxWait)
	}
	r := request{method: http.MethodGet, path: "/v1/formula-runs/" + id, retry: true}
	if wait > 0 {
		r.query = url.Values{"wait": {seconds(wait)}}
		r.extraTimeout = wait
	}
	var out FormulaRun
	if err := s.client.do(ctx, r, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FormulaWaitParams are the options of [FormulaRunService.Wait].
type FormulaWaitParams struct {
	// Timeout is how long to wait for the run to end; zero means [DefaultFormulaTimeout]. A deadline of the
	// context that comes sooner ends the wait too.
	Timeout time.Duration
	// PollInterval is the least time between two looks at the run; zero means [DefaultPollInterval].
	PollInterval time.Duration
	// OnPoll, when set, is called with the run every time it was read.
	OnPoll func(*FormulaRun)
}

// Wait waits until a formula run has ended, however it ended, and returns it. Each look asks the API to hold the
// request until the run ends (up to [FormulaWait]), and two looks are at least params.PollInterval apart.
//
// When the time is up the run goes on, and Wait returns it as last seen with an error that matches
// context.DeadlineExceeded ([errors.Is]). A run that ended failed, unsupported or cancelled is returned without an
// error; call [FormulaRun.Err] to turn its status into one.
func (s *FormulaRunService) Wait(ctx context.Context, runID string, params *FormulaWaitParams) (*FormulaRun, error) {
	timeout, interval := DefaultFormulaTimeout, DefaultPollInterval
	var onPoll func(*FormulaRun)
	if params != nil {
		if params.Timeout > 0 {
			timeout = params.Timeout
		}
		if params.PollInterval > 0 {
			interval = params.PollInterval
		}
		onPoll = params.OnPoll
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var last *FormulaRun
	for {
		started := time.Now()
		wait := FormulaWait
		if deadline, ok := ctx.Deadline(); ok {
			// Whole seconds, and less than what is left, so that the answer comes before the time is up.
			wait = max(0, min(wait, time.Until(deadline).Truncate(time.Second)))
		}
		run, err := s.GetWait(ctx, runID, wait)
		if err != nil {
			if ctx.Err() != nil {
				return last, ctx.Err()
			}
			return last, err
		}
		last = run
		if onPoll != nil {
			onPoll(run)
		}
		if run.Status.IsFinal() {
			return run, nil
		}
		// A server that answered at once is not asked again before the interval is over.
		if rest := interval - time.Since(started); rest > 0 {
			if err := s.client.sleep(ctx, rest); err != nil {
				return run, err
			}
		}
	}
}

// Cancel asks a formula run to stop and returns it as it is at that moment: usually still running, as calls in
// flight are asked to stop. Rows already written stay, and what was already done stays charged. A run that has
// ended is returned unchanged. It is safe to repeat.
func (s *FormulaRunService) Cancel(ctx context.Context, runID string) (*FormulaRun, error) {
	id, err := pathEscape("formula run id", runID)
	if err != nil {
		return nil, err
	}
	var out FormulaRun
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/formula-runs/" + id + "/cancel", retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FormulaRowListParams narrow and page [FormulaRunService.Rows] and [FormulaRunService.RowsPage].
type FormulaRowListParams struct {
	// Fields keeps only these columns in each row's Data.
	Fields []string
	// MinScore, when set, keeps the rows scored at least this.
	MinScore *float64
	// View is [ViewRows] (the default) or [ViewCompany]: one row per company, for a formula that has a company
	// view ([FormulaDetails.CompanyView]).
	View string
	// Limit is the rows per page: 1 to 1000. Zero leaves it to the API (100).
	Limit int
	// Cursor is the NextCursor of the page before.
	Cursor string
}

func rowQuery(fields []string, minScore *float64, limit int, cursor string) url.Values {
	q := url.Values{}
	if len(fields) > 0 {
		q.Set("fields", strings.Join(fields, ","))
	}
	if minScore != nil {
		q.Set("min_score", strconv.FormatFloat(*minScore, 'f', -1, 64))
	}
	setInt(q, "limit", limit)
	setString(q, "cursor", cursor)
	return q
}

// RowsPage returns one page of a formula run's rows, best first. Rows can be read while the run is still going.
func (s *FormulaRunService) RowsPage(ctx context.Context, runID string, params *FormulaRowListParams) (*FormulaRowPage, error) {
	id, err := pathEscape("formula run id", runID)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	if params != nil {
		q = rowQuery(params.Fields, params.MinScore, params.Limit, params.Cursor)
		setString(q, "view", params.View)
	}
	var out FormulaRowPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/formula-runs/" + id + "/rows", query: q, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Rows iterates over a formula run's rows, best first, fetching page after page. Rows are kept 90 days, or as
// long as the plan keeps them.
func (s *FormulaRunService) Rows(ctx context.Context, runID string, params *FormulaRowListParams) *Iter[FormulaRow] {
	p := FormulaRowListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]FormulaRow, string, error) {
		p.Cursor = cursor
		page, err := s.RowsPage(ctx, runID, &p)
		if err != nil {
			return nil, "", err
		}
		return page.Rows, page.NextCursor, nil
	})
}

// FormulaExportParams are the options of [FormulaRunService.Export] and [FormulaRunService.ExportAll].
type FormulaExportParams struct {
	// Format is [FormatCSV] or [FormatNDJSON].
	Format string
	// Fields keeps only these columns, in this order.
	Fields []string
	// MinScore, when set, keeps the rows scored at least this.
	MinScore *float64
	// View is [ViewRows] (the default) or [ViewCompany].
	View string
	// Limit is the rows per page: 1 to 1000. Zero leaves it to the API (100), and ExportAll reads 1000 a page.
	Limit int
	// Cursor is the NextCursor of the page before; where ExportAll starts.
	Cursor string
}

// Export returns one page of a formula run's rows written by the API, as CSV or NDJSON.
//
// The CSV has a column for each of the formula's columns, in its order (or for each of Fields), and every page
// begins with the header row; a cell a spreadsheet would run as a formula starts with an apostrophe. An NDJSON
// line is one row's columns. NextCursor continues; [FormulaRunService.ExportAll] follows it for you.
func (s *FormulaRunService) Export(ctx context.Context, runID string, params FormulaExportParams) (*RowsExport, error) {
	id, err := pathEscape("formula run id", runID)
	if err != nil {
		return nil, err
	}
	q := rowQuery(params.Fields, params.MinScore, params.Limit, params.Cursor)
	setString(q, "view", params.View)
	return exportRows(ctx, s.client, "/v1/formula-runs/"+id+"/rows", params.Format, q)
}

// ExportAll writes every row of a formula run to w, as CSV or NDJSON, page after page: the result is one file (a
// CSV has its header row once).
func (s *FormulaRunService) ExportAll(ctx context.Context, runID string, params FormulaExportParams, w io.Writer) error {
	if params.Limit == 0 {
		params.Limit = 1000
	}
	return exportAll(params.Format, params.Cursor, w, func(cursor string) (*RowsExport, error) {
		params.Cursor = cursor
		return s.Export(ctx, runID, params)
	})
}

// exportRows asks for one page of rows written as CSV or NDJSON.
func exportRows(ctx context.Context, c *Client, path, format string, q url.Values) (*RowsExport, error) {
	accept := map[string]string{FormatCSV: "text/csv", FormatNDJSON: "application/x-ndjson"}[format]
	if accept == "" {
		return nil, fmt.Errorf("enrichednow: the format %q is not csv or ndjson; for JSON use RowsPage", format)
	}
	q.Set("format", format)
	got, err := c.send(ctx, request{method: http.MethodGet, path: path, query: q, retry: true, accept: accept})
	if err != nil {
		return nil, err
	}
	return &RowsExport{Format: format, Body: got.data, NextCursor: got.header.Get("X-Next-Cursor")}, nil
}

// exportAll writes page after page to w, the header row of a CSV once.
func exportAll(format, cursor string, w io.Writer, page func(cursor string) (*RowsExport, error)) error {
	for first := true; ; first = false {
		got, err := page(cursor)
		if err != nil {
			return err
		}
		body := got.Body
		if !first && format == FormatCSV {
			if end := bytes.IndexByte(body, '\n'); end >= 0 {
				body = body[end+1:]
			} else {
				body = nil
			}
		}
		if _, err := w.Write(body); err != nil {
			return err
		}
		if got.NextCursor == "" {
			return nil
		}
		cursor = got.NextCursor
	}
}

// RunFormulaParams are the options of [Client.RunFormula].
type RunFormulaParams struct {
	// Preview runs at most 10 rows and 50 units; see [FormulaRunParams.Preview].
	Preview bool
	// UnitCap is the most the run may spend; zero: the formula's default.
	UnitCap int
	// RowCap is the most rows it may produce; zero: the formula's default.
	RowCap int
	// Timeout is how long to wait for the run to end; zero means [DefaultFormulaTimeout]. When it passes the run
	// goes on, and RunFormula returns it with an error that matches context.DeadlineExceeded.
	Timeout time.Duration
	// PollInterval is the least time between two looks at the run; zero means [DefaultPollInterval].
	PollInterval time.Duration
	// IdempotencyKey names this start; see [FormulaRunParams.IdempotencyKey].
	IdempotencyKey string
	// AllowFailed returns a run that ended failed, unsupported or cancelled without an error, and with an iterator
	// over the rows it wrote (a cancelled run keeps those written before it stopped). Without it such a run comes
	// back with a [*FormulaRunFailedError].
	AllowFailed bool
}

// RunFormula runs a formula to the end: it starts a run, waits until it has ended and returns it with an iterator
// over its rows, best first, read 1000 a page.
//
//	run, rows, err := client.RunFormula(ctx, "job-order-hunter", enrichednow.JobOrderHunterParams{
//		JobTitles: []string{"Java Developer"},
//		Locations: []string{"London"},
//	}, nil)
//
// A run that ended failed, unsupported or cancelled is returned with a [*FormulaRunFailedError] (it matches
// [ErrRunFailed]), unless opts.AllowFailed; partial is not an error: it has rows, and its Reason says what was
// missed. To narrow the rows, use [FormulaRunService.Rows] with the run's ID; to download them as one file,
// [FormulaRunService.ExportAll].
func (c *Client) RunFormula(ctx context.Context, formulaID string, params any, opts *RunFormulaParams) (*FormulaRun, *Iter[FormulaRow], error) {
	o := RunFormulaParams{}
	if opts != nil {
		o = *opts
	}
	timeout := o.Timeout
	if timeout <= 0 {
		timeout = DefaultFormulaTimeout
	}
	deadline := time.Now().Add(timeout)
	run, err := c.Formulas.Run(ctx, formulaID, params, &FormulaRunParams{Preview: o.Preview, UnitCap: o.UnitCap, RowCap: o.RowCap,
		Wait: min(ScrapeWait, timeout), IdempotencyKey: o.IdempotencyKey})
	if err != nil {
		return nil, nil, err
	}
	if !run.Status.IsFinal() {
		left := time.Until(deadline)
		if left <= 0 {
			return run, nil, context.DeadlineExceeded
		}
		latest, err := c.FormulaRuns.Wait(ctx, run.ID, &FormulaWaitParams{Timeout: left, PollInterval: o.PollInterval})
		if latest != nil {
			run = latest
		}
		if err != nil {
			return run, nil, err
		}
	}
	if err := run.Err(); err != nil && !o.AllowFailed {
		return run, nil, err
	}
	return run, c.FormulaRuns.Rows(ctx, run.ID, &FormulaRowListParams{Limit: 1000}), nil
}
