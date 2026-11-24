package enrichednow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

// fixture is one exchange recorded from the real API: testdata/formulas/<name>.json.
type fixture struct {
	Request struct {
		Method  string            `json:"method"`
		Path    string            `json:"path"`
		Query   map[string]string `json:"query"`
		Headers map[string]string `json:"headers"`
		Body    json.RawMessage   `json:"body"`
		Text    *string           `json:"text"`
	} `json:"request"`
	Response struct {
		Status  int               `json:"status"`
		Headers map[string]string `json:"headers"`
		Body    json.RawMessage   `json:"body"`
		Text    *string           `json:"text"`
	} `json:"response"`
}

func loadFixture(t *testing.T, name string) fixture {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "formulas", name+".json"))
	if err != nil {
		t.Fatal(err)
	}
	var fx fixture
	if err := json.Unmarshal(data, &fx); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return fx
}

// write answers a request with the fixture's recorded response.
func (fx fixture) write(w http.ResponseWriter) {
	for name, value := range fx.Response.Headers {
		w.Header().Set(name, value)
	}
	w.WriteHeader(fx.Response.Status)
	switch {
	case fx.Response.Text != nil:
		io.WriteString(w, *fx.Response.Text)
	case len(fx.Response.Body) > 0:
		w.Write(fx.Response.Body)
	}
}

// sentRequest is a request the fake API was sent.
type sentRequest struct {
	Method, Path string
	Query        url.Values
	Header       http.Header
	Body         []byte
}

// fixtureAPI is a fake API that answers with recorded responses and keeps what it was sent.
type fixtureAPI struct {
	mu   sync.Mutex
	sent []sentRequest
}

func (f *fixtureAPI) requests() []sentRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]sentRequest(nil), f.sent...)
}

