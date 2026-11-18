package enrichednow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const testKey = "en_live_0123456789abcdef0123456789abcdef"

// recorder is a fake API: it keeps every request it was sent and answers with handler.
type recorder struct {
	mu       sync.Mutex
	requests []recorded
	delays   []time.Duration
}

type recorded struct {
	Method, Path, Query, Auth, UserAgent, ContentType, Body string
}

func (r *recorder) calls() []recorded {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]recorded(nil), r.requests...)
}

// newTestClient serves handler and returns a client that talks to it and never really sleeps.
func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) (*Client, *recorder) {
	t.Helper()
	rec := &recorder{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		rec.mu.Lock()
		rec.requests = append(rec.requests, recorded{Method: req.Method, Path: req.URL.Path, Query: req.URL.RawQuery,
			Auth: req.Header.Get("Authorization"), UserAgent: req.Header.Get("User-Agent"),
			ContentType: req.Header.Get("Content-Type"), Body: string(body)})
		rec.mu.Unlock()
		handler(w, req)
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(append([]Option{WithAPIKey(testKey), WithBaseURL(server.URL)}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	client.sleep = func(ctx context.Context, d time.Duration) error {
		rec.mu.Lock()
		rec.delays = append(rec.delays, d)
		rec.mu.Unlock()
		return ctx.Err()
	}
	client.jitter = func() float64 { return 1 }
	return client, rec
}

func writeJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	io.WriteString(w, body)
}

func apiErrorBody(code, message, details string) string {
	if details == "" {
		details = "{}"
	}
	return fmt.Sprintf(`{"error":{"code":%q,"message":%q,"details":%s}}`, code, message, details)
}

func runJSON(id, status, reason string, records int) string {
	r := "null"
	if reason != "" {
		r = fmt.Sprintf("%q", reason)
	}
	finished := "null"
	if RunStatus(status).IsFinal() {
		finished = `"2026-11-16T22:37:51.066313Z"`
	}
	return fmt.Sprintf(`{"id":%q,"scraper_id":"ats.board","scraper_version":1,"status":%q,"reason":%s,`+
		`"input":{"board_url":"https://boards.greenhouse.io/acme"},"record_count":%d,"units":0,"diagnostics":null,`+
		`"created_at":"2026-11-16T22:37:50.170950Z","started_at":null,"finished_at":%s,"a_field_added_later":true}`,
		id, status, r, records, finished)
}

func TestNewClientAPIKey(t *testing.T) {
	t.Setenv(EnvAPIKey, "")
	if _, err := NewClient(); !errors.Is(err, ErrMissingAPIKey) {
		t.Fatalf("without a key: got %v, want ErrMissingAPIKey", err)
	}
	if !strings.Contains(ErrMissingAPIKey.Error(), "ENRICHEDNOW_API_KEY") {
		t.Errorf("the error should name the environment variable: %v", ErrMissingAPIKey)
	}
	t.Setenv(EnvAPIKey, "en_live_from_env")
	c, err := NewClient()
	if err != nil || c.apiKey != "en_live_from_env" {
		t.Fatalf("from the environment: %v, key %q", err, c.apiKey)
	}
	c, _ = NewClient(WithAPIKey("en_live_explicit"))
	if c.apiKey != "en_live_explicit" {
		t.Errorf("an explicit key should win over the environment, got %q", c.apiKey)
	}
	if c.baseURL != DefaultBaseURL || c.maxRetries != 2 || c.timeout != DefaultTimeout {
		t.Errorf("defaults: %q %d %v", c.baseURL, c.maxRetries, c.timeout)
	}
	if _, err := NewClient(WithAPIKey("k"), WithBaseURL("not a url")); err == nil {
		t.Error("a base URL that is no address should be refused")
	}
}

func TestAnonymousClient(t *testing.T) {
	t.Setenv(EnvAPIKey, "")
	rec := &recorder{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rec.mu.Lock()
		rec.requests = append(rec.requests, recorded{Path: req.URL.Path, Auth: req.Header.Get("Authorization")})
		rec.mu.Unlock()
		writeJSON(w, 200, `{"status":"ok","scrapers":[]}`)
	}))
	defer server.Close()
	c, err := NewClient(WithoutAPIKey(), WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	if h, err := c.Health(context.Background()); err != nil || h.Status != "ok" {
		t.Fatalf("health: %v %v", h, err)
	}
	if _, err := c.Scrapers.List(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, r := range rec.calls() {
		if r.Auth != "" {
			t.Errorf("%s: no Authorization header expected, got %q", r.Path, r.Auth)
		}
	}
	if _, err := c.Runs.Get(context.Background(), "run_1"); !errors.Is(err, ErrMissingAPIKey) {
		t.Errorf("a call that needs a key: got %v, want ErrMissingAPIKey", err)
	}
	if n := len(rec.calls()); n != 2 {
		t.Errorf("the keyless call must not reach the API; %d requests", n)
	}
}

func TestHeaders(t *testing.T) {
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, runJSON("run_1", "ok", "", 0)) },
		WithUserAgent("my-app/1.2"))
	if _, err := c.Scrapers.Run(context.Background(), "ats.board", map[string]any{"board_url": "https://x.example"}, nil); err != nil {
		t.Fatal(err)
	}
	got := rec.calls()[0]
	if got.Auth != "Bearer "+testKey {
		t.Errorf("Authorization = %q", got.Auth)
	}
	if want := "my-app/1.2 enrichednow-go/" + Version; got.UserAgent != want {
		t.Errorf("User-Agent = %q, want %q", got.UserAgent, want)
	}
	if got.ContentType != "application/json" {
		t.Errorf("Content-Type = %q", got.ContentType)
	}
}

