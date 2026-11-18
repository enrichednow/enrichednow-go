package enrichednow

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ScraperService is the catalogue, and where runs are started.
type ScraperService struct{ client *Client }

// List returns the catalogue. It needs no API key.
func (s *ScraperService) List(ctx context.Context) ([]Scraper, error) {
	var out struct {
		Scrapers []Scraper `json:"scrapers"`
	}
	err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/scrapers", public: true, retry: true}, &out)
	return out.Scrapers, err
}

// Get returns one scraper with the JSON Schema of its input and of its records. It needs no API key.
func (s *ScraperService) Get(ctx context.Context, scraperID string) (*ScraperDetail, error) {
	id, err := pathEscape("scraper id", scraperID)
	if err != nil {
		return nil, err
	}
	var out ScraperDetail
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/scrapers/" + id, public: true, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MaxWait is the longest the API holds a request open for a run to end.
const MaxWait = 60 * time.Second

// RunParams are the options of [ScraperService.Run].
type RunParams struct {
	// Wait asks the API to hold the request open until the run ends or this much time has passed, whichever
	// comes first: from 0 to [MaxWait]. A run that has not ended by then is returned queued or running;
	// follow it with [RunService.Wait].
	Wait time.Duration
	// IdempotencyKey names this call so that it can be repeated without starting a second run: for 24 hours the
	// same key with the same request gives the run the first call started, and with a different request fails
	// with [ErrConflict]. Any text of yours, up to 255 characters, no spaces. Empty: a new key for each call,
	// which is what makes its retries safe.
	IdempotencyKey string
}

// newIdempotencyKey is 32 random hex characters: one call's own key.
func newIdempotencyKey() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16) // no randomness to be had: unique enough for one process
	}
	return hex.EncodeToString(b[:])
}

