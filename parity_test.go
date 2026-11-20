package enrichednow

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func setStanding(w http.ResponseWriter, remaining string) {
	h := w.Header()
	h.Set("X-Units-Limit", "1000")
	h.Set("X-Units-Used", "120")
	h.Set("X-Units-Remaining", remaining)
	h.Set("X-Units-Topup", "50")
	h.Set("X-Units-Reset", "2026-12-01T00:00:00Z")
	h.Set("X-RateLimit-Limit", "2")
	h.Set("X-RateLimit-Remaining", "1")
}

func TestStandingIsKeptFromTheLastAnswerThatCarriedIt(t *testing.T) {
	calls := 0
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch calls {
		case 1:
			writeJSON(w, 200, `{"scrapers":[]}`)
		case 2:
			setStanding(w, "880")
			writeJSON(w, 200, runJSON("run_1", "ok", "", 1))
		default:
			setStanding(w, "0")
			writeJSON(w, 402, apiErrorBody("quota_exceeded", "No.", ""))
		}
	})
	ctx := context.Background()
	if _, err := client.Scrapers.List(ctx); err != nil || client.Standing() != nil {
		t.Fatalf("the catalogue says nothing of an account: %v %+v", err, client.Standing())
	}
	if _, err := client.Runs.Get(ctx, "run_1"); err != nil {
		t.Fatal(err)
	}
	s := client.Standing()
	if s == nil || *s.UnitsLimit != 1000 || *s.UnitsUsed != 120 || *s.UnitsRemaining != 880 || *s.UnitsTopup != 50 ||
		*s.RunsLimit != 2 || *s.RunsRemaining != 1 || s.Unmetered || !s.UnitsReset.Equal(time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("standing: %+v", s)
	}
	if s.UnitsAvailable() != 930 {
		t.Errorf("available: %d", s.UnitsAvailable())
	}
	if _, err := client.Scrapers.Run(ctx, "ats.board", nil, nil); !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("got %v", err)
	}
	if client.Standing().UnitsAvailable() != 50 {
		t.Errorf("an error carries the standing too: %+v", client.Standing())
	}
}

func TestUsageStandingAsksForIt(t *testing.T) {
	client, rec := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		setStanding(w, "880")
		w.Header().Set("X-Units-Unmetered", "true")
		writeJSON(w, 200, `{"runs":[],"next_cursor":null}`)
	})
	s, err := client.Usage.Standing(context.Background())
	if err != nil || !s.Unmetered {
		t.Fatalf("%v %+v", err, s)
	}
	if c := rec.calls()[0]; c.Path != "/v1/runs" || c.Query != "limit=1" {
		t.Errorf("sent %+v", c)
	}
}

func TestUsageSaysUnmetered(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, `{"period":{"from":"2026-11-01","to":"2026-11-19"},"total_units":3,"allowance_units":1000,"unmetered":true,"rows":[]}`)
	})
	u, err := client.Usage.Get(context.Background(), nil)
	if err != nil || !u.Unmetered {
		t.Fatalf("%v %+v", err, u)
	}
}