// Every endpoint: the method, the path and the query it is called with.
func TestEndpoints(t *testing.T) {
	ctx := context.Background()
	yes := true
	seen := time.Date(2026, 11, 1, 12, 0, 0, 0, time.FixedZone("x", 3600))
	cases := []struct {
		name   string
		call   func(c *Client) error
		method string
		path   string
		query  string
		body   string
	}{
		{"health", func(c *Client) error { _, err := c.Health(ctx); return err }, "GET", "/v1/health", "", ""},
		{"scrapers list", func(c *Client) error { _, err := c.Scrapers.List(ctx); return err }, "GET", "/v1/scrapers", "", ""},
		{"scrapers get", func(c *Client) error { _, err := c.Scrapers.Get(ctx, "ats.board"); return err }, "GET", "/v1/scrapers/ats.board", "", ""},
		{"run", func(c *Client) error {
			_, err := c.Scrapers.Run(ctx, "careers.company", map[string]any{"domain": "acme.com"}, nil)
			return err
		}, "POST", "/v1/scrapers/careers.company/runs", "", `{"input":{"domain":"acme.com"}}`},
		{"run with wait", func(c *Client) error {
			_, err := c.Scrapers.Run(ctx, "careers.company", struct {
				Domain string `json:"domain"`
			}{"acme.com"}, &RunParams{Wait: 1500 * time.Millisecond})
			return err
		}, "POST", "/v1/scrapers/careers.company/runs", "wait=1.5", `{"input":{"domain":"acme.com"}}`},
		{"run without input", func(c *Client) error { _, err := c.Scrapers.Run(ctx, "x", nil, nil); return err },
			"POST", "/v1/scrapers/x/runs", "", `{"input":{}}`},
		{"runs get", func(c *Client) error { _, err := c.Runs.Get(ctx, "run_1"); return err }, "GET", "/v1/runs/run_1", "", ""},
		{"runs list", func(c *Client) error { _, err := c.Runs.List(ctx, nil); return err }, "GET", "/v1/runs", "", ""},
		{"runs list filtered", func(c *Client) error {
			_, err := c.Runs.List(ctx, &RunListParams{Limit: 50, Cursor: "abc", ScraperID: "careers.company", Status: RunPartial})
			return err
		}, "GET", "/v1/runs", "cursor=abc&limit=50&scraper_id=careers.company&status=partial", ""},
		{"runs cancel", func(c *Client) error { _, err := c.Runs.Cancel(ctx, "run_1"); return err }, "POST", "/v1/runs/run_1/cancel", "", ""},
		{"records page", func(c *Client) error {
			_, err := c.Runs.RecordsPage(ctx, "run_1", &RecordListParams{Limit: 2, Cursor: "Mg"})
			return err
		}, "GET", "/v1/runs/run_1/records", "cursor=Mg&limit=2", ""},
		{"usage", func(c *Client) error { _, err := c.Usage.Get(ctx, nil); return err }, "GET", "/v1/usage", "", ""},
		{"usage params", func(c *Client) error {
			_, err := c.Usage.Get(ctx, &UsageParams{From: "2026-11-01", To: "2026-11-16", GroupBy: "scraper", Product: "scrapers"})
			return err
		}, "GET", "/v1/usage", "from=2026-11-01&group_by=scraper&product=scrapers&to=2026-11-16", ""},
		{"companies list", func(c *Client) error {
			_, err := c.Data.Companies.List(ctx, &CompanyListParams{Query: "fig", ATS: "greenhouse", Limit: 2, Cursor: "c"})
			return err
		}, "GET", "/v1/data/companies", "ats=greenhouse&cursor=c&limit=2&q=fig", ""},
		{"company get", func(c *Client) error { _, err := c.Data.Companies.Get(ctx, "figma.com"); return err }, "GET", "/v1/data/companies/figma.com", "", ""},
		{"company postings", func(c *Client) error {
			_, err := c.Data.Companies.Postings(ctx, "figma.com", &PostingListParams{CompanyDomain: "ignored.com", Status: "open", Limit: 5})
			return err
		}, "GET", "/v1/data/companies/figma.com/postings", "limit=5&status=open", ""},
		{"postings list", func(c *Client) error {
			_, err := c.Data.Postings.List(ctx, &PostingListParams{
				CompanyDomain: "figma.com", Title: "program manager", Location: "New York", Country: "US", Remote: &yes,
				EmploymentType: "full_time", Seniority: "manager", Function: "product", Status: "open", Source: "greenhouse",
				Provenance: "feed,page", Unmatched: true, PostedFrom: "2026-06-01", PostedTo: "2026-06-30",
				FirstSeenFrom: seen, Sort: "last_seen_at", Limit: 10, Cursor: "n"})
			return err
		}, "GET", "/v1/data/postings", "company_domain=figma.com&country=US&cursor=n&employment_type=full_time" +
			"&first_seen_from=2026-11-01T11%3A00%3A00Z&function=product&limit=10&location=New+York&posted_from=2026-06-01" +
			"&posted_to=2026-06-30&provenance=feed%2Cpage&remote=true&seniority=manager&sort=last_seen_at&source=greenhouse" +
			"&status=open&title=program+manager&unmatched=true", ""},
		{"postings list plain", func(c *Client) error { _, err := c.Data.Postings.List(ctx, nil); return err }, "GET", "/v1/data/postings", "", ""},
		{"posting get", func(c *Client) error { _, err := c.Data.Postings.Get(ctx, "job_2a8c"); return err }, "GET", "/v1/data/postings/job_2a8c", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, `{}`) })
			if err := tc.call(c); err != nil {
				t.Fatal(err)
			}
			calls := rec.calls()
			if len(calls) != 1 {
				t.Fatalf("%d requests, want 1", len(calls))
			}
			got := calls[0]
			if got.Method != tc.method || got.Path != tc.path || got.Query != tc.query {
				t.Errorf("got  %s %s ?%s\nwant %s %s ?%s", got.Method, got.Path, got.Query, tc.method, tc.path, tc.query)
			}
			if got.Body != tc.body {
				t.Errorf("body = %s, want %s", got.Body, tc.body)
			}
		})
	}
}

