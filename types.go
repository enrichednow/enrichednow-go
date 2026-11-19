package enrichednow

import (
	"encoding/json"
	"time"
)

// Health is the answer of GET /v1/health.
type Health struct {
	rawBody
	Status string `json:"status"`
}

// Scraper is one entry of the catalogue.
type Scraper struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	// Mode is "sync" or "job".
	Mode string `json:"mode"`
	// BillingUnit is "request" (a charged run costs 1 unit) or "record" (1 unit for each record).
	BillingUnit  string `json:"billing_unit"`
	DataClass    string `json:"data_class"`
	Availability string `json:"availability"`
}

// ScraperDetail is a scraper with the JSON Schema of its input and of one of its records.
type ScraperDetail struct {
	rawBody
	Scraper
	// InputSchema is the JSON Schema of a run's input.
	InputSchema json.RawMessage `json:"input_schema"`
	// RecordSchema is the JSON Schema of one record.
	RecordSchema json.RawMessage `json:"record_schema"`
}

// RunStatus is where a run stands.
type RunStatus string

// The statuses of a run. Queued and running runs are still going; the other five are final.
const (
	RunQueued      RunStatus = "queued"
	RunRunning     RunStatus = "running"
	RunOK          RunStatus = "ok"
	RunPartial     RunStatus = "partial"
	RunFailed      RunStatus = "failed"
	RunUnsupported RunStatus = "unsupported"
	RunCancelled   RunStatus = "cancelled"
)

// IsFinal reports whether a run with this status has ended.
func (s RunStatus) IsFinal() bool { return s != RunQueued && s != RunRunning && s != "" }

// IsSuccess reports whether the run ended with records to read: ok or partial.
func (s RunStatus) IsSuccess() bool { return s == RunOK || s == RunPartial }

// Run is one execution of one scraper with one input.
type Run struct {
	rawBody
	ID             string    `json:"id"`
	ScraperID      string    `json:"scraper_id"`
	ScraperVersion int       `json:"scraper_version"`
	Status         RunStatus `json:"status"`
	// Reason is a stable code when the status is not ok, e.g. "render_required" or "fetch_failed". Empty
	// when the API sent null.
	Reason string `json:"reason"`
	// Input is the input the run was started with.
	Input       json.RawMessage `json:"input"`
	RecordCount int             `json:"record_count"`
	// Units is what the run was charged: 0 until it ends ok or partial.
	Units int `json:"units"`
	// Diagnostics says what the run did; nil until it ends.
	Diagnostics *RunDiagnostics `json:"diagnostics"`
	CreatedAt   time.Time       `json:"created_at"`
	StartedAt   *time.Time      `json:"started_at"`
	FinishedAt  *time.Time      `json:"finished_at"`
}

// Err returns a [*RunFailedError] when the run ended failed, unsupported or cancelled, and nil otherwise:
// for a run that ended ok or partial, and for one that is still going.
func (r *Run) Err() error {
	if r.Status.IsFinal() && !r.Status.IsSuccess() {
		return &RunFailedError{Run: r}
	}
	return nil
}

// RunDiagnostics is what a run did.
type RunDiagnostics struct {
	// Requests the run made, in total.
	Requests int `json:"requests"`
	// Bytes fetched; nil when they were not measured.
	Bytes *int64 `json:"bytes"`
	// WallMS is how long the run took, in milliseconds.
	WallMS *int64 `json:"wall_ms"`
	// Facts is what the scraper learned about the target, e.g. careers_url, ats, pages_read. Which facts
	// depends on the scraper.
	Facts map[string]any `json:"facts"`
}

// RunPage is one page of runs, newest first.
type RunPage struct {
	rawBody
	Runs []Run `json:"runs"`
	// NextCursor is the cursor of the next page; empty on the last one.
	NextCursor string `json:"next_cursor"`
}

// RecordPage is one page of a run's records, in the order the run stored them.
type RecordPage struct {
	rawBody
	// Records are scraper-specific JSON objects; see the scraper's record schema. Decode one with [Decode].
	Records []json.RawMessage `json:"records"`
	// NextCursor is the cursor of the next page; empty on the last one.
	NextCursor string `json:"next_cursor"`
	// Expired is true only when the records are gone: they are kept for 7 days after the run ended.
	Expired bool `json:"expired"`
}