// routed serves the fixture that pick names for each request, and returns a client that talks to it and never
// really sleeps (the pauses it asked for are in rec.delays).
func routed(t *testing.T, pick func(r *http.Request, n int) string, opts ...Option) (*Client, *fixtureAPI, *recorder) {
	t.Helper()
	api := &fixtureAPI{}
	cache := map[string]fixture{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		api.mu.Lock()
		api.sent = append(api.sent, sentRequest{Method: r.Method, Path: r.URL.Path, Query: r.URL.Query(), Header: r.Header.Clone(), Body: body})
		n := len(api.sent)
		api.mu.Unlock()
		name := pick(r, n)
		fx, ok := cache[name]
		if !ok {
			fx = loadFixture(t, name)
			cache[name] = fx
		}
		fx.write(w)
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(append([]Option{WithAPIKey(testKey), WithBaseURL(server.URL)}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	rec := &recorder{}
	client.sleep = func(ctx context.Context, d time.Duration) error {
		rec.mu.Lock()
		rec.delays = append(rec.delays, d)
		rec.mu.Unlock()
		return ctx.Err()
	}
	client.jitter = func() float64 { return 1 }
	return client, api, rec
}

// replay answers every request with the named fixtures, one after the other (the last one again after that).
func replay(t *testing.T, names ...string) (*Client, *fixtureAPI) {
	t.Helper()
	client, api, _ := routed(t, func(_ *http.Request, n int) string { return names[min(n, len(names))-1] })
	return client, api
}

// checkSent fails the test unless a request is the fixture's: the method, the path, the query (as a map, without
// the names in loose), the JSON body or the text, and the Idempotency-Key and Content-Type the fixture names.
func checkSent(t *testing.T, got sentRequest, fx fixture, loose ...string) {
	t.Helper()
	want := fx.Request
	if got.Method != want.Method || got.Path != want.Path {
		t.Errorf("sent %s %s, want %s %s", got.Method, got.Path, want.Method, want.Path)
	}
	query := map[string]string{}
	for name, values := range got.Query {
		query[name] = values[0]
	}
	wanted := map[string]string{}
	for name, value := range want.Query {
		wanted[name] = value
	}
	for _, name := range loose {
		delete(query, name)
		delete(wanted, name)
	}
	if !reflect.DeepEqual(query, wanted) {
		t.Errorf("%s %s: query %v, want %v", got.Method, got.Path, query, wanted)
	}
	switch {
	case want.Text != nil:
		if string(got.Body) != *want.Text {
			t.Errorf("%s %s: body %q, want %q", got.Method, got.Path, got.Body, *want.Text)
		}
	case len(want.Body) > 0:
		var a, b any
		if err := json.Unmarshal(got.Body, &a); err != nil {
			t.Fatalf("%s %s: the body is not JSON: %q", got.Method, got.Path, got.Body)
		}
		_ = json.Unmarshal(want.Body, &b)
		if !reflect.DeepEqual(a, b) {
			t.Errorf("%s %s: body\n %s\nwant\n %s", got.Method, got.Path, got.Body, want.Body)
		}
		if ct := got.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("%s %s: Content-Type %q", got.Method, got.Path, ct)
		}
	default:
		if len(got.Body) != 0 {
			t.Errorf("%s %s: no body expected, got %q", got.Method, got.Path, got.Body)
		}
	}
	for name, value := range want.Headers {
		if got.Header.Get(name) != value && !slicesContain(loose, name) {
			t.Errorf("%s %s: %s = %q, want %q", got.Method, got.Path, name, got.Header.Get(name), value)
		}
	}
	if got.Header.Get("Authorization") != "Bearer "+testKey {
		t.Errorf("%s %s: Authorization %q", got.Method, got.Path, got.Header.Get("Authorization"))
	}
}

func slicesContain(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

// jobOrder is the params of the recorded runs of job-order-hunter, as its parameter type.
var jobOrder = JobOrderHunterParams{
	JobTitles:            []string{"Java Developer"},
	Locations:            []string{"London"},
	PostedWithinDays:     Ptr(2),
	TitleExcludeKeywords: []string{"intern"},
	ExcludeCompanies:     []string{"Blocked By Us Inc"},
	CompanySize:          &IntRange{Min: 11, Max: 5000},
}

const (
	run1 = "frun_000000000000000000000001"
	run2 = "frun_000000000000000000000002"
	sf1  = "sf_000000000000000000000001"
	sf2  = "sf_000000000000000000000002"
)

// Every call: the request it sends is the one recorded from the real API, and the answer it parses is the
// recorded one.
func TestFormulaCallsSendTheRecordedRequests(t *testing.T) {
	ctx := context.Background()
	since := time.Date(2026, 11, 25, 9, 2, 55, 250_000_000, time.UTC)
	cases := []struct {
		fixture string
		call    func(c *Client) (any, error)
		check   func(t *testing.T, got any)
	}{
		{"formulas_list", func(c *Client) (any, error) { return c.Formulas.List(ctx, nil) }, func(t *testing.T, got any) {
			list := got.([]Formula)
			if len(list) != 3 || list[1].ID != "job-order-hunter" || !list[1].NeedsAttestation || list[1].Health == nil ||
				list[1].Health.Status != "ok" || list[1].Health.SuccessRate7d != nil || strings.Join(list[1].Schedules, ",") != "daily,weekdays,weekly" ||
				list[0].Available || list[1].Status != FormulaBeta || list[1].Personas[0] != "recruiter" {
				t.Errorf("formulas = %+v", list)
			}
			if b := list[1].Basics; len(b) != 3 || b[0].Name != "job_titles" || !b[0].Required || b[0].Type != "list of strings" ||
				b[0].Help == "" || b[2].Name != "posted_within_days" || b[2].Required {
				t.Errorf("basics = %+v", b)
			}
		}},
		{"formulas_list_filtered", func(c *Client) (any, error) {
			return c.Formulas.List(ctx, &FormulaListParams{Query: "hiring", Category: "recruiting"})
		}, func(t *testing.T, got any) {
			if list := got.([]Formula); len(list) != 1 || list[0].Category != "recruiting" {
				t.Errorf("formulas = %+v", list)
			}
		}},
		{"formulas_get", func(c *Client) (any, error) { return c.Formulas.Get(ctx, "job-order-hunter") }, func(t *testing.T, got any) {
			f := got.(*FormulaDetails)
			var schema struct {
				Required []string `json:"required"`
			}
			_ = json.Unmarshal(f.Parameters, &schema)
			if f.ID != "job-order-hunter" || f.Version != "1.1.0" || f.Title != "Job Order Hunter" || f.Description == "" ||
				strings.Join(schema.Required, ",") != "job_titles,locations" || len(f.Presets) != 5 || f.Presets[0].ID != "it-contract-uk" ||
				f.Steps[0].Label != "Find" || len(f.Columns) != 24 || !f.Columns[0].DefaultVisible || f.CompanyView == nil ||
				f.CompanyView.Key != "company_domain" || f.Requirements.NotOffered[0] != "google.jobs" || f.Limits["min_unit_cap"] != 20 ||
				len(f.SampleRows) == 0 || f.SampleIsFixture || len(f.Changelog) != 2 || len(f.RawJSON()) == 0 || !f.NeedsProvider {
				t.Errorf("formula = %+v", f)
			}
		}},
		{"formulas_estimate", func(c *Client) (any, error) {
			return c.Formulas.Estimate(ctx, "job-order-hunter", jobOrder, &FormulaEstimateParams{UnitCap: 2000})
		}, func(t *testing.T, got any) {
			e := got.(*FormulaEstimate)
			if !e.Valid || e.Estimate == nil || e.Estimate.ExpectedUnits != 63 || e.Estimate.MaxUnits != 1425 || e.Estimate.Basis == "" ||
				len(e.Errors) != 0 || e.NotOffered[0] != "google.jobs" {
				t.Errorf("estimate = %+v", e)
			}
		}},
		{"formulas_estimate_invalid", func(c *Client) (any, error) {
			return c.Formulas.Estimate(ctx, "job-order-hunter", map[string]any{"locations": []string{"London"}}, nil)
		}, func(t *testing.T, got any) {
			e := got.(*FormulaEstimate)
			if e.Valid || e.Estimate != nil || len(e.Errors) != 1 || e.Errors[0] != (FieldError{"params.job_titles", "This is required."}) {
				t.Errorf("estimate = %+v", e)
			}
		}},
		{"formulas_run", func(c *Client) (any, error) {
			return c.Formulas.Run(ctx, "job-order-hunter", jobOrder, &FormulaRunParams{Preview: true, Wait: 60 * time.Second, IdempotencyKey: "c0ffee-run-2"})
		}, func(t *testing.T, got any) {
			run := got.(*FormulaRun)
			params, err := Decode[JobOrderHunterParams](run.Params)
			if err != nil || run.ID != "frun_000000000000000000000003" || run.Status != RunPartial || run.Reason != "contacts_not_connected" ||
				run.Mode != ModePreview || run.Trigger != "api" || len(run.Reasons) != 2 || run.Reasons[1].Code != "source_not_offered" ||
				run.Reasons[1].Fields["scraper"] != "google.jobs" || run.Estimate.MaxUnits != 50 || run.UnitCap != 50 || run.RowCap != 10 ||
				run.UnitsSpent != 17 || run.RowsTotal != 4 || run.SavedFormulaID != "" || run.FinishedAt == nil || len(run.Steps) == 0 ||
				*params.PostedWithinDays != 2 || params.CompanySize.Max != 5000 || run.Err() != nil {
				t.Errorf("run = %+v (%v)", run, err)
			}
		}},
		{"formulas_run_queued", func(c *Client) (any, error) {
			return c.Formulas.Run(ctx, "job-order-hunter", &jobOrder, &FormulaRunParams{UnitCap: 2000, IdempotencyKey: "c0ffee-run-1"})
		}, func(t *testing.T, got any) {
			if run := got.(*FormulaRun); run.Status != RunQueued || run.Status.IsFinal() || run.StartedAt != nil || run.Err() != nil {
				t.Errorf("run = %+v", run)
			}
		}},
		{"formula_runs_get_running", func(c *Client) (any, error) { return c.FormulaRuns.Get(ctx, run2) }, func(t *testing.T, got any) {
			if run := got.(*FormulaRun); run.Status != RunRunning || run.FinishedAt != nil {
				t.Errorf("run = %+v", run)
			}
		}},
		{"formula_runs_get", func(c *Client) (any, error) { return c.FormulaRuns.GetWait(ctx, run2, 60*time.Second) }, func(t *testing.T, got any) {
			run := got.(*FormulaRun)
			var done *FormulaRunStep
			for i, step := range run.Steps {
				if step.ID == "search_linkedin" {
					done = &run.Steps[i]
				}
			}
			if run.Status != RunPartial || run.UnitsSpent != 24 || done == nil || done.Status != "done" || done.RecordsOut != 8 ||
				done.Calls != 1 || done.Units != 8 || done.StartedAt == nil || done.Errors == nil {
				t.Errorf("run = %+v, step %+v", run, done)
			}
		}},
		{"formula_runs_list", func(c *Client) (any, error) { return c.FormulaRuns.List(ctx, &FormulaRunListParams{Limit: 1}) }, func(t *testing.T, got any) {
			page := got.(*FormulaRunPage)
			if len(page.Runs) != 1 || page.NextCursor == "" || page.Runs[0].Steps != nil || page.Runs[0].Mode != ModePreview {
				t.Errorf("page = %+v", page)
			}
		}},
		{"formula_runs_rows", func(c *Client) (any, error) {
			return c.FormulaRuns.RowsPage(ctx, run2, &FormulaRowListParams{Limit: 2})
		}, func(t *testing.T, got any) {
			page := got.(*FormulaRowPage)
			row := page.Rows[0]
			data, _ := Decode[map[string]any](row.Data)
			evidence, _ := Decode[map[string]any](row.Evidence)
			if len(page.Rows) != 2 || page.NextCursor != "Mg" || row.RunID != run2 || *row.Score != 20 || data["company_name"] != "Corvid" ||
				evidence["sources"] == nil || row.Status != "new" || row.FirstSeenAt.IsZero() || row.RowKey == "" {
				t.Errorf("page = %+v", page)
			}
		}},
		{"formula_runs_rows_page2", func(c *Client) (any, error) {
			return c.FormulaRuns.RowsPage(ctx, run2, &FormulaRowListParams{Limit: 2, Cursor: "Mg"})
		}, func(t *testing.T, got any) {
			if page := got.(*FormulaRowPage); len(page.Rows) != 2 || page.NextCursor != "" {
				t.Errorf("page = %+v", page)
			}
		}},
		{"formula_runs_rows_fields", func(c *Client) (any, error) {
			return c.FormulaRuns.RowsPage(ctx, run2, &FormulaRowListParams{Fields: []string{"company_name", "job_title", "contact_email"}, MinScore: Ptr(10.0)})
		}, func(t *testing.T, got any) {
			page := got.(*FormulaRowPage)
			data, _ := Decode[map[string]any](page.Rows[0].Data)
			if len(data) != 3 || data["job_title"] == nil {
				t.Errorf("data = %v", data)
			}
		}},
		{"formula_runs_rows_company", func(c *Client) (any, error) {
			return c.FormulaRuns.RowsPage(ctx, run2, &FormulaRowListParams{View: ViewCompany})
		}, func(t *testing.T, got any) {
			page := got.(*FormulaRowPage)
			data, _ := Decode[map[string]any](page.Rows[0].Data)
			if len(page.Rows) != 3 || data["rows"] != float64(2) || len(data["job_title"].([]any)) != 2 {
				t.Errorf("page = %+v", page)
			}
		}},
		{"formula_runs_rows_csv", func(c *Client) (any, error) {
			return c.FormulaRuns.Export(ctx, run2, FormulaExportParams{Format: FormatCSV, Limit: 2})
		}, func(t *testing.T, got any) {
			page := got.(*RowsExport)
			if page.Format != FormatCSV || page.NextCursor != "Mg" || !bytes.HasPrefix(page.Body, []byte("company_name,company_domain,")) {
				t.Errorf("page = %+v", page)
			}
		}},
		{"formula_runs_rows_csv_page2", func(c *Client) (any, error) {
			return c.FormulaRuns.Export(ctx, run2, FormulaExportParams{Format: FormatCSV, Limit: 2, Cursor: "Mg"})
		}, func(t *testing.T, got any) {
			if page := got.(*RowsExport); page.NextCursor != "" || !bytes.Contains(page.Body, []byte("Northwind Ltd")) {
				t.Errorf("page = %+v", page)
			}
		}},
		{"formula_runs_rows_ndjson", func(c *Client) (any, error) {
			return c.FormulaRuns.Export(ctx, run2, FormulaExportParams{Format: FormatNDJSON})
		}, func(t *testing.T, got any) {
			if page := got.(*RowsExport); bytes.Count(page.Body, []byte("\n")) != 4 || page.NextCursor != "" {
				t.Errorf("page = %+v", page)
			}
		}},
		{"formula_runs_cancel", func(c *Client) (any, error) { return c.FormulaRuns.Cancel(ctx, run1) }, func(t *testing.T, got any) {
			if run := got.(*FormulaRun); run.ID != run1 || run.Status != RunRunning {
				t.Errorf("run = %+v", run)
			}
		}},
		{"formula_rows_feedback", func(c *Client) (any, error) {
			return c.FormulaRows.Feedback(ctx, "frow_000000000000000000000001", VerdictWrongCompany,
				&RowFeedbackParams{Note: "This is the agency, not the employer.", CorrectDomain: "corvid.example"})
		}, func(t *testing.T, got any) {
			f := got.(*RowFeedback)
			if f.Status != "flagged" || f.Verdict != VerdictWrongCompany || f.Correction == nil || f.Correction.Domain != "corvid.example" || f.Correction.Country != "UK" {
				t.Errorf("feedback = %+v", f)
			}
		}},
		{"saved_create", func(c *Client) (any, error) {
			return c.SavedFormulas.Create(ctx, "job-order-hunter", jobOrder, &SavedFormulaCreateParams{
				Schedule: &ScheduleParams{Cadence: CadenceWeekdays, LocalTime: "07:30", Timezone: "Europe/London"}, Notify: NotifyDigest})
		}, func(t *testing.T, got any) {
			s := got.(*SavedFormula)
			if s.ID != sf2 || !s.NameAuto || s.Status != SavedFormulaActive || s.Notify != NotifyDigest || s.Schedule == nil ||
				s.Schedule.Cadence != CadenceWeekdays || s.Schedule.NextFireAt == nil || s.Schedule.DayOfWeek != nil || s.MonthlyUnitCap != nil ||
				s.RunsWithKey == "" || s.Estimate.ExpectedRows != 6 || s.LastRun != nil || s.UnitCap != 3000 || s.Major != 1 || s.ArchivedAt != nil {
				t.Errorf("saved formula = %+v", s)
			}
		}},
		{"saved_list", func(c *Client) (any, error) { return c.SavedFormulas.List(ctx, &SavedFormulaListParams{Limit: 10}) }, func(t *testing.T, got any) {
			if page := got.(*SavedFormulaPage); len(page.SavedFormulas) != 1 || page.NextCursor != "" {
				t.Errorf("page = %+v", page)
			}
		}},
		{"saved_get", func(c *Client) (any, error) { return c.SavedFormulas.Get(ctx, sf2) }, func(t *testing.T, got any) {
			if s := got.(*SavedFormula); s.ID != sf2 || s.FormulaID != "job-order-hunter" {
				t.Errorf("saved formula = %+v", s)
			}
		}},
		{"saved_update", func(c *Client) (any, error) {
			return c.SavedFormulas.Update(ctx, sf2, SavedFormulaUpdateParams{Name: "Java in London, weekdays", MonthlyUnitCap: 20000})
		}, func(t *testing.T, got any) {
			if s := got.(*SavedFormula); s.NameAuto || *s.MonthlyUnitCap != 20000 {
				t.Errorf("saved formula = %+v", s)
			}
		}},
		{"saved_run", func(c *Client) (any, error) {
			return c.SavedFormulas.Run(ctx, sf2, &SavedFormulaRunParams{Wait: 60 * time.Second, IdempotencyKey: "c0ffee-saved-1"})
		}, func(t *testing.T, got any) {
			if run := got.(*FormulaRun); run.SavedFormulaID != sf2 || run.Status != RunPartial {
				t.Errorf("run = %+v", run)
			}
		}},
		{"saved_rows", func(c *Client) (any, error) {
			return c.SavedFormulas.RowsPage(ctx, sf2, &SavedFormulaRowListParams{Since: since.In(time.FixedZone("x", 3600)), Limit: 2})
		}, func(t *testing.T, got any) {
			if page := got.(*FormulaRowPage); len(page.Rows) != 2 || page.NextCursor == "" || page.Rows[0].RunID != "frun_000000000000000000000004" {
				t.Errorf("page = %+v", page)
			}
		}},
		{"saved_rows_csv", func(c *Client) (any, error) {
			return c.SavedFormulas.Export(ctx, sf2, SavedFormulaExportParams{Format: FormatCSV, Fields: []string{"company_name", "job_title", "job_url"}})
		}, func(t *testing.T, got any) {
			if page := got.(*RowsExport); !bytes.HasPrefix(page.Body, []byte("company_name,job_title,job_url\r\n")) || page.NextCursor != "" {
				t.Errorf("page = %+v", page)
			}
		}},
		{"saved_pause", func(c *Client) (any, error) { return c.SavedFormulas.Pause(ctx, sf2) }, func(t *testing.T, got any) {
			s := got.(*SavedFormula)
			if s.Status != SavedFormulaPaused || s.PausedReason != "manual" || s.Schedule.NextFireAt != nil || s.LastRun == nil ||
				s.LastRun.RowsNew != 4 || s.LastRun.Status != RunPartial {
				t.Errorf("saved formula = %+v", s)
			}
		}},
		{"saved_resume", func(c *Client) (any, error) { return c.SavedFormulas.Resume(ctx, sf2) }, func(t *testing.T, got any) {
			if s := got.(*SavedFormula); s.Status != SavedFormulaActive {
				t.Errorf("saved formula = %+v", s)
			}
		}},
		{"saved_seen_reset", func(c *Client) (any, error) { return c.SavedFormulas.ResetSeen(ctx, sf2) }, func(t *testing.T, got any) {
			if r := got.(*SeenReset); r.Cleared != 4 || r.SavedFormula.ID != sf2 {
				t.Errorf("reset = %+v", r)
			}
		}},
		{"saved_duplicate", func(c *Client) (any, error) {
			return c.SavedFormulas.Duplicate(ctx, sf2, &SavedFormulaDuplicateParams{Name: "Java in London (copy)", CopySeen: true})
		}, func(t *testing.T, got any) {
			if s := got.(*SavedFormula); s.ID != "sf_000000000000000000000003" || s.Name != "Java in London (copy)" || s.Schedule != nil {
				t.Errorf("saved formula = %+v", s)
			}
		}},
		{"saved_limits", func(c *Client) (any, error) { return c.SavedFormulas.Limits(ctx) }, func(t *testing.T, got any) {
			if l := got.(*PlanLimits); l.Schedules != 3 || l.SchedulesActive != 1 || l.MinIntervalMinutes != 1440 || l.RowRetentionDays != 90 || l.Destinations != 3 {
				t.Errorf("limits = %+v", l)
			}
		}},
		{"saved_archive", func(c *Client) (any, error) { return c.SavedFormulas.Archive(ctx, sf2) }, func(t *testing.T, got any) {
			if s := got.(*SavedFormula); s.Status != SavedFormulaArchived || s.ArchivedAt == nil {
				t.Errorf("saved formula = %+v", s)
			}
		}},
		{"saved_snapshots", func(c *Client) (any, error) {
			return c.SavedFormulas.Snapshots(ctx, sf1, &SnapshotListParams{EntityKey: "https://board.example/1"})
		}, func(t *testing.T, got any) {
			list := got.([]Snapshot)
			if len(list) != 2 || list[0].Fields["job.title"] != "Senior Java Developer" || list[0].RunID == "" || list[0].Entity != "job" {
				t.Errorf("snapshots = %+v", list)
			}
		}},
		{"suppressions_add", func(c *Client) (any, error) {
			return c.Suppressions.Add(ctx, []SuppressionEntry{{SuppressDomain, "https://www.globex.example/"}, {SuppressEmail, "Ana@Initech.example"},
				{SuppressCompanyName, "Umbrella, Inc."}}, nil)
		}, func(t *testing.T, got any) {
			if got.(int) != 3 {
				t.Errorf("added %v", got)
			}
		}},
		{"suppressions_list", func(c *Client) (any, error) { return c.Suppressions.List(ctx, &SuppressionListParams{Limit: 2}) }, func(t *testing.T, got any) {
			page := got.(*SuppressionPage)
			if len(page.Suppressions) != 2 || page.NextCursor == "" || page.Suppressions[0].Value != "umbrella" || page.Suppressions[1].Kind != SuppressDomain {
				t.Errorf("page = %+v", page)
			}
		}},
		{"suppressions_upload", func(c *Client) (any, error) {
			return c.Suppressions.Upload(ctx, strings.NewReader("kind,value\ndomain,hooli.example\nemail,gavin@hooli.example\ncompany_name,Pied Piper\nemail,not an address\n"), nil)
		}, func(t *testing.T, got any) {
			u := got.(*SuppressionUpload)
			if u.Added != 3 || u.Already != 0 || len(u.Invalid) != 1 || u.Invalid[0].Index != 5 {
				t.Errorf("upload = %+v", u)
			}
		}},
		{"suppressions_delete", func(c *Client) (any, error) { return nil, c.Suppressions.Delete(ctx, SuppressDomain, "globex.example") }, nil},
		{"attestations_list", func(c *Client) (any, error) { return c.Attestations.List(ctx) }, func(t *testing.T, got any) {
			list := got.([]Attestation)
			if len(list) != 1 || list[0].Kind != "personal_data" || !list[0].Accepted || list[0].AcceptedAt == nil || list[0].AcceptedBy != "ada@example.com" ||
				list[0].TextVersion != "2026-11-24" {
				t.Errorf("attestations = %+v", list)
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.fixture, func(t *testing.T) {
			client, api := replay(t, tc.fixture)
			got, err := tc.call(client)
			if err != nil {
				t.Fatal(err)
			}
			sent := api.requests()
			if len(sent) != 1 {
				t.Fatalf("%d requests, want 1", len(sent))
			}
			fx := loadFixture(t, tc.fixture)
			checkSent(t, sent[0], fx)
			if accept := sent[0].Header.Get("Accept"); fx.Request.Query["format"] == "csv" && accept != "text/csv" ||
				fx.Request.Query["format"] == "ndjson" && accept != "application/x-ndjson" {
				t.Errorf("Accept = %q", accept)
			}
			if tc.check != nil {
				tc.check(t, got)
			}
		})
	}
}

// Every error the formula routes answer: what errors.Is matches, and what the error carries.
func TestFormulaErrors(t *testing.T) {
	ctx := context.Background()
	terms := map[string]any{"terms": []string{"java"}}
	cases := []struct {
		fixture string
		call    func(c *Client) error
		is      []error
		isNot   []error
		check   func(t *testing.T, e *APIError)
	}{
		{"error_attestation_required", func(c *Client) error {
			_, err := c.Formulas.Run(ctx, "job-order-hunter", jobOrder, &FormulaRunParams{IdempotencyKey: "c0ffee-attest"})
			return err
		}, []error{ErrAttestationRequired, ErrConflict}, []error{ErrPermission, ErrFormulaUnavailable}, func(t *testing.T, e *APIError) {
			if e.Details["kind"] != "personal_data" || e.Details["text_version"] != "2026-11-24" {
				t.Errorf("details = %v", e.Details)
			}
		}},
		{"error_formula_not_found", func(c *Client) error { _, err := c.Formulas.Get(ctx, "no-such-formula"); return err },
			[]error{ErrNotFound}, []error{ErrConflict, ErrValidation}, func(t *testing.T, e *APIError) {
				if e.Code != CodeFormulaNotFound {
					t.Errorf("code = %q", e.Code)
				}
			}},
		{"error_formula_unavailable", func(c *Client) error {
			_, err := c.Formulas.Run(ctx, "retired-formula", terms, &FormulaRunParams{IdempotencyKey: "c0ffee-retired"})
			return err
		}, []error{ErrFormulaUnavailable, ErrConflict}, []error{ErrRequirementMissing, ErrNotFound}, func(t *testing.T, e *APIError) {
			if e.Details["status"] != FormulaDeprecated {
				t.Errorf("details = %v", e.Details)
			}
		}},
		{"error_requirement_missing", func(c *Client) error {
			_, err := c.Formulas.Run(ctx, "data-formula", terms, &FormulaRunParams{IdempotencyKey: "c0ffee-data"})
			return err
		}, []error{ErrRequirementMissing, ErrConflict}, []error{ErrAttestationRequired}, func(t *testing.T, e *APIError) {
			if m := e.Missing(); len(m) != 1 || m[0] != "product:data" {
				t.Errorf("missing = %v", m)
			}
		}},
		{"error_unit_cap_too_low", func(c *Client) error {
			_, err := c.Formulas.Run(ctx, "job-order-hunter", jobOrder, &FormulaRunParams{UnitCap: 1, IdempotencyKey: "c0ffee-cap"})
			return err
		}, []error{ErrValidation}, []error{ErrConflict, ErrQuotaExceeded}, func(t *testing.T, e *APIError) {
			if n, ok := e.MinUnitCap(); !ok || n != 20 || e.Code != CodeUnitCapTooLow {
				t.Errorf("min unit cap = %d %v", n, ok)
			}
			if _, ok := (&APIError{Details: map[string]any{}}).MinUnitCap(); ok {
				t.Error("an error without min_unit_cap should say so")
			}
		}},
		{"error_limit_reached", func(c *Client) error {
			_, err := c.SavedFormulas.Run(ctx, sf1, &SavedFormulaRunParams{IdempotencyKey: "c0ffee-saved-cap"})
			return err
		}, []error{ErrLimitReached, ErrQuotaExceeded}, []error{ErrPermission}, func(t *testing.T, e *APIError) {
			if e.Details["scope"] != "saved_formula" || e.Details["saved_formula_id"] != sf1 || e.Details["units_needed"] != float64(10) {
				t.Errorf("details = %v", e.Details)
			}
		}},
		{"error_plan_limit_reached", func(c *Client) error {
			_, err := c.SavedFormulas.Create(ctx, "posting-watch", terms, &SavedFormulaCreateParams{
				Schedule: &ScheduleParams{Cadence: CadenceHourly, LocalTime: "07:30", Timezone: "Europe/London"}})
			return err
		}, []error{ErrPlanLimitReached, ErrPermission}, []error{ErrQuotaExceeded, ErrLimitReached}, func(t *testing.T, e *APIError) {
			if e.Details["limit"] != "min_interval_minutes" || e.Details["allowed"] != float64(1440) || e.Details["asked"] != float64(60) {
				t.Errorf("details = %v", e.Details)
			}
		}},
		{"error_saved_formula_archived", func(c *Client) error {
			_, err := c.SavedFormulas.Update(ctx, "sf_000000000000000000000002", SavedFormulaUpdateParams{Notify: NotifyNone})
			return err
		}, []error{ErrSavedFormulaArchived, ErrConflict}, []error{ErrNotFound}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.fixture, func(t *testing.T) {
			client, api := replay(t, tc.fixture)
			err := tc.call(client)
			var e *APIError
			if !errors.As(err, &e) {
				t.Fatalf("got %T %v, want *APIError", err, err)
			}
			fx := loadFixture(t, tc.fixture)
			if sent := api.requests(); len(sent) != 1 {
				t.Errorf("sent %d times: an answer of %d is not tried again", len(sent), fx.Response.Status)
			} else {
				checkSent(t, sent[0], fx)
			}
			if e.StatusCode != fx.Response.Status || !strings.Contains(err.Error(), e.Code) || e.Message == "" {
				t.Errorf("error = %+v", e)
			}
			for _, target := range tc.is {
				if !errors.Is(err, target) {
					t.Errorf("errors.Is(%v) should be true", target)
				}
			}
			for _, target := range tc.isNot {
				if errors.Is(err, target) {
					t.Errorf("errors.Is(%v) should be false", target)
				}
			}
			if tc.check != nil {
				tc.check(t, e)
			}
		})
	}
}

// A parameter type encodes as the params the real API was sent, and is refused, before anything is sent, for
// another formula.
func TestFormulaParamTypes(t *testing.T) {
	fx := loadFixture(t, "formulas_estimate")
	var body struct {
		Params any `json:"params"`
	}
	_ = json.Unmarshal(fx.Request.Body, &body)
	encoded, err := json.Marshal(jobOrder)
	if err != nil {
		t.Fatal(err)
	}
	var params any
	_ = json.Unmarshal(encoded, &params)
	if !reflect.DeepEqual(params, body.Params) {
		t.Errorf("JobOrderHunterParams encodes as\n %s\nwant\n %v", encoded, body.Params)
	}

	client, api := replay(t, "formulas_run_queued")
	ctx := context.Background()
	refused := []func() error{
		func() error { _, err := client.Formulas.Run(ctx, "account-research", jobOrder, nil); return err },
		func() error { _, err := client.Formulas.Run(ctx, "account-research", &jobOrder, nil); return err },
		func() error { _, err := client.Formulas.Estimate(ctx, "account-research", jobOrder, nil); return err },
		func() error {
			_, err := client.SavedFormulas.Create(ctx, "account-research", jobOrder, nil)
			return err
		},
		func() error { _, _, err := client.RunFormula(ctx, "account-research", jobOrder, nil); return err },
	}
	for i, call := range refused {
		if err := call(); err == nil || !strings.Contains(err.Error(), `"job-order-hunter", not of "account-research"`) {
			t.Errorf("%d: got %v", i, err)
		}
	}
	if n := len(api.requests()); n != 0 {
		t.Fatalf("a mismatched parameter type must not be sent; %d requests were", n)
	}

	// No params, a nil map and a nil pointer are all {}; a map or a struct of yours is sent as it is.
	var none *JobOrderHunterParams
	var nilMap map[string]any
	for _, p := range []any{nil, none, nilMap, map[string]any{"job_titles": []string{"Nurse"}}, struct {
		Locations []string `json:"locations"`
	}{[]string{"Leeds"}}} {
		if _, err := client.Formulas.Run(ctx, "job-order-hunter", p, nil); err != nil {
			t.Fatal(err)
		}
	}
	var bodies []string
	for _, r := range api.requests() {
		bodies = append(bodies, string(r.Body))
	}
	want := []string{`{"params":{}}`, `{"params":{}}`, `{"params":{}}`, `{"params":{"job_titles":["Nurse"]}}`, `{"params":{"locations":["Leeds"]}}`}
	if !reflect.DeepEqual(bodies, want) {
		t.Errorf("bodies = %q", bodies)
	}
	if p := Ptr(3); *p != 3 {
		t.Error("Ptr")
	}
}

func TestFormulaRunStartsCarryAnIdempotencyKey(t *testing.T) {
	var keys []string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		keys = append(keys, r.Header.Get("Idempotency-Key"))
		if len(keys) == 1 {
			writeJSON(w, 503, apiErrorBody("internal", "x", ""))
			return
		}
		loadFixture(t, "formulas_run_queued").write(w)
	})
	ctx := context.Background()
	if _, err := c.Formulas.Run(ctx, "job-order-hunter", jobOrder, nil); err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 || !regexp.MustCompile(`^[0-9a-f]{32}$`).MatchString(keys[0]) || keys[0] != keys[1] {
		t.Errorf("keys = %v; want one new key, sent again with the retry", keys)
	}
	if _, err := c.SavedFormulas.Run(ctx, sf2, nil); err != nil {
		t.Fatal(err)
	}
	if len(keys) != 3 || keys[2] == keys[0] || len(keys[2]) != 32 {
		t.Errorf("keys = %v; each start has its own key", keys)
	}
	if _, err := c.Formulas.Run(ctx, "job-order-hunter", nil, &FormulaRunParams{Wait: MaxWait + time.Second}); err == nil {
		t.Error("a wait the API would refuse should not be sent")
	}
	if _, err := c.FormulaRuns.GetWait(ctx, run2, -time.Second); err == nil {
		t.Error("a negative wait should not be sent")
	}
	if len(keys) != 3 {
		t.Errorf("%d requests", len(keys))
	}
}

func TestSavedFormulaCreateAndDuplicateAreNotSentAgain(t *testing.T) {
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 500, apiErrorBody("internal", "Oops.", ""))
	})
	ctx := context.Background()
	if _, err := c.SavedFormulas.Create(ctx, "job-order-hunter", jobOrder, nil); !errors.Is(err, ErrServer) {
		t.Fatalf("got %v", err)
	}
	if _, err := c.SavedFormulas.Duplicate(ctx, sf2, nil); !errors.Is(err, ErrServer) {
		t.Fatalf("got %v", err)
	}
	if n := len(rec.calls()); n != 2 {
		t.Errorf("sent %d times, want once each", n)
	}
	if body := rec.calls()[1].Body; body != `{}` {
		t.Errorf("a duplicate with no options sends {}, got %s", body)
	}
}