func TestEmptyIDsAreRefusedLocally(t *testing.T) {
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, `{}`) })
	ctx := context.Background()
	if _, err := c.Runs.Get(ctx, ""); err == nil {
		t.Error("an empty run id should be an error")
	}
	if _, err := c.Scrapers.Get(ctx, " "); err == nil {
		t.Error("an empty scraper id should be an error")
	}
	if n := len(rec.calls()); n != 0 {
		t.Errorf("%d requests were sent", n)
	}
}

func TestRunDecoding(t *testing.T) {
	const body = `{"id":"run_841742f37553189ff9c27ba1","scraper_id":"ats.board","scraper_version":1,"status":"ok","reason":null,` +
		`"input":{"board_url":"https://boards.greenhouse.io/acme"},"record_count":21,"units":1,"diagnostics":{"requests":1,` +
		`"bytes":259214,"wall_ms":885,"facts":{"careers_url":"https://boards.greenhouse.io/acme","ats":"greenhouse","pages_read":1}},` +
		`"created_at":"2026-11-16T22:37:50.170950Z","started_at":"2026-11-16T22:37:50.178331Z","finished_at":"2026-11-16T22:37:51.066313Z"}`
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, body) })
	run, err := c.Runs.Get(context.Background(), "run_841742f37553189ff9c27ba1")
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != RunOK || run.Reason != "" || run.RecordCount != 21 || run.Units != 1 || run.ScraperVersion != 1 {
		t.Errorf("run = %+v", run)
	}
	if run.Diagnostics == nil || *run.Diagnostics.Bytes != 259214 || *run.Diagnostics.WallMS != 885 || run.Diagnostics.Facts["ats"] != "greenhouse" {
		t.Errorf("diagnostics = %+v", run.Diagnostics)
	}
	if run.FinishedAt == nil || run.FinishedAt.Sub(run.CreatedAt) <= 0 || run.CreatedAt.Year() != 2026 {
		t.Errorf("times: %v %v", run.CreatedAt, run.FinishedAt)
	}
	if string(run.RawJSON()) != body {
		t.Errorf("RawJSON should be the body as sent")
	}
	if run.Err() != nil {
		t.Errorf("an ok run has no error, got %v", run.Err())
	}
}

func TestErrorMapping(t *testing.T) {
	cases := []struct {
		status   int
		code     string
		details  string
		is       []error
		isNot    []error
		wantText string
	}{
		{401, "unauthorized", "", []error{ErrAuthentication}, []error{ErrNotFound, ErrServer}, "401 unauthorized"},
		{403, "forbidden", `{"reason":"email_not_verified"}`, []error{ErrPermission}, []error{ErrAuthentication}, "403 forbidden"},
		{404, "not_found", "", []error{ErrNotFound}, []error{ErrValidation}, "404 not_found"},
		{422, "invalid_input", `{"fields":[{"path":"input.board_url","message":"Field required"},{"path":"limit","message":"too big"}]}`,
			[]error{ErrValidation}, []error{ErrNotFound}, "[input.board_url: Field required]"},
		{402, "quota_exceeded", `{"used_units":1,"allowance_units":1}`, []error{ErrQuotaExceeded}, []error{ErrRateLimit}, "402 quota_exceeded"},
		{429, "too_many_runs", `{"active_runs":2,"max_concurrent_runs":2}`, []error{ErrRateLimit, ErrTooManyRuns}, []error{ErrServer}, "too_many_runs"},
		{429, "rate_limited", "", []error{ErrRateLimit}, []error{ErrTooManyRuns}, "rate_limited"},
		{500, "internal", "", []error{ErrServer}, []error{ErrDataUnavailable}, "500 internal"},
		{503, "data_unavailable", "", []error{ErrServer, ErrDataUnavailable}, []error{ErrRateLimit}, "data_unavailable"},
	}
	for _, tc := range cases {
		t.Run(tc.code, func(t *testing.T) {
			c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("X-Request-Id", "req_123")
				writeJSON(w, tc.status, apiErrorBody(tc.code, "A sentence.", tc.details))
			}, WithMaxRetries(0))
			_, err := c.Runs.Get(context.Background(), "run_1")
			var api *APIError
			if !errors.As(err, &api) {
				t.Fatalf("got %T %v, want *APIError", err, err)
			}
			if api.StatusCode != tc.status || api.Code != tc.code || api.Message != "A sentence." || api.RequestID != "req_123" {
				t.Errorf("error = %+v", api)
			}
			if api.Method != "GET" || api.Path != "/v1/runs/run_1" || api.Details == nil {
				t.Errorf("error = %+v", api)
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
			if !strings.Contains(err.Error(), tc.wantText) {
				t.Errorf("Error() = %q, should contain %q", err.Error(), tc.wantText)
			}
		})
	}
}