// Usage is the units used over a period.
type Usage struct {
	rawBody
	Period UsagePeriod `json:"period"`
	// TotalUnits is the sum of the rows.
	TotalUnits int `json:"total_units"`
	// AllowanceUnits is the account's units a calendar month (UTC): the scraper allowance, or the data
	// allowance when the request named the data product.
	AllowanceUnits int        `json:"allowance_units"`
	Rows           []UsageRow `json:"rows"`
	// Unmetered is true for an account with no unit limit (the service's own): nothing it does is refused for units.
	Unmetered bool `json:"unmetered"`
}

// Standing is where the account stands, as the API says it in the headers of every answer about runs
// (X-Units-*, X-RateLimit-*). A number the answer did not carry is nil.
type Standing struct {
	// UnitsLimit is the month's allowance of units. For an unmetered account, the largest number the API's
	// integers hold.
	UnitsLimit *int64
	// UnitsUsed is the units of the month's allowance used so far.
	UnitsUsed *int64
	// UnitsRemaining is the units of the month's allowance left. Top-up units are counted apart.
	UnitsRemaining *int64
	// UnitsTopup is the top-up units left: spent once the month's allowance is used up; they do not expire.
	UnitsTopup *int64
	// UnitsReset is when the allowance starts again: the first of next month, UTC.
	UnitsReset *time.Time
	// RunsLimit is how many runs the account may have in progress at the same time.
	RunsLimit *int64
	// RunsRemaining is how many more runs can be started now.
	RunsRemaining *int64
	// Unmetered is true for an account with no unit limit.
	Unmetered bool
}

// UnitsAvailable is the units that can still be spent: what is left of the allowance, and the top-up units.
// -1 when the answer did not say.
func (s *Standing) UnitsAvailable() int64 {
	if s == nil || s.UnitsRemaining == nil {
		return -1
	}
	n := *s.UnitsRemaining
	if s.UnitsTopup != nil {
		n += *s.UnitsTopup
	}
	return n
}

// The formats [RunService.Export] asks the API to write records in.
const (
	// FormatCSV is a header row, then a row a record: a column for each field of the scraper's record schema.
	FormatCSV = "csv"
	// FormatNDJSON is one JSON record a line.
	FormatNDJSON = "ndjson"
)

// RecordsExport is one page of a run's records as the API wrote it.
type RecordsExport struct {
	Format string
	// Body is the page, UTF-8 text.
	Body []byte
	// NextCursor is where the next page starts (the X-Next-Cursor header); empty on the last page.
	NextCursor string
	// Expired is true when the records are gone: they are kept 7 days after a run ends.
	Expired bool
}

// UsagePeriod is the first and last day (UTC, YYYY-MM-DD) a usage report covers, both included.
type UsagePeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// UsageRow is one row of a usage report.
type UsageRow struct {
	// Key is a day (YYYY-MM-DD), a scraper id, or a product ("scrapers", "data"), depending on the grouping.
	Key   string `json:"key"`
	Units int    `json:"units"`
	// Runs is what was charged: runs that ended ok or partial and, for the data API, requests that returned
	// at least one record.
	Runs int `json:"runs"`
}

// Company is a company of the dataset.
type Company struct {
	rawBody
	Domain      string        `json:"domain"`
	Name        *string       `json:"name"`
	Names       []CompanyName `json:"names"`
	CareersURL  *string       `json:"careers_url"`
	ATS         *string       `json:"ats"`
	FirstSeenAt time.Time     `json:"first_seen_at"`
	LastSeenAt  time.Time     `json:"last_seen_at"`
	LastRead    LastRead      `json:"last_read"`
	Postings    PostingCounts `json:"postings"`
}

// CompanyName is a name a company was seen under.
type CompanyName struct {
	Name        string    `json:"name"`
	FirstSeenAt time.Time `json:"first_seen_at"`
	LastSeenAt  time.Time `json:"last_seen_at"`
}

// LastRead is the last time a company's openings were read.
type LastRead struct {
	At *time.Time `json:"at"`
	// Status is "complete", "incomplete" or "failed"; nil when the company has not been read.
	Status         *string    `json:"status"`
	LastCompleteAt *time.Time `json:"last_complete_at"`
}

// PostingCounts is how many of a company's postings are open and closed.
type PostingCounts struct {
	Open   int `json:"open"`
	Closed int `json:"closed"`
}