// Update sends only what is set, and the API's null for what is taken away.
func TestSavedFormulaUpdateSendsNulls(t *testing.T) {
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) { loadFixture(t, "saved_get").write(w) })
	ctx := context.Background()
	for _, change := range []SavedFormulaUpdateParams{
		{AutoName: true},
		{RemoveSchedule: true, RemoveMonthlyUnitCap: true},
		{Schedule: &ScheduleParams{Cadence: CadenceWeekly, LocalTime: "08:00", Timezone: "UTC", DayOfWeek: Ptr(0)}, Notify: NotifyEachRun},
		{Params: jobOrder, ResetSeen: true, UnitCap: 500, RowCap: 50},
	} {
		if _, err := c.SavedFormulas.Update(ctx, sf2, change); err != nil {
			t.Fatal(err)
		}
	}
	want := []string{
		`{"name":null}`,
		`{"monthly_unit_cap":null,"schedule":null}`,
		`{"notify":"each_run","schedule":{"cadence":"weekly","local_time":"08:00","timezone":"UTC","day_of_week":0}}`,
		`{"params":{"job_titles":["Java Developer"],"locations":["London"],"posted_within_days":2,"company_size":{"min":11,"max":5000},` +
			`"exclude_companies":["Blocked By Us Inc"],"title_exclude_keywords":["intern"]},"reset_seen":true,"row_cap":50,"unit_cap":500}`,
	}
	for i, r := range rec.calls() {
		if r.Method != "PATCH" || r.Path != "/v1/saved-formulas/"+sf2 || r.Body != want[i] {
			t.Errorf("%d: %s %s %s\nwant %s", i, r.Method, r.Path, r.Body, want[i])
		}
	}
	for _, bad := range []SavedFormulaUpdateParams{{}, {Name: "x", AutoName: true}, {MonthlyUnitCap: 5, RemoveMonthlyUnitCap: true},
		{Schedule: &ScheduleParams{}, RemoveSchedule: true}} {
		if _, err := c.SavedFormulas.Update(ctx, sf2, bad); err == nil {
			t.Errorf("%+v should be refused", bad)
		}
	}
	if n := len(rec.calls()); n != 4 {
		t.Errorf("a refused change must not be sent; %d requests", n)
	}
}