func TestErrorDetails(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 422, apiErrorBody("invalid_input", "The input is not valid.",
			`{"fields":[{"path":"input.board_url","message":"Field required"},{"path":"input.board","message":"Extra inputs are not permitted"}]}`))
	})
	_, err := c.Scrapers.Run(context.Background(), "ats.board", map[string]any{"board": "x"}, nil)
	var api *APIError
	if !errors.As(err, &api) {
		t.Fatal(err)
	}
	fields := api.Fields()
	if len(fields) != 2 || fields[0] != (FieldError{"input.board_url", "Field required"}) || fields[1].Path != "input.board" {
		t.Errorf("fields = %+v", fields)
	}
}

func TestErrorThatIsNotTheAPIs(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(502)
		io.WriteString(w, "<html>bad gateway</html>")
	}, WithMaxRetries(0))
	_, err := c.Runs.Get(context.Background(), "run_1")
	var api *APIError
	if !errors.As(err, &api) || api.StatusCode != 502 || api.Code != "" || api.Message != "Bad Gateway" || !errors.Is(err, ErrServer) {
		t.Fatalf("got %#v", err)
	}
}

func TestRetriesSafeRequests(t *testing.T) {
	var n atomic.Int32
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		n.Add(1)
		if n.Load() < 3 {
			writeJSON(w, 500, apiErrorBody("internal", "Something went wrong on our side.", ""))
			return
		}
		writeJSON(w, 200, runJSON("run_1", "ok", "", 0))
	})
	run, err := c.Runs.Get(context.Background(), "run_1")
	if err != nil || run.ID != "run_1" {
		t.Fatalf("got %v, %v", run, err)
	}
	if n.Load() != 3 {
		t.Errorf("%d attempts, want 3 (the default is 2 retries)", n.Load())
	}
	// With the jitter pinned at its top: 0.5s, then 1s.
	if len(rec.delays) != 2 || rec.delays[0] != 500*time.Millisecond || rec.delays[1] != time.Second {
		t.Errorf("delays = %v", rec.delays)
	}
}

func TestRetriesGiveUp(t *testing.T) {
	var n atomic.Int32
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		n.Add(1)
		writeJSON(w, 503, apiErrorBody("data_unavailable", "Try again.", ""))
	}, WithMaxRetries(1))
	_, err := c.Data.Postings.List(context.Background(), nil)
	if !errors.Is(err, ErrDataUnavailable) || n.Load() != 2 {
		t.Fatalf("err = %v after %d attempts; want data_unavailable after 2", err, n.Load())
	}
}

func TestBackoffGrowsAndIsCapped(t *testing.T) {
	c, _ := NewClient(WithAPIKey("k"))
	c.jitter = func() float64 { return 1 }
	want := []time.Duration{500 * time.Millisecond, time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 8 * time.Second}
	for i, w := range want {
		if got := c.delay(i+1, nil); got != w {
			t.Errorf("retry %d: %v, want %v", i+1, got, w)
		}
	}
	c.jitter = func() float64 { return 0 }
	if got := c.delay(1, nil); got != 250*time.Millisecond {
		t.Errorf("with no jitter added: %v, want half of the step", got)
	}
	if got := c.delay(60, nil); got != 4*time.Second {
		t.Errorf("a late retry: %v, want half of the cap", got)
	}
}

func TestRetryAfterIsHonoured(t *testing.T) {
	var n atomic.Int32
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		n.Add(1)
		if n.Load() == 1 {
			w.Header().Set("Retry-After", "7")
			writeJSON(w, 429, apiErrorBody("rate_limited", "Too many attempts; try again in a minute.", ""))
			return
		}
		writeJSON(w, 200, `{"runs":[],"next_cursor":null}`)
	})
	if _, err := c.Runs.List(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if len(rec.delays) != 1 || rec.delays[0] != 7*time.Second {
		t.Errorf("delays = %v, want [7s]", rec.delays)
	}
	if got := retryAfter(time.Now().Add(90 * time.Second).UTC().Format(http.TimeFormat)); got < 80*time.Second || got > 91*time.Second {
		t.Errorf("an HTTP date: %v", got)
	}
	if retryAfter("soon") != 0 || retryAfter("-3") != 0 || retryAfter("") != 0 {
		t.Error("a Retry-After that says nothing should be zero")
	}
}

func TestStartingARunIsRetriedWithTheSameIdempotencyKey(t *testing.T) {
	var keys []string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		keys = append(keys, r.Header.Get("Idempotency-Key"))
		switch len(keys) {
		case 1:
			writeJSON(w, 500, apiErrorBody("internal", "x", ""))
		case 2:
			hijackAndClose(t, w)
		default:
			writeJSON(w, 202, runJSON("run_1", "queued", "", 0))
		}
	})
	run, err := c.Scrapers.Run(context.Background(), "ats.board", nil, nil)
	if err != nil || run.ID != "run_1" {
		t.Fatalf("run = %v, err = %v", run, err)
	}
	if len(keys) != 3 || len(keys[0]) != 32 || keys[0] != keys[1] || keys[1] != keys[2] {
		t.Errorf("keys = %v; want one 32-character key sent three times", keys)
	}
}