// CompanyPage is one page of companies.
type CompanyPage struct {
	rawBody
	Companies []Company `json:"companies"`
	// NextCursor is the cursor of the next page; empty on the last one.
	NextCursor string `json:"next_cursor"`
}

// Posting is a job posting of the dataset. A value that is not known is nil.
type Posting struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	TitleClean     string          `json:"title_clean"`
	Company        CompanyRef      `json:"company"`
	Location       PostingLocation `json:"location"`
	EmploymentType *string         `json:"employment_type"`
	Salary         *Salary         `json:"salary"`
	// PostedAt and ValidThrough are dates, YYYY-MM-DD.
	PostedAt     *string       `json:"posted_at"`
	ValidThrough *string       `json:"valid_through"`
	Seniority    *string       `json:"seniority"`
	Function     *string       `json:"function"`
	URL          string        `json:"url"`
	ApplyURL     *string       `json:"apply_url"`
	Source       PostingSource `json:"source"`
	// Status is "open" or "closed".
	Status      string     `json:"status"`
	FirstSeenAt time.Time  `json:"first_seen_at"`
	LastSeenAt  time.Time  `json:"last_seen_at"`
	ClosedAt    *time.Time `json:"closed_at"`
}

// CompanyRef names the company of a posting.
type CompanyRef struct {
	Domain *string `json:"domain"`
	Name   *string `json:"name"`
}

// PostingLocation is where a job is.
type PostingLocation struct {
	Raw         *string `json:"raw"`
	City        *string `json:"city"`
	Region      *string `json:"region"`
	Country     *string `json:"country"`
	CountryCode *string `json:"country_code"`
	Remote      *bool   `json:"remote"`
}

// Salary is what a posting says it pays.
type Salary struct {
	Raw      *string  `json:"raw"`
	Min      *float64 `json:"min"`
	Max      *float64 `json:"max"`
	Currency *string  `json:"currency"`
	Period   *string  `json:"period"`
}

// PostingSource is where a posting was read.
type PostingSource struct {
	Name string `json:"name"`
	// Provenance is "feed", "page" or "board".
	Provenance string  `json:"provenance"`
	Trusted    bool    `json:"trusted"`
	ExternalID *string `json:"external_id"`
}

// PostingPage is one page of postings.
type PostingPage struct {
	rawBody
	Postings []Posting `json:"postings"`
	// NextCursor is the cursor of the next page; empty on the last one.
	NextCursor string `json:"next_cursor"`
}

// PostingDetail is one posting with its description and its history.
type PostingDetail struct {
	rawBody
	Posting
	Description    *string               `json:"description"`
	Classification PostingClassification `json:"classification"`
	History        PostingHistory        `json:"history"`
}

// PostingClassification says which rules gave a posting its seniority and function.
type PostingClassification struct {
	SeniorityRule *string `json:"seniority_rule"`
	FunctionRule  *string `json:"function_rule"`
}

// PostingHistory is what happened to a posting over time.
type PostingHistory struct {
	Events        []PostingEvent   `json:"events"`
	Changes       []PostingChange  `json:"changes"`
	Sightings     PostingSightings `json:"sightings"`
	ReopenedCount int              `json:"reopened_count"`
}

// PostingEvent is one event in a posting's life.
type PostingEvent struct {
	// Kind is "first_seen", "seen_earlier", "closed" or "reopened".
	Kind string    `json:"kind"`
	At   time.Time `json:"at"`
}

// PostingChange is one change of one field of a posting.
type PostingChange struct {
	Field string          `json:"field"`
	Old   json.RawMessage `json:"old"`
	New   json.RawMessage `json:"new"`
	At    time.Time       `json:"at"`
}

// PostingSightings is how often a posting was seen, and the latest times.
type PostingSightings struct {
	Count  int         `json:"count"`
	Recent []time.Time `json:"recent"`
}

// rawBody keeps the JSON a response was decoded from.
type rawBody struct{ raw json.RawMessage }

func (r *rawBody) setRaw(data []byte) { r.raw = data }

// RawJSON returns the JSON this value was decoded from, exactly as the API sent it — including fields this
// version of the SDK does not know. It is nil for a value that was not itself the answer of a call: one you
// built, or an item of a page (use the page's RawJSON).
func (r rawBody) RawJSON() json.RawMessage { return r.raw }