// Rows are read page after page, as the cursor says.
func TestFormulaRowsAcrossPages(t *testing.T) {
	client, api, _ := routed(t, func(r *http.Request, _ int) string {
		if r.URL.Query().Get("cursor") == "Mg" {
			return "formula_runs_rows_page2"
		}
		return "formula_runs_rows"
	})
	rows, err := client.FormulaRuns.Rows(context.Background(), run2, &FormulaRowListParams{Limit: 2}).Collect()
	if err != nil || len(rows) != 4 || rows[3].ID != "frow_000000000000000000000004" {
		t.Fatalf("%d rows, %v", len(rows), err)
	}
	sent := api.requests()
	if len(sent) != 2 {
		t.Fatalf("%d requests, want 2", len(sent))
	}
	checkSent(t, sent[0], loadFixture(t, "formula_runs_rows"))
	checkSent(t, sent[1], loadFixture(t, "formula_runs_rows_page2"))

	// The saved formula's rows: its first page, then the page its cursor names.
	client, api, _ = routed(t, func(r *http.Request, n int) string { return "saved_rows" })
	it := client.SavedFormulas.Rows(context.Background(), sf2, &SavedFormulaRowListParams{Since: time.Date(2026, 11, 25, 9, 2, 55, 250_000_000, time.UTC), Limit: 2})
	for i := 0; i < 3 && it.Next(); i++ {
	}
	sent = api.requests()
	checkSent(t, sent[0], loadFixture(t, "saved_rows"))
	if len(sent) != 2 || sent[1].Query.Get("cursor") == "" || sent[1].Query.Get("since") != "2026-11-25T09:02:55.250000+00:00" {
		t.Errorf("the second page: %+v", sent)
	}

	// Runs, saved formulas and the suppression list page the same way.
	client, api, _ = routed(t, func(r *http.Request, n int) string {
		switch {
		case strings.HasPrefix(r.URL.Path, "/v1/formula-runs"):
			return "formula_runs_list"
		case strings.HasPrefix(r.URL.Path, "/v1/saved-formulas"):
			return "saved_list"
		}
		return "suppressions_list"
	})
	ctx := context.Background()
	runs := client.FormulaRuns.ListAll(ctx, &FormulaRunListParams{FormulaID: "job-order-hunter", Status: RunPartial})
	if !runs.Next() || !runs.Next() {
		t.Fatal(runs.Err())
	}
	saved, err := client.SavedFormulas.ListAll(ctx, &SavedFormulaListParams{Status: SavedFormulaActive}).Collect()
	if err != nil || len(saved) != 1 {
		t.Fatalf("%v %v", saved, err)
	}
	entries := client.Suppressions.ListAll(ctx, &SuppressionListParams{Kind: SuppressDomain})
	if !entries.Next() || !entries.Next() || !entries.Next() {
		t.Fatal(entries.Err())
	}
	sent = api.requests()
	if q := sent[0].Query.Encode(); q != "formula_id=job-order-hunter&status=partial" {
		t.Errorf("runs: %s", q)
	}
	if sent[1].Query.Get("cursor") == "" || sent[2].Query.Encode() != "status=active" || sent[3].Query.Encode() != "kind=domain" || sent[4].Query.Get("cursor") == "" {
		t.Errorf("sent %+v", sent)
	}
}