func TestEachCallHasItsOwnKeyUnlessOneIsGiven(t *testing.T) {
	var keys []string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		keys = append(keys, r.Header.Get("Idempotency-Key"))
		if len(keys) == 4 {
			writeJSON(w, 409, apiErrorBody("idempotency_conflict", "Used.", ""))
			return
		}
		writeJSON(w, 202, runJSON("run_1", "queued", "", 0))
	})
	ctx := context.Background()
	_, _ = c.Scrapers.Run(ctx, "ats.board", nil, nil)
	_, _ = c.Scrapers.Run(ctx, "ats.board", nil, &RunParams{Wait: time.Second})
	_, _ = c.Scrapers.Run(ctx, "ats.board", nil, &RunParams{IdempotencyKey: "import-acme"})
	_, err := c.Scrapers.Run(ctx, "ats.board", nil, &RunParams{IdempotencyKey: "import-acme"})
	if keys[0] == keys[1] || keys[2] != "import-acme" {
		t.Errorf("keys = %v", keys)
	}
	if !errors.Is(err, ErrConflict) || len(keys) != 4 {
		t.Errorf("a 409 should be ErrConflict and not retried: %v after %d requests", err, len(keys))
	}
}

func TestTooManyRunsIsNotRetriedWhenStartingARun(t *testing.T) {
	var n atomic.Int32
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		n.Add(1)
		writeJSON(w, 429, apiErrorBody("too_many_runs", "x", ""))
	})
	if _, err := c.Scrapers.Run(context.Background(), "ats.board", nil, nil); !errors.Is(err, ErrTooManyRuns) {
		t.Fatalf("err = %v", err)
	}
	if n.Load() != 1 || len(rec.delays) != 0 {
		t.Errorf("%d attempts, want exactly 1: it clears when a run ends, not after a pause", n.Load())
	}
}

func TestConnectionErrorsAreRetriedForReads(t *testing.T) {
	var n atomic.Int32
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		n.Add(1)
		if n.Load() == 1 {
			hijackAndClose(t, w)
			return
		}
		writeJSON(w, 200, runJSON("run_1", "running", "", 0))
	})
	run, err := c.Runs.Get(context.Background(), "run_1")
	if err != nil || run.Status != RunRunning || n.Load() != 2 || len(rec.delays) != 1 {
		t.Fatalf("run %v, err %v, %d attempts", run, err, n.Load())
	}
}

func hijackAndClose(t *testing.T, w http.ResponseWriter) {
	t.Helper()
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}

func TestCancelIsRetried(t *testing.T) {
	var n atomic.Int32
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		n.Add(1)
		if n.Load() == 1 {
			writeJSON(w, 500, apiErrorBody("internal", "x", ""))
			return
		}
		writeJSON(w, 200, runJSON("run_1", "cancelled", "cancelled", 0))
	})
	run, err := c.Runs.Cancel(context.Background(), "run_1")
	if err != nil || run.Status != RunCancelled || n.Load() != 2 {
		t.Fatalf("run %v, err %v, %d attempts", run, err, n.Load())
	}
}

func TestClientErrorsAreNotRetried(t *testing.T) {
	for _, tc := range []struct {
		status int
		code   string
	}{{401, "unauthorized"}, {403, "forbidden"}, {404, "not_found"}, {422, "invalid_input"}, {402, "quota_exceeded"}, {429, "too_many_runs"}} {
		var n atomic.Int32
		c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			n.Add(1)
			writeJSON(w, tc.status, apiErrorBody(tc.code, "x", ""))
		})
		if _, err := c.Runs.Get(context.Background(), "run_1"); err == nil || n.Load() != 1 {
			t.Errorf("%s: err %v after %d attempts; want 1 attempt", tc.code, err, n.Load())
		}
	}
}

func TestContextCancelsARetryWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		cancel()
		writeJSON(w, 500, apiErrorBody("internal", "x", ""))
	})
	c.sleep = sleepContext
	_, err := c.Runs.Get(ctx, "run_1")
	if err == nil {
		t.Fatal("want an error")
	}
}

func TestTimeout(t *testing.T) {
	release := make(chan struct{})
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) { <-release }, WithTimeout(30*time.Millisecond), WithMaxRetries(0))
	defer close(release)
	_, err := c.Runs.Get(context.Background(), "run_1")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got %v, want a deadline error", err)
	}
}

func TestWaitExtendsTheTimeout(t *testing.T) {
	// The attempt may take the timeout plus the wait that was asked for: 20ms alone would not be enough here.
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(120 * time.Millisecond)
		writeJSON(w, 200, runJSON("run_1", "ok", "", 0))
	}, WithTimeout(20*time.Millisecond))
	if _, err := c.Scrapers.Run(context.Background(), "ats.board", nil, &RunParams{Wait: 5 * time.Second}); err != nil {
		t.Fatal(err)
	}
}