// Run starts a run of a scraper and returns it: usually queued, or already ended when params.Wait was long
// enough. input is the scraper's input (a map or a struct that encodes to a JSON object); see the scraper's
// input schema.
//
// The request carries an Idempotency-Key, so it is safe to send again: after a connection error or a 5xx it is
// retried with the same key, and the API answers with the run the first attempt started, never a second one.
// Set [RunParams.IdempotencyKey] to make a call of yours repeatable across processes.
//
// A run that ended failed, unsupported or cancelled is returned without an error, as the API answers it; call
// [Run.Err] to turn its status into one.
func (s *ScraperService) Run(ctx context.Context, scraperID string, input any, params *RunParams) (*Run, error) {
	id, err := pathEscape("scraper id", scraperID)
	if err != nil {
		return nil, err
	}
	if input == nil {
		input = map[string]any{}
	}
	r := request{method: http.MethodPost, path: "/v1/scrapers/" + id + "/runs", body: map[string]any{"input": input}, retry: true}
	if params != nil {
		r.idempotencyKey = params.IdempotencyKey
	}
	if r.idempotencyKey == "" {
		r.idempotencyKey = newIdempotencyKey()
	}
	if params != nil && params.Wait > 0 {
		r.query = url.Values{"wait": {strconv.FormatFloat(params.Wait.Seconds(), 'f', -1, 64)}}
		r.extraTimeout = params.Wait
	}
	var out Run
	if err := s.client.do(ctx, r, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RunService reads, waits for and cancels runs, and reads their records.
type RunService struct{ client *Client }

// Get returns a run as it is now.
func (s *RunService) Get(ctx context.Context, runID string) (*Run, error) {
	id, err := pathEscape("run id", runID)
	if err != nil {
		return nil, err
	}
	var out Run
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/runs/" + id, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RunListParams filter and page [RunService.List].
type RunListParams struct {
	// Limit is the runs per page: 1 to 100. Zero leaves it to the API (20).
	Limit int
	// Cursor is the NextCursor of the page before.
	Cursor string
	// ScraperID keeps only the runs of this scraper.
	ScraperID string
	// Status keeps only the runs with this status.
	Status RunStatus
}

// List returns one page of the account's runs, newest first.
func (s *RunService) List(ctx context.Context, params *RunListParams) (*RunPage, error) {
	q := url.Values{}
	if params != nil {
		setInt(q, "limit", params.Limit)
		setString(q, "cursor", params.Cursor)
		setString(q, "scraper_id", params.ScraperID)
		setString(q, "status", string(params.Status))
	}
	var out RunPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/runs", query: q, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAll iterates over the account's runs, newest first, fetching page after page. params.Cursor, when set, is
// where it starts.
func (s *RunService) ListAll(ctx context.Context, params *RunListParams) *Iter[Run] {
	p := RunListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]Run, string, error) {
		p.Cursor = cursor
		page, err := s.List(ctx, &p)
		if err != nil {
			return nil, "", err
		}
		return page.Runs, page.NextCursor, nil
	})
}

// Cancel asks a run to stop and returns it as it is at that moment. A queued run ends cancelled at once; a
// running one stops at its next safe point, so the answer usually still says running — follow it with
// [RunService.Wait] to see how it ended. A run that has already ended is returned unchanged.
func (s *RunService) Cancel(ctx context.Context, runID string) (*Run, error) {
	id, err := pathEscape("run id", runID)
	if err != nil {
		return nil, err
	}
	var out Run
	// Cancelling twice changes nothing, so this POST is safe to retry.
	if err := s.client.do(ctx, request{method: http.MethodPost, path: "/v1/runs/" + id + "/cancel", retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DefaultPollInterval is the pause between two looks at a run in [RunService.Wait].
const DefaultPollInterval = 2 * time.Second

// WaitParams are the options of [RunService.Wait].
type WaitParams struct {
	// PollInterval is the pause between two looks at the run. Zero means [DefaultPollInterval].
	PollInterval time.Duration
	// OnPoll, when set, is called with the run every time it was read.
	OnPoll func(*Run)
}

// Wait polls a run until it has ended and returns it. How long it waits is up to the context: give it a
// deadline to give up. A run that ended failed, unsupported or cancelled is returned without an error; call
// [Run.Err] to turn its status into one.
func (s *RunService) Wait(ctx context.Context, runID string, params *WaitParams) (*Run, error) {
	interval := DefaultPollInterval
	var onPoll func(*Run)
	if params != nil {
		if params.PollInterval > 0 {
			interval = params.PollInterval
		}
		onPoll = params.OnPoll
	}
	for {
		run, err := s.Get(ctx, runID)
		if err != nil {
			return nil, err
		}
		if onPoll != nil {
			onPoll(run)
		}
		if run.Status.IsFinal() {
			return run, nil
		}
		if err := s.client.sleep(ctx, interval); err != nil {
			return run, err
		}
	}
}

// RecordListParams page [RunService.Records] and [RunService.RecordsPage].
type RecordListParams struct {
	// Limit is the records per page: 1 to 1000. Zero leaves it to the API (100).
	Limit int
	// Cursor is the NextCursor of the page before.
	Cursor string
}

// RecordsPage returns one page of a run's records.
func (s *RunService) RecordsPage(ctx context.Context, runID string, params *RecordListParams) (*RecordPage, error) {
	id, err := pathEscape("run id", runID)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	if params != nil {
		setInt(q, "limit", params.Limit)
		setString(q, "cursor", params.Cursor)
	}
	var out RecordPage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/runs/" + id + "/records", query: q, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Records iterates over a run's records in the order the run stored them, fetching page after page. Each record
// is a scraper-specific JSON object; decode it with [Decode], or all of them with [DecodeAll].
//
// Records exist for runs that ended ok or partial, and for a cancelled run that had found some. They are kept
// for 7 days after the run ended; after that the iterator stops with [ErrRecordsExpired].
func (s *RunService) Records(ctx context.Context, runID string, params *RecordListParams) *Iter[json.RawMessage] {
	p := RecordListParams{}
	if params != nil {
		p = *params
	}
	return newIter(ctx, p.Cursor, func(ctx context.Context, cursor string) ([]json.RawMessage, string, error) {
		p.Cursor = cursor
		page, err := s.RecordsPage(ctx, runID, &p)
		if err != nil {
			return nil, "", err
		}
		if page.Expired {
			return nil, "", ErrRecordsExpired
		}
		return page.Records, page.NextCursor, nil
	})
}

// Scrape runs a scraper to the end: it starts a run, waits until the run has ended and returns it with an
// iterator over its records.
//
//	run, records, err := client.Scrape(ctx, "ats.board", map[string]any{"board_url": "https://boards.greenhouse.io/acme"})
//
// A run that ended failed, unsupported or cancelled is returned with a [*RunFailedError], which carries the
// status and the reason; the run is returned too. How long Scrape waits is up to the context. If the context
// ends first, the run goes on: the run returned with the context's error has the id to follow or cancel it.
func (c *Client) Scrape(ctx context.Context, scraperID string, input any) (*Run, *Iter[json.RawMessage], error) {
	run, err := c.Scrapers.Run(ctx, scraperID, input, &RunParams{Wait: MaxWait})
	if err != nil {
		return nil, nil, err
	}
	if !run.Status.IsFinal() {
		latest, err := c.Runs.Wait(ctx, run.ID, nil)
		if latest != nil {
			run = latest
		}
		if err != nil {
			return run, nil, err
		}
	}
	if err := run.Err(); err != nil {
		return run, nil, err
	}
	return run, c.Runs.Records(ctx, run.ID, &RecordListParams{Limit: 1000}), nil
}

// UsageService reports the units used.
type UsageService struct{ client *Client }

// UsageParams choose the period and the grouping of [UsageService.Get].
type UsageParams struct {
	// From and To are the first and last day (UTC, YYYY-MM-DD), both included. Without them the period is the
	// current calendar month up to today. A period is at most 366 days.
	From string
	To   string
	// GroupBy is "day" (the default), "scraper" or "product".
	GroupBy string
	// Product keeps one product's usage: "scrapers" or "data". Empty counts both.
	Product string
}

// Get returns the units used over a period.
func (s *UsageService) Get(ctx context.Context, params *UsageParams) (*Usage, error) {
	q := url.Values{}
	if params != nil {
		setString(q, "from", params.From)
		setString(q, "to", params.To)
		setString(q, "group_by", params.GroupBy)
		setString(q, "product", params.Product)
	}
	var out Usage
	if err := s.client.do(ctx, request{method: http.MethodGet, path: "/v1/usage", query: q, retry: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func setString(q url.Values, name, value string) {
	if value != "" {
		q.Set(name, value)
	}
}

func setInt(q url.Values, name string, value int) {
	if value != 0 {
		q.Set(name, strconv.Itoa(value))
	}
}

func setTime(q url.Values, name string, value time.Time) {
	if !value.IsZero() {
		q.Set(name, value.UTC().Format(time.RFC3339))
	}
}