func TestWaitUpTo300SecondsAndEveryRunFilter(t *testing.T) {
	client, rec := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/runs" {
			writeJSON(w, 200, `{"runs":[],"next_cursor":null}`)
			return
		}
		writeJSON(w, 200, runJSON("run_1", "ok", "", 1))
	})
	ctx := context.Background()
	if _, err := client.Runs.GetWait(ctx, "run_1", 120*time.Second); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Scrapers.Run(ctx, "ats.board", map[string]any{}, &RunParams{Wait: 300 * time.Second, WebhookURL: "https://example.com/hook"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Scrapers.Run(ctx, "ats.board", nil, &RunParams{Wait: 301 * time.Second}); err == nil {
		t.Error("a wait the API would refuse should not be sent")
	}
	if _, err := client.Runs.GetWait(ctx, "run_1", -time.Second); err == nil {
		t.Error("a negative wait should not be sent")
	}
	from := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	if _, err := client.Runs.List(ctx, &RunListParams{APIKeyID: "key_1", From: from, To: from.AddDate(0, 0, 18), Status: RunOK}); err != nil {
		t.Fatal(err)
	}
	calls := rec.calls()
	if len(calls) != 3 {
		t.Fatalf("sent %d requests: %+v", len(calls), calls)
	}
	if calls[0].Query != "wait=120" || calls[1].Query != "wait=300" || !strings.Contains(calls[1].Body, `"webhook_url":"https://example.com/hook"`) {
		t.Errorf("sent %+v", calls[:2])
	}
	q, _ := url.ParseQuery(calls[2].Query)
	if q.Get("api_key_id") != "key_1" || q.Get("from") != "2026-11-01T00:00:00Z" || q.Get("to") != "2026-11-19T00:00:00Z" || q.Get("status") != "ok" {
		t.Errorf("list sent %v", q)
	}
}

func TestExportOnePageAndAWholeRun(t *testing.T) {
	pages := map[string]string{"": "a,b\r\n1,2\r\n", "c2": "a,b\r\n3,4\r\n", "c3": "a,b\r\n5,6\r\n"}
	next := map[string]string{"": "c2", "c2": "c3"}
	var accepts []string
	client, rec := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		accepts = append(accepts, r.Header.Get("Accept"))
		cursor := r.URL.Query().Get("cursor")
		if n := next[cursor]; n != "" {
			w.Header().Set("X-Next-Cursor", n)
		}
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		io.WriteString(w, pages[cursor])
	})
	ctx := context.Background()
	page, err := client.Runs.Export(ctx, "run_1", ExportParams{Format: FormatCSV, Limit: 1})
	if err != nil || string(page.Body) != "a,b\r\n1,2\r\n" || page.NextCursor != "c2" || page.Expired {
		t.Fatalf("%v %+v", err, page)
	}
	if q := rec.calls()[0].Query; q != "format=csv&limit=1" || accepts[0] != "text/csv" {
		t.Errorf("sent %q accept %q", q, accepts[0])
	}
	var out bytes.Buffer
	if err := client.Runs.ExportAll(ctx, "run_1", FormatCSV, &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != "a,b\r\n1,2\r\n3,4\r\n5,6\r\n" {
		t.Errorf("the whole run: %q", out.String())
	}
	if _, err := client.Runs.Export(ctx, "run_1", ExportParams{Format: "xml"}); err == nil {
		t.Error("a format the API does not write should be refused")
	}
}

func TestExportOfRecordsThatAreGone(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Records-Expired", "true")
		w.WriteHeader(200)
	})
	page, err := client.Runs.Export(context.Background(), "run_1", ExportParams{Format: FormatNDJSON})
	if err != nil || !page.Expired {
		t.Fatalf("%v %+v", err, page)
	}
	if err := client.Runs.ExportAll(context.Background(), "run_1", FormatNDJSON, io.Discard); !errors.Is(err, ErrRecordsExpired) {
		t.Errorf("got %v", err)
	}
}

func TestAFileARunProduced(t *testing.T) {
	client, rec := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte{0x89, 'P', 'N', 'G'})
	})
	got, err := client.Runs.File(context.Background(), "run_1", "screenshot.png")
	if err != nil || string(got) != "\x89PNG" || rec.calls()[0].Path != "/v1/runs/run_1/files/screenshot.png" {
		t.Fatalf("%v %q %+v", err, got, rec.calls())
	}
}

const endpointJSON = `{"id":"whk_1","url":"https://example.com/hook","description":"prod","events":["run.succeeded"],"enabled":true,` +
	`"disabled_reason":null,"disabled_at":null,"last_success_at":null,"last_failure_at":null,"created_at":"2026-11-19T10:00:00Z"}`