func TestRecordsIteratorPages(t *testing.T) {
	c, rec := newTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Query().Get("cursor") {
		case "":
			writeJSON(w, 200, `{"records":[{"title":"A"},{"title":"B"}],"next_cursor":"Mg"}`)
		case "Mg":
			writeJSON(w, 200, `{"records":[{"title":"C"}],"next_cursor":"Mw"}`)
		case "Mw":
			writeJSON(w, 200, `{"records":[],"next_cursor":null}`)
		}
	})
	type job struct {
		Title string `json:"title"`
	}
	it := c.Runs.Records(context.Background(), "run_1", &RecordListParams{Limit: 2})
	var titles []string
	for it.Next() {
		j, err := Decode[job](it.Current())
		if err != nil {
			t.Fatal(err)
		}
		titles = append(titles, j.Title)
	}
	if it.Err() != nil || strings.Join(titles, "") != "ABC" {
		t.Fatalf("titles %v, err %v", titles, it.Err())
	}
	if it.Next() {
		t.Error("Next after the end should stay false")
	}
	calls := rec.calls()
	if len(calls) != 3 || calls[0].Query != "limit=2" || calls[1].Query != "cursor=Mg&limit=2" || calls[2].Query != "cursor=Mw&limit=2" {
		t.Errorf("requests = %+v", calls)
	}

	all, err := DecodeAll[map[string]any](c.Runs.Records(context.Background(), "run_1", nil))
	if err != nil || len(all) != 3 || all[2]["title"] != "C" {
		t.Errorf("DecodeAll = %v, %v", all, err)
	}
}

func TestIteratorAll(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("cursor") == "" {
			writeJSON(w, 200, `{"records":[1,2],"next_cursor":"x"}`)
			return
		}
		writeJSON(w, 404, apiErrorBody("not_found", "gone", ""))
	})
	// All is an iter.Seq2; it is called by hand here so that the test also builds with Go 1.22.
	var got []string
	var last error
	c.Runs.Records(context.Background(), "run_1", nil).All()(func(record json.RawMessage, err error) bool {
		if err != nil {
			last = err
			return false
		}
		got = append(got, string(record))
		return true
	})
	if strings.Join(got, ",") != "1,2" || !errors.Is(last, ErrNotFound) {
		t.Errorf("got %v, last error %v", got, last)
	}

	stopped := 0
	c.Runs.Records(context.Background(), "run_1", nil).All()(func(json.RawMessage, error) bool { stopped++; return false })
	if stopped != 1 {
		t.Errorf("a yield that returns false must stop the iteration; called %d times", stopped)
	}
}

func TestRecordsExpired(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, `{"records":[],"next_cursor":null,"expired":true}`)
	})
	page, err := c.Runs.RecordsPage(context.Background(), "run_1", nil)
	if err != nil || !page.Expired || len(page.Records) != 0 || page.NextCursor != "" {
		t.Fatalf("page %+v, err %v", page, err)
	}
	it := c.Runs.Records(context.Background(), "run_1", nil)
	if it.Next() || !errors.Is(it.Err(), ErrRecordsExpired) {
		t.Errorf("the iterator should stop with ErrRecordsExpired, got %v", it.Err())
	}
}

func TestRunsListAll(t *testing.T) {
	c, rec := newTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("cursor") == "" {
			writeJSON(w, 200, `{"runs":[`+runJSON("run_1", "ok", "", 1)+`,`+runJSON("run_2", "failed", "fetch_failed", 0)+`],"next_cursor":"p2"}`)
			return
		}
		writeJSON(w, 200, `{"runs":[`+runJSON("run_3", "running", "", 0)+`],"next_cursor":null}`)
	})
	runs, err := c.Runs.ListAll(context.Background(), &RunListParams{ScraperID: "ats.board"}).Collect()
	if err != nil || len(runs) != 3 || runs[1].Reason != "fetch_failed" || runs[2].ID != "run_3" {
		t.Fatalf("runs %+v, err %v", runs, err)
	}
	if calls := rec.calls(); len(calls) != 2 || calls[1].Query != "cursor=p2&scraper_id=ats.board" {
		t.Errorf("requests = %+v", calls)
	}
}

func TestWaitPollsUntilTheRunEnds(t *testing.T) {
	statuses := []string{"queued", "running", "running", "partial"}
	var n atomic.Int32
	c, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, runJSON("run_1", statuses[n.Load()], map[bool]string{true: "render_required"}[statuses[n.Load()] == "partial"], 4))
		n.Add(1)
	})
	var seen []RunStatus
	run, err := c.Runs.Wait(context.Background(), "run_1", &WaitParams{PollInterval: 3 * time.Second, OnPoll: func(r *Run) { seen = append(seen, r.Status) }})
	if err != nil || run.Status != RunPartial || run.Reason != "render_required" {
		t.Fatalf("run %+v, err %v", run, err)
	}
	if n.Load() != 4 || len(seen) != 4 || seen[0] != RunQueued {
		t.Errorf("%d polls, seen %v", n.Load(), seen)
	}
	if len(rec.delays) != 3 || rec.delays[0] != 3*time.Second {
		t.Errorf("pauses = %v, want three of 3s", rec.delays)
	}
	if run.Err() != nil {
		t.Errorf("partial is not an error: %v", run.Err())
	}
}

func TestWaitStopsWithTheContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 200, runJSON("run_1", "running", "", 0))
	})
	run, err := c.Runs.Wait(ctx, "run_1", &WaitParams{OnPoll: func(*Run) { cancel() }})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
	if run == nil || run.ID != "run_1" {
		t.Errorf("the run as last seen should come back with the error, got %v", run)
	}
}