// Wait looks until the run has ended: each look asks the API to hold the request, and two looks are at least the
// poll interval apart.
func TestFormulaRunWait(t *testing.T) {
	client, api, rec := routed(t, func(r *http.Request, n int) string {
		if n < 3 {
			return "formula_runs_get_running"
		}
		return "formula_runs_get"
	})
	ctx := context.Background()
	run, err := client.FormulaRuns.Get(ctx, run2)
	if err != nil || run.Status != RunRunning {
		t.Fatalf("%v %v", run, err)
	}
	checkSent(t, api.requests()[0], loadFixture(t, "formula_runs_get_running"))
	var seen []RunStatus
	run, err = client.FormulaRuns.Wait(ctx, run2, &FormulaWaitParams{PollInterval: 3 * time.Second, OnPoll: func(r *FormulaRun) { seen = append(seen, r.Status) }})
	if err != nil || run.Status != RunPartial || run.Reason != "contacts_not_connected" {
		t.Fatalf("run %+v, err %v", run, err)
	}
	sent := api.requests()
	if len(sent) != 3 || len(seen) != 2 || seen[0] != RunRunning {
		t.Fatalf("%d requests, seen %v", len(sent), seen)
	}
	for _, r := range sent[1:] {
		checkSent(t, r, loadFixture(t, "formula_runs_get")) // ?wait=60
	}
	if len(rec.delays) != 1 || rec.delays[0] <= 2*time.Second || rec.delays[0] > 3*time.Second {
		t.Errorf("pauses = %v; want one of the rest of 3s", rec.delays)
	}

	// The wait asked for is what is left of the timeout, in whole seconds.
	client, api, _ = routed(t, func(*http.Request, int) string { return "formula_runs_get" })
	if _, err := client.FormulaRuns.Wait(ctx, run2, &FormulaWaitParams{Timeout: 30 * time.Second}); err != nil {
		t.Fatal(err)
	}
	if q := api.requests()[0].Query.Get("wait"); q != "29" {
		t.Errorf("wait=%s, want 29", q)
	}
}