func TestWebhookEndpoints(t *testing.T) {
	client, rec := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/webhooks":
			writeJSON(w, 200, `{"endpoints":[`+endpointJSON+`],"run_secret_set":true}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/webhooks":
			writeJSON(w, 201, strings.TrimSuffix(endpointJSON, "}")+`,"secret":"whsec_new"}`)
		case r.Method == http.MethodPatch, r.Method == http.MethodGet:
			writeJSON(w, 200, endpointJSON)
		case r.Method == http.MethodDelete:
			w.WriteHeader(204)
		case strings.HasSuffix(r.URL.Path, "/rotate-secret"):
			writeJSON(w, 200, strings.TrimSuffix(endpointJSON, "}")+`,"secret":"whsec_rotated"}`)
		case r.URL.Path == "/v1/webhooks/run-secret":
			writeJSON(w, 200, `{"secret":"whsec_run"}`)
		}
	})
	ctx := context.Background()
	list, err := client.Webhooks.List(ctx)
	if err != nil || len(list.Endpoints) != 1 || !list.RunSecretSet || list.Endpoints[0].Events[0] != EventRunSucceeded {
		t.Fatalf("%v %+v", err, list)
	}
	made, err := client.Webhooks.Create(ctx, WebhookEndpointParams{URL: "https://example.com/hook", Events: []string{EventRunSucceeded}, Description: "prod"})
	if err != nil || made.Secret != "whsec_new" || made.ID != "whk_1" {
		t.Fatalf("%v %+v", err, made)
	}
	off := false
	if _, err := client.Webhooks.Update(ctx, "whk_1", WebhookEndpointChange{Enabled: &off}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Webhooks.Update(ctx, "whk_1", WebhookEndpointChange{}); err == nil {
		t.Error("an empty change should not be sent")
	}
	if _, err := client.Webhooks.Get(ctx, "whk_1"); err != nil {
		t.Fatal(err)
	}
	if rotated, err := client.Webhooks.RotateSecret(ctx, "whk_1"); err != nil || rotated.Secret != "whsec_rotated" {
		t.Fatalf("%v %+v", err, rotated)
	}
	if secret, err := client.Webhooks.CreateRunSecret(ctx); err != nil || secret.Secret != "whsec_run" {
		t.Fatalf("%v %+v", err, secret)
	}
	if err := client.Webhooks.Delete(ctx, "whk_1"); err != nil {
		t.Fatal(err)
	}
	calls := rec.calls()
	if calls[1].Body != `{"url":"https://example.com/hook","events":["run.succeeded"],"description":"prod"}` || calls[2].Body != `{"enabled":false}` {
		t.Errorf("sent %q and %q", calls[1].Body, calls[2].Body)
	}
}

func TestCreatingAnEndpointIsNotSentAgain(t *testing.T) {
	client, rec := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 500, apiErrorBody("internal", "Oops.", ""))
	})
	if _, err := client.Webhooks.Create(context.Background(), WebhookEndpointParams{URL: "https://x", Events: WebhookEvents}); !errors.Is(err, ErrServer) {
		t.Fatalf("got %v", err)
	}
	if n := len(rec.calls()); n != 1 {
		t.Errorf("sent %d times", n)
	}
}

func TestWebhookDeliveries(t *testing.T) {
	client, rec := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		delivery := `{"id":"whd_1","endpoint_id":"whk_1","run_id":"run_1","event_id":"evt_1","event":"run.succeeded","url":"https://x",` +
			`"status":"failed","attempts":3,"max_attempts":8,"next_attempt_at":null,"response_status":500,"error":"HTTP 500",` +
			`"duration_ms":12,"last_attempt_at":"2026-11-19T10:00:00Z","created_at":"2026-11-19T10:00:00Z","finished_at":null,"resend_of":null}`
		if r.Method == http.MethodPost {
			writeJSON(w, 201, delivery)
			return
		}
		if r.URL.Query().Get("cursor") == "" {
			writeJSON(w, 200, `{"deliveries":[`+delivery+`],"next_cursor":"c2"}`)
			return
		}
		writeJSON(w, 200, `{"deliveries":[`+delivery+`],"next_cursor":null}`)
	})
	ctx := context.Background()
	all, err := client.Webhooks.Deliveries.ListAll(ctx, &WebhookDeliveryListParams{Status: "failed", RunID: "run_1"}).Collect()
	if err != nil || len(all) != 2 || *all[0].ResponseStatus != 500 {
		t.Fatalf("%v %+v", err, all)
	}
	if q := rec.calls()[0].Query; q != "run_id=run_1&status=failed" {
		t.Errorf("sent %q", q)
	}
	if again, err := client.Webhooks.Deliveries.Resend(ctx, "whd_1"); err != nil || again.EventID != "evt_1" {
		t.Fatalf("%v %+v", err, again)
	}
}

func sign(secret, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "."))
	mac.Write(body)
	return "v1=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyAndParseWebhook(t *testing.T) {
	secret := "whsec_test"
	body := []byte(`{"id":"evt_1","type":"run.succeeded","created_at":"2026-11-19T10:00:00Z","data":{"run":` + runJSON("run_1", "ok", "", 2) + `}}`)
	now := strconv.FormatInt(time.Now().Unix(), 10)
	header := http.Header{}
	header.Set(HeaderWebhookTimestamp, now)
	header.Set(HeaderWebhookSignature, sign(secret, now, body))
	event, err := ParseWebhook(secret, header, body)
	if err != nil || event.ID != "evt_1" || event.Type != EventRunSucceeded || event.Data.Run.ID != "run_1" || event.Data.Run.RecordCount != 2 {
		t.Fatalf("%v %+v", err, event)
	}
	if event.Data.Run.UnitsPerRequest != nil || event.Data.Run.UnitsPerRecord != nil {
		t.Errorf("an event from before costs differed has no rate: %+v", event.Data.Run)
	}
	rated := []byte(strings.Replace(string(body), `"units":0,`, `"units":0,"units_per_request":null,"units_per_record":3,`, 1))
	header.Set(HeaderWebhookSignature, sign(secret, now, rated))
	event, err = ParseWebhook(secret, header, rated)
	if err != nil || event.Data.Run.UnitsPerRecord == nil || *event.Data.Run.UnitsPerRecord != 3 || event.Data.Run.UnitsPerRequest != nil {
		t.Fatalf("%v %+v", err, event)
	}
	header.Set(HeaderWebhookSignature, sign(secret, now, body))
	old := strconv.FormatInt(time.Now().Unix()-301, 10)
	cases := map[string][4]string{
		"another secret":       {"whsec_other", now, string(body), sign(secret, now, body)},
		"a changed body":       {secret, now, string(body) + " ", sign(secret, now, body)},
		"too old":              {secret, old, string(body), sign(secret, old, body)},
		"upper-case hex":       {secret, now, string(body), strings.ToUpper(sign(secret, now, body))},
		"no signature":         {secret, now, string(body), ""},
		"a timestamp of words": {secret, "soon", string(body), sign(secret, "soon", body)},
		"a signed timestamp":   {secret, "+" + now, string(body), sign(secret, "+"+now, body)},
		"no secret":            {"", now, string(body), sign("", now, body)},
	}
	for name, c := range cases {
		if err := VerifyWebhookSignature(c[0], c[1], []byte(c[2]), c[3], 0); !errors.Is(err, ErrWebhookVerification) {
			t.Errorf("%s: got %v", name, err)
		}
	}
	if err := VerifyWebhookSignature(secret, old, body, sign(secret, old, body), 400*time.Second); err != nil {
		t.Errorf("a wider tolerance should accept it: %v", err)
	}
	notEvent := []byte(`[1,2]`)
	header.Set(HeaderWebhookSignature, sign(secret, now, notEvent))
	if _, err := ParseWebhook(secret, header, notEvent); !errors.Is(err, ErrWebhookVerification) {
		t.Errorf("a body that is not an event: %v", err)
	}
}