func TestScrape(t *testing.T) {
	c, rec := newTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		switch {
		case req.Method == "POST" && req.URL.Path == "/v1/scrapers/ats.board/runs":
			writeJSON(w, 202, runJSON("run_9", "running", "", 0))
		case req.URL.Path == "/v1/runs/run_9":
			writeJSON(w, 200, runJSON("run_9", "ok", "", 2))
		case req.URL.Path == "/v1/runs/run_9/records":
			writeJSON(w, 200, `{"records":[{"title":"A"},{"title":"B"}],"next_cursor":null}`)
		default:
			t.Errorf("unexpected %s %s", req.Method, req.URL.Path)
		}
	})
	run, records, err := c.Scrape(context.Background(), "ats.board", map[string]any{"board_url": "https://boards.greenhouse.io/acme"})
	if err != nil {
		t.Fatal(err)
	}
	if run.ID != "run_9" || run.Status != RunOK || run.RecordCount != 2 {
		t.Errorf("run = %+v", run)
	}
	all, err := records.Collect()
	if err != nil || len(all) != 2 {
		t.Fatalf("records %v, err %v", all, err)
	}
	calls := rec.calls()
	if len(calls) != 3 || calls[0].Query != "wait=60" || calls[0].Body != `{"input":{"board_url":"https://boards.greenhouse.io/acme"}}` {
		t.Errorf("requests = %+v", calls)
	}
	if calls[2].Query != "limit=1000" {
		t.Errorf("the records are read in pages of 1000, got ?%s", calls[2].Query)
	}
}

func TestScrapeThatEndsWithinTheWait(t *testing.T) {
	c, rec := newTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		writeJSON(w, 200, runJSON("run_9", "ok", "", 0))
	})
	run, records, err := c.Scrape(context.Background(), "ats.board", nil)
	if err != nil || run.Status != RunOK || records == nil {
		t.Fatalf("run %v, err %v", run, err)
	}
	if n := len(rec.calls()); n != 1 {
		t.Errorf("%d requests; a run that came back ended needs no poll", n)
	}
}

func TestScrapeRunFailed(t *testing.T) {
	for _, tc := range []struct{ status, reason string }{{"failed", "fetch_failed"}, {"unsupported", "render_required"}, {"cancelled", "cancelled"}} {
		t.Run(tc.status, func(t *testing.T) {
			c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				writeJSON(w, 200, runJSON("run_9", tc.status, tc.reason, 0))
			})
			run, records, err := c.Scrape(context.Background(), "ats.board", nil)
			var failed *RunFailedError
			if !errors.As(err, &failed) {
				t.Fatalf("got %T %v, want *RunFailedError", err, err)
			}
			if string(failed.Status()) != tc.status || failed.Reason() != tc.reason || failed.Run.ID != "run_9" {
				t.Errorf("error = %+v", failed.Run)
			}
			if !errors.Is(err, ErrRunFailed) {
				t.Error("errors.Is(err, ErrRunFailed) should be true")
			}
			if want := "run run_9 ended " + tc.status + " (reason: " + tc.reason + ")"; !strings.Contains(err.Error(), want) {
				t.Errorf("Error() = %q, should contain %q", err.Error(), want)
			}
			if run == nil || records != nil {
				t.Errorf("the run comes back with the error, the records do not: %v %v", run, records)
			}
		})
	}
}

func TestScrapeAPIError(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, 402, apiErrorBody("quota_exceeded", "This month's units are used up.", `{"used_units":1,"allowance_units":1}`))
	})
	run, records, err := c.Scrape(context.Background(), "ats.board", nil)
	if !errors.Is(err, ErrQuotaExceeded) || run != nil || records != nil {
		t.Fatalf("run %v, err %v", run, err)
	}
	var api *APIError
	errors.As(err, &api)
	if api.Details["allowance_units"] != float64(1) {
		t.Errorf("details = %v", api.Details)
	}
}