func TestFormulaRunWaitTimesOut(t *testing.T) {
	client, api, _ := routed(t, func(*http.Request, int) string { return "formula_runs_get_running" })
	client.sleep = sleepContext
	run, err := client.FormulaRuns.Wait(context.Background(), run2, &FormulaWaitParams{Timeout: 300 * time.Millisecond, PollInterval: 50 * time.Millisecond})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got %v, want a deadline error", err)
	}
	if run == nil || run.ID != run2 || run.Status != RunRunning {
		t.Errorf("the run as last seen should come back with the error, got %v", run)
	}
	if n := len(api.requests()); n < 2 || n > 8 {
		t.Errorf("%d looks in 300ms, 50ms apart", n)
	}

	// The caller's context ends it too, with the context's error.
	ctx, cancel := context.WithCancel(context.Background())
	client, _, _ = routed(t, func(*http.Request, int) string { return "formula_runs_get_running" })
	_, err = client.FormulaRuns.Wait(ctx, run2, &FormulaWaitParams{OnPoll: func(*FormulaRun) { cancel() }})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("got %v, want context.Canceled", err)
	}
}

// RunFormula: the start, waiting, and the rows.
func TestRunFormula(t *testing.T) {
	client, api, _ := routed(t, func(r *http.Request, _ int) string {
		switch {
		case r.Method == http.MethodPost:
			return "formulas_run"
		case r.URL.Query().Get("cursor") == "Mg":
			return "formula_runs_rows_page2"
		}
		return "formula_runs_rows"
	})
	ctx := context.Background()
	run, rows, err := client.RunFormula(ctx, "job-order-hunter", jobOrder, &RunFormulaParams{Preview: true, IdempotencyKey: "c0ffee-run-2"})
	if err != nil || run.Status != RunPartial || run.Mode != ModePreview {
		t.Fatalf("run %+v, err %v", run, err)
	}
	all, err := rows.Collect()
	if err != nil || len(all) != 4 {
		t.Fatalf("%d rows, %v", len(all), err)
	}
	sent := api.requests()
	checkSent(t, sent[0], loadFixture(t, "formulas_run")) // ?wait=60, mode preview
	if len(sent) != 3 || sent[1].Path != "/v1/formula-runs/frun_000000000000000000000003/rows" || sent[1].Query.Encode() != "limit=1000" ||
		sent[2].Query.Encode() != "cursor=Mg&limit=1000" {
		t.Errorf("sent %+v", sent[1:])
	}

	// A run still queued when the start answers is waited for.
	client, api, _ = routed(t, func(r *http.Request, _ int) string {
		switch {
		case r.Method == http.MethodPost:
			return "formulas_run_queued"
		case strings.HasSuffix(r.URL.Path, "/rows"):
			return "formula_runs_rows_page2"
		}
		return "formula_runs_get"
	})
	run, rows, err = client.RunFormula(ctx, "job-order-hunter", jobOrder, &RunFormulaParams{UnitCap: 2000})
	if err != nil || run.ID != run2 || run.Status != RunPartial {
		t.Fatalf("run %+v, err %v", run, err)
	}
	if all, err := rows.Collect(); err != nil || len(all) != 2 {
		t.Fatalf("%d rows, %v", len(all), err)
	}
	sent = api.requests()
	checkSent(t, sent[0], loadFixture(t, "formulas_run_queued"), "wait", "Idempotency-Key")
	checkSent(t, sent[1], loadFixture(t, "formula_runs_get"))
	if sent[0].Query.Get("wait") != "60" || len(sent[0].Header.Get("Idempotency-Key")) != 32 {
		t.Errorf("the start: %+v", sent[0])
	}

	// A shorter timeout asks the API to wait no longer than that.
	client, api, _ = routed(t, func(*http.Request, int) string { return "formulas_run" })
	if _, _, err := client.RunFormula(ctx, "job-order-hunter", jobOrder, &RunFormulaParams{Timeout: 20 * time.Second}); err != nil {
		t.Fatal(err)
	}
	if q := api.requests()[0].Query.Get("wait"); q != "20" {
		t.Errorf("wait=%s", q)
	}
}

func TestRunFormulaThatFailed(t *testing.T) {
	fx := loadFixture(t, "formulas_run")
	failed := bytes.Replace(bytes.Replace(fx.Response.Body, []byte(`"status": "partial"`), []byte(`"status": "failed"`), 1),
		[]byte(`"reason": "contacts_not_connected"`), []byte(`"reason": "source_failed:linkedin.jobs_search"`), 1)
	if bytes.Equal(failed, fx.Response.Body) {
		t.Fatal("the fixture did not change")
	}
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, string(failed)) })
	run, rows, err := c.RunFormula(context.Background(), "job-order-hunter", jobOrder, nil)
	var e *FormulaRunFailedError
	if !errors.As(err, &e) || !errors.Is(err, ErrRunFailed) {
		t.Fatalf("got %T %v", err, err)
	}
	if e.Status() != RunFailed || e.Reason() != "source_failed:linkedin.jobs_search" || e.Run.ID != "frun_000000000000000000000003" ||
		!strings.Contains(err.Error(), "formula run frun_000000000000000000000003 ended failed (reason: source_failed:linkedin.jobs_search)") {
		t.Errorf("error = %v", err)
	}
	if run == nil || rows != nil || len(rec.calls()) != 1 {
		t.Errorf("the run comes back with the error, the rows do not: %v %v", run, rows)
	}
	run, rows, err = c.RunFormula(context.Background(), "job-order-hunter", jobOrder, &RunFormulaParams{AllowFailed: true})
	if err != nil || run.Err() == nil || rows == nil {
		t.Errorf("AllowFailed: %v %v %v", run, rows, err)
	}
	for status, want := range map[RunStatus]bool{RunOK: false, RunPartial: false, RunRunning: false, RunFailed: true, RunUnsupported: true, RunCancelled: true} {
		if got := (&FormulaRun{Status: status}).Err() != nil; got != want {
			t.Errorf("%s: Err() != nil is %v", status, got)
		}
	}
}