func TestDataDecoding(t *testing.T) {
	const posting = `{"id":"job_2a8c2e32f38fa96099cc53b2","title":"Voice of the Customer Program Manager","title_clean":"Voice of the Customer Program Manager",` +
		`"company":{"domain":"figma.com","name":"Figma"},"location":{"raw":"San Francisco, CA","city":"San Francisco","region":"CA","country":"United States","country_code":null,"remote":false},` +
		`"employment_type":null,"salary":{"raw":"$100k","min":100000,"max":null,"currency":"USD","period":"year"},"posted_at":"2026-06-04","valid_through":null,"seniority":"manager","function":null,` +
		`"url":"https://example.com/jobs","apply_url":"https://boards.greenhouse.io/figma/jobs/6008343004","source":{"name":"greenhouse","provenance":"feed","trusted":true,"external_id":"6008343004"},` +
		`"status":"open","first_seen_at":"2026-11-10T22:37:56Z","last_seen_at":"2026-11-16T22:38:01.152673Z","closed_at":null,"description":"Figma is growing.",` +
		`"classification":{"seniority_rule":"manager","function_rule":null},"history":{"events":[{"kind":"first_seen","at":"2026-11-10T22:37:56Z"}],` +
		`"changes":[{"field":"title","old":"A","new":"B","at":"2026-11-12T00:00:00Z"}],"sightings":{"count":3,"recent":["2026-11-16T22:38:01.152673Z"]},"reopened_count":1}}`
	const company = `{"domain":"figma.com","name":"Figma","names":[{"name":"Figma","first_seen_at":"2026-11-10T22:37:56Z","last_seen_at":"2026-11-16T22:38:01.152673Z"}],` +
		`"careers_url":"https://boards.greenhouse.io/figma","ats":"greenhouse","first_seen_at":"2026-11-10T22:37:56Z","last_seen_at":"2026-11-16T22:38:01.152673Z",` +
		`"last_read":{"at":"2026-11-16T22:38:01.152673Z","status":"complete","last_complete_at":"2026-11-16T22:38:01.152673Z"},"postings":{"open":140,"closed":0}}`
	c, _ := newTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/data/postings/job_2a8c2e32f38fa96099cc53b2":
			writeJSON(w, 200, posting)
		case "/v1/data/postings":
			writeJSON(w, 200, `{"postings":[`+posting+`],"next_cursor":"abc"}`)
		case "/v1/data/companies":
			writeJSON(w, 200, `{"companies":[`+company+`],"next_cursor":null}`)
		case "/v1/data/companies/figma.com":
			writeJSON(w, 200, company)
		}
	})
	ctx := context.Background()
	p, err := c.Data.Postings.Get(ctx, "job_2a8c2e32f38fa96099cc53b2")
	if err != nil {
		t.Fatal(err)
	}
	if p.Title == "" || *p.Company.Domain != "figma.com" || p.Location.CountryCode != nil || *p.Location.Remote || p.EmploymentType != nil ||
		*p.Salary.Min != 100000 || p.Salary.Max != nil || *p.PostedAt != "2026-06-04" || !p.Source.Trusted || p.ClosedAt != nil ||
		*p.Description != "Figma is growing." || *p.Classification.SeniorityRule != "manager" || p.History.ReopenedCount != 1 ||
		p.History.Events[0].Kind != "first_seen" || string(p.History.Changes[0].New) != `"B"` || p.History.Sightings.Count != 3 {
		t.Errorf("posting = %+v", p)
	}
	if string(p.RawJSON()) != posting {
		t.Error("RawJSON should be the body as sent")
	}
	page, err := c.Data.Postings.List(ctx, nil)
	if err != nil || len(page.Postings) != 1 || page.NextCursor != "abc" || page.Postings[0].ID != p.ID {
		t.Errorf("page = %+v, %v", page, err)
	}
	companies, err := c.Data.Companies.List(ctx, nil)
	if err != nil || len(companies.Companies) != 1 || companies.NextCursor != "" {
		t.Fatalf("companies = %+v, %v", companies, err)
	}
	co, err := c.Data.Companies.Get(ctx, "figma.com")
	if err != nil || co.Domain != "figma.com" || *co.Name != "Figma" || co.Postings.Open != 140 || *co.LastRead.Status != "complete" || len(co.Names) != 1 {
		t.Errorf("company = %+v, %v", co, err)
	}
}

func TestUsageAndScraperDecoding(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/v1/usage":
			writeJSON(w, 200, `{"period":{"from":"2026-11-01","to":"2026-11-16"},"total_units":13,"allowance_units":1000,"rows":[{"key":"data.postings.list","units":7,"runs":5},{"key":"ats.board","units":1,"runs":1}]}`)
		case "/v1/scrapers":
			writeJSON(w, 200, `{"scrapers":[{"id":"ats.board","version":1,"title":"ATS job board","summary":"Job openings.","mode":"job","billing_unit":"request","data_class":"public","availability":"public"}]}`)
		case "/v1/scrapers/ats.board":
			writeJSON(w, 200, `{"id":"ats.board","version":1,"title":"ATS job board","summary":"Job openings.","mode":"job","billing_unit":"request","data_class":"public","availability":"public","input_schema":{"type":"object","required":["board_url"]},"record_schema":{"type":"object"}}`)
		}
	})
	ctx := context.Background()
	u, err := c.Usage.Get(ctx, nil)
	if err != nil || u.Period.From != "2026-11-01" || u.TotalUnits != 13 || u.AllowanceUnits != 1000 || len(u.Rows) != 2 || u.Rows[0].Runs != 5 {
		t.Errorf("usage = %+v, %v", u, err)
	}
	list, err := c.Scrapers.List(ctx)
	if err != nil || len(list) != 1 || list[0].BillingUnit != "request" || list[0].Mode != "job" {
		t.Errorf("scrapers = %+v, %v", list, err)
	}
	s, err := c.Scrapers.Get(ctx, "ats.board")
	if err != nil || s.ID != "ats.board" || s.Title != "ATS job board" || string(s.InputSchema) != `{"type":"object","required":["board_url"]}` {
		t.Errorf("scraper = %+v, %v", s, err)
	}
}

func TestRunStatus(t *testing.T) {
	for status, want := range map[RunStatus][2]bool{
		RunQueued: {false, false}, RunRunning: {false, false}, RunOK: {true, true}, RunPartial: {true, true},
		RunFailed: {true, false}, RunUnsupported: {true, false}, RunCancelled: {true, false},
	} {
		if status.IsFinal() != want[0] || status.IsSuccess() != want[1] {
			t.Errorf("%s: final %v, success %v", status, status.IsFinal(), status.IsSuccess())
		}
	}
	if (&Run{Status: RunRunning}).Err() != nil {
		t.Error("a run that is still going has no error")
	}
}