// ExportAll writes one file: the header of a CSV once, every page after it.
func TestFormulaExportAll(t *testing.T) {
	client, api, _ := routed(t, func(r *http.Request, _ int) string {
		switch {
		case r.URL.Query().Get("format") == FormatNDJSON:
			return "formula_runs_rows_ndjson"
		case strings.HasPrefix(r.URL.Path, "/v1/saved-formulas"):
			return "saved_rows_csv"
		case r.URL.Query().Get("cursor") == "Mg":
			return "formula_runs_rows_csv_page2"
		}
		return "formula_runs_rows_csv"
	})
	ctx := context.Background()
	var out bytes.Buffer
	if err := client.FormulaRuns.ExportAll(ctx, run2, FormulaExportParams{Format: FormatCSV, Limit: 2}, &out); err != nil {
		t.Fatal(err)
	}
	first, second := loadFixture(t, "formula_runs_rows_csv"), loadFixture(t, "formula_runs_rows_csv_page2")
	header := (*first.Response.Text)[:strings.Index(*first.Response.Text, "\n")+1]
	want := *first.Response.Text + strings.TrimPrefix(*second.Response.Text, header)
	if out.String() != want || strings.Count(out.String(), header) != 1 || strings.Count(out.String(), "\r\n") != 5 {
		t.Errorf("the whole run:\n%s", out.String())
	}
	sent := api.requests()
	checkSent(t, sent[0], first)
	checkSent(t, sent[1], second)

	out.Reset()
	if err := client.FormulaRuns.ExportAll(ctx, run2, FormulaExportParams{Format: FormatNDJSON}, &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != *loadFixture(t, "formula_runs_rows_ndjson").Response.Text {
		t.Errorf("ndjson:\n%s", out.String())
	}
	checkSent(t, api.requests()[2], loadFixture(t, "formula_runs_rows_ndjson"), "limit")
	if q := api.requests()[2].Query.Get("limit"); q != "1000" {
		t.Errorf("ExportAll reads 1000 a page, asked %s", q)
	}

	out.Reset()
	if err := client.SavedFormulas.ExportAll(ctx, sf2, SavedFormulaExportParams{Format: FormatCSV, Fields: []string{"company_name", "job_title", "job_url"}}, &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != *loadFixture(t, "saved_rows_csv").Response.Text {
		t.Errorf("saved csv:\n%s", out.String())
	}
	checkSent(t, api.requests()[3], loadFixture(t, "saved_rows_csv"), "limit")

	if _, err := client.FormulaRuns.Export(ctx, run2, FormulaExportParams{Format: "json"}); err == nil {
		t.Error("a format the API does not write should be refused")
	}
	if err := client.SavedFormulas.ExportAll(ctx, sf2, SavedFormulaExportParams{Format: "xml"}, &out); err == nil {
		t.Error("a format the API does not write should be refused")
	}
	if n := len(api.requests()); n != 4 {
		t.Errorf("%d requests", n)
	}
}

// The catalogue of formulas needs no key; everything else does.
func TestFormulaCatalogueIsPublic(t *testing.T) {
	t.Setenv(EnvAPIKey, "")
	var auths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auths = append(auths, r.Header.Get("Authorization"))
		if r.URL.Path == "/v1/formulas" {
			loadFixture(t, "formulas_list").write(w)
			return
		}
		loadFixture(t, "formulas_get").write(w)
	}))
	defer server.Close()
	c, err := NewClient(WithoutAPIKey(), WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if list, err := c.Formulas.List(ctx, nil); err != nil || len(list) != 3 {
		t.Fatalf("%v %v", list, err)
	}
	if f, err := c.Formulas.Get(ctx, "job-order-hunter"); err != nil || f.ID != "job-order-hunter" {
		t.Fatalf("%v %v", f, err)
	}
	if _, err := c.Formulas.Estimate(ctx, "job-order-hunter", jobOrder, nil); !errors.Is(err, ErrMissingAPIKey) {
		t.Errorf("estimate: %v", err)
	}
	if _, err := c.FormulaRuns.Get(ctx, run2); !errors.Is(err, ErrMissingAPIKey) {
		t.Errorf("a run: %v", err)
	}
	if len(auths) != 2 || auths[0] != "" || auths[1] != "" {
		t.Errorf("Authorization headers = %q", auths)
	}
}

func TestFormulaCallsRefuseEmptyIDsAndEntries(t *testing.T) {
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, `{}`) })
	ctx := context.Background()
	calls := []func() error{
		func() error { _, err := c.Formulas.Get(ctx, ""); return err },
		func() error { _, err := c.Formulas.Run(ctx, " ", nil, nil); return err },
		func() error { _, err := c.FormulaRuns.Cancel(ctx, ""); return err },
		func() error { _, err := c.SavedFormulas.Pause(ctx, ""); return err },
		func() error { _, err := c.SavedFormulas.Create(ctx, "", nil, nil); return err },
		func() error { return c.Suppressions.Delete(ctx, SuppressDomain, "") },
		func() error { _, err := c.Suppressions.Add(ctx, nil, nil); return err },
		func() error {
			_, err := c.Suppressions.Add(ctx, make([]SuppressionEntry, MaxSuppressionEntries+1), nil)
			return err
		},
		func() error { _, err := c.FormulaRows.Feedback(ctx, "", VerdictGood, nil); return err },
	}
	for i, call := range calls {
		if call() == nil {
			t.Errorf("%d: should be refused", i)
		}
	}
	if n := len(rec.calls()); n != 0 {
		t.Errorf("%d requests were sent", n)
	}
}

// A reason and a note keep every field the API sent, and write it back.
func TestReasonsAndNotesKeepTheirFields(t *testing.T) {
	var r FormulaRunReason
	if err := json.Unmarshal([]byte(`{"code":"source_not_offered","step":"search_google_jobs","scraper":"google.jobs"}`), &r); err != nil {
		t.Fatal(err)
	}
	again, _ := json.Marshal(r)
	if r.Code != "source_not_offered" || r.Fields["step"] != "search_google_jobs" || string(again) != `{"code":"source_not_offered","scraper":"google.jobs","step":"search_google_jobs"}` {
		t.Errorf("%+v %s", r, again)
	}
	var n SavedFormulaNote
	if err := json.Unmarshal([]byte(`{"code":"source_disabled","message":"Indeed is off for now.","at":"2026-11-25T09:00:00Z","scraper":"indeed.jobs_search"}`), &n); err != nil {
		t.Fatal(err)
	}
	if n.Code != "source_disabled" || n.Message != "Indeed is off for now." || n.At.Hour() != 9 || n.Fields["scraper"] != "indeed.jobs_search" {
		t.Errorf("%+v", n)
	}
	if again, _ := json.Marshal(n); !strings.Contains(string(again), `"scraper":"indeed.jobs_search"`) || !strings.Contains(string(again), `"code":"source_disabled"`) {
		t.Errorf("%s", again)
	}
}
