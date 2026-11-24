# EnrichedNow Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/enrichednow/enrichednow-go.svg)](https://pkg.go.dev/github.com/enrichednow/enrichednow-go)
[![CI](https://github.com/enrichednow/enrichednow-go/actions/workflows/ci.yml/badge.svg)](https://github.com/enrichednow/enrichednow-go/actions/workflows/ci.yml)

The official Go client for the [EnrichedNow API](https://docs.enrichednow.com): run scrapers, read their records,
run formulas for finished lists of rows, and query the datasets.

- Standard library only — no dependencies.
- Go 1.22 or later.
- Every call takes a `context.Context`.

## Install

```bash
go get github.com/enrichednow/enrichednow-go@latest
```

The package is imported as `enrichednow`:

```go
import enrichednow "github.com/enrichednow/enrichednow-go"
```

## Quickstart

Create a key in the console under **API keys**, and keep it out of your code:

```bash
export ENRICHEDNOW_API_KEY="en_live_..."
```

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	enrichednow "github.com/enrichednow/enrichednow-go"
)

type Job struct {
	Title    string `json:"title"`
	ApplyURL string `json:"apply_url"`
}

func main() {
	client, err := enrichednow.NewClient() // reads ENRICHEDNOW_API_KEY
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Start a run, wait for it to end, and get its records.
	run, records, err := client.Scrape(ctx, "ats.board", map[string]any{
		"board_url": "https://boards.greenhouse.io/acme",
	})
	var failed *enrichednow.RunFailedError
	if errors.As(err, &failed) {
		log.Fatalf("the run ended %s: %s", failed.Status(), failed.Reason())
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Println(run.ID, run.Status, run.RecordCount, "records,", run.Units, "units")

	for records.Next() {
		job, err := enrichednow.Decode[Job](records.Current())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(job.Title, job.ApplyURL)
	}
	if err := records.Err(); err != nil {
		log.Fatal(err)
	}
}
```

## The client

```go
client, err := enrichednow.NewClient(
	enrichednow.WithAPIKey("en_live_..."),         // default: ENRICHEDNOW_API_KEY
	enrichednow.WithBaseURL("https://api.enrichednow.com"),
	enrichednow.WithTimeout(30*time.Second),       // per attempt
	enrichednow.WithMaxRetries(2),                 // the default
	enrichednow.WithHTTPClient(&http.Client{}),
	enrichednow.WithUserAgent("my-app/1.2"),       // put in front of enrichednow-go/<version>
)
```

Without a key, `NewClient` returns `enrichednow.ErrMissingAPIKey`. The catalogue and the health check need no key;
`enrichednow.WithoutAPIKey()` builds a client for just those.

| | |
| --- | --- |
| `client.Scrape(ctx, id, input)` | Start a run, wait for its end, return it with a records iterator. |
| `client.Scrapers.List(ctx)` | The catalogue. Each `Scraper` says what it costs: `UnitCost()`, from `UnitsPerRequest` or `UnitsPerRecord`. |
| `client.Scrapers.Get(ctx, id)` | One scraper, with the JSON Schema of its input and of one record. |
| `client.Scrapers.Run(ctx, id, input, params)` | Start a run. `params.Wait` holds the request open, up to 300 s; `params.WebhookURL` is told when it ends. |
| `client.Runs.Get(ctx, id)` / `GetWait(ctx, id, wait)` | A run as it is now, or when it has ended (up to 300 s). |
| `client.Runs.List(ctx, params)` / `ListAll` | Your runs, newest first: one page, or an iterator over all. |
| `client.Runs.Wait(ctx, id, params)` | Poll until the run has ended. |
| `client.Runs.Cancel(ctx, id)` | Ask a run to stop. |
| `client.Runs.Records(ctx, id, params)` / `RecordsPage` | A run's records: an iterator, or one page. |
| `client.Runs.Export(ctx, id, params)` / `ExportAll(ctx, id, format, w)` | A run's records as CSV or NDJSON, written by the API. |
| `client.Runs.File(ctx, id, name)` | A file a run made. |
| `client.Usage.Get(ctx, params)` | Units used over a period. |
| `client.Usage.Standing(ctx)` / `client.Standing()` | Where the account stands: units left, top-up, runs that can start. |
| `client.Webhooks.List` / `Create` / `Get` / `Update` / `Delete` / `RotateSecret` / `CreateRunSecret` | Webhook endpoints and secrets. |
| `client.Webhooks.Deliveries.List` / `ListAll` / `Resend` | The delivery log. |
| `client.Data.Companies.List` / `ListAll` / `Get` / `Postings` / `PostingsAll` | Datasets: companies. |
| `client.Data.Postings.List` / `ListAll` / `Get` | Datasets: job postings. |
| `client.RunFormula(ctx, id, params, opts)` | Run a formula, wait for its end, return it with a rows iterator. |
| `client.Formulas.List` / `Get` / `Estimate` / `Run` | The formulas: the catalogue, one with its parameters, what a run would cost, starting one. |
| `client.FormulaRuns.List` / `ListAll` / `Get` / `GetWait` / `Wait` / `Cancel` | Formula runs. |
| `client.FormulaRuns.Rows` / `RowsPage` / `Export` / `ExportAll` | A formula run's rows: an iterator, one page, or CSV and NDJSON written by the API. |
| `client.SavedFormulas.Create` / `List` / `ListAll` / `Get` / `Update` / `Archive` / `Duplicate` | Saved formulas: params and caps kept, run by hand or on a schedule. |
| `client.SavedFormulas.Run` / `Pause` / `Resume` / `ResetSeen` / `Rows` / `RowsPage` / `Export` / `ExportAll` / `Snapshots` / `Limits` | Running them, their schedule, their rows, what the plan allows. |
| `client.Suppressions.List` / `ListAll` / `Add` / `Upload` / `Delete` | The suppression list: what no formula returns. |
| `client.Attestations.List` | The terms some formulas need accepted. |
| `client.FormulaRows.Feedback` | Say what is wrong with a row. |
| `client.Health(ctx)` | Is the service up. |

## Runs

`Scrape` is the one-liner. The steps are there too, for when you want to start a run and come back for it:

```go
run, err := client.Scrapers.Run(ctx, "careers.company", map[string]any{"domain": "acme.com"}, nil)
// run.Status is "queued": store run.ID, do something else …

run, err = client.Runs.Wait(ctx, run.ID, &enrichednow.WaitParams{
	OnPoll: func(r *enrichednow.Run) { log.Println(r.ID, r.Status) },
})
if err != nil {
	return err // the context ended, or the API could not be reached
}
if err := run.Err(); err != nil {
	return err // *RunFailedError: failed, unsupported or cancelled
}
```

A run has two kinds of "did not work", and the SDK keeps them apart:

- **The request failed** — a 4xx or 5xx. You get an `*APIError`.
- **The scrape failed** — the request worked and the run ended `failed`, `unsupported` or `cancelled`.
  `Scrapers.Run`, `Runs.Get` and `Runs.Wait` return such a run without an error, as the API does; `run.Err()` and
  `client.Scrape` turn it into a `*RunFailedError` that carries the status and the reason code.

`ok` with `RecordCount == 0` is a real answer: the target was read and had nothing. `partial` has records, and a
`Reason` that says what could not be read.

How long `Wait` and `Scrape` wait is up to the context — give it a deadline. When the context ends first the run
goes on; the run that comes back with the error has the id to follow or cancel it.

## Records and paging

Records are scraper-specific JSON objects, so they come as `json.RawMessage`; the scraper's `RecordSchema` says what
is in them. Decode one with `Decode`, or all of them with `DecodeAll`:

```go
jobs, err := enrichednow.DecodeAll[Job](client.Runs.Records(ctx, run.ID, nil))
```

Every list has an iterator that fetches page after page:

```go
it := client.Runs.Records(ctx, run.ID, &enrichednow.RecordListParams{Limit: 1000})
for it.Next() {
	record := it.Current()
	// …
}
if err := it.Err(); err != nil { /* … */ }
```

With Go 1.23 or later, range over it:

```go
for record, err := range client.Runs.Records(ctx, run.ID, nil).All() {
	if err != nil { /* … */ }
	// …
}
```

And a page at a time, when you want the cursor yourself:

```go
page, err := client.Runs.RecordsPage(ctx, run.ID, &enrichednow.RecordListParams{Limit: 100})
// page.Records, page.NextCursor ("" on the last page), page.Expired
```

Records are kept for 7 days after a run ended. After that a records iterator stops with
`enrichednow.ErrRecordsExpired`; the run itself stays.

### CSV and NDJSON, written by the API

```go
f, _ := os.Create("jobs.csv")
defer f.Close()
err := client.Runs.ExportAll(ctx, run.ID, enrichednow.FormatCSV, f) // every page; the header row once

page, err := client.Runs.Export(ctx, run.ID, enrichednow.ExportParams{Format: enrichednow.FormatNDJSON, Limit: 500})
// page.Body, page.NextCursor, page.Expired
```

The CSV has a column for each field of the scraper's record schema, in its order; a nested object becomes
`parent.child` columns, a list one cell of JSON, and a cell a spreadsheet would run as a formula starts with an
apostrophe. A file a run made (a screenshot) is `client.Runs.File(ctx, run.ID, "screenshot.png")`.

### Waiting longer in one call

`RunParams.Wait` and `Runs.GetWait(ctx, id, wait)` hold the request open until the run ends, up to `MaxWait`
(300 seconds). `RunListParams` also has `APIKeyID`, `From` and `To`.

## Where the account stands

Every answer about runs carries the account's units and runs in progress in its headers. `client.Standing()` is the
last one seen; `client.Usage.Standing(ctx)` asks for it:

```go
s, err := client.Usage.Standing(ctx)
fmt.Println(*s.UnitsRemaining, *s.UnitsTopup, s.UnitsReset, *s.RunsRemaining, s.Unmetered, s.UnitsAvailable())
```

`Usage.Unmetered` is true for an account with no unit limit.

## Webhooks

Endpoints are told when a run ends, signed with their secret:

```go
made, err := client.Webhooks.Create(ctx, enrichednow.WebhookEndpointParams{
	URL:    "https://example.com/hooks/enrichednow",
	Events: []string{enrichednow.EventRunSucceeded, enrichednow.EventRunFailed},
})
// made.Secret is shown this once: keep it.

list, err := client.Webhooks.List(ctx)
off := false
_, err = client.Webhooks.Update(ctx, made.ID, enrichednow.WebhookEndpointChange{Enabled: &off})
_, err = client.Webhooks.RotateSecret(ctx, made.ID)
err = client.Webhooks.Delete(ctx, made.ID)

secret, err := client.Webhooks.CreateRunSecret(ctx) // for RunParams.WebhookURL, one run's own URL
for d, err := range client.Webhooks.Deliveries.ListAll(ctx, &enrichednow.WebhookDeliveryListParams{Status: "failed"}).All() { /* … */ }
_, err = client.Webhooks.Deliveries.Resend(ctx, deliveryID)
```

And a request that arrives is checked before it is acted on: give the body exactly as it arrived.

```go
http.HandleFunc("/hooks/enrichednow", func(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	event, err := enrichednow.ParseWebhook(os.Getenv("ENRICHEDNOW_WEBHOOK_SECRET"), r.Header, body)
	if err != nil {
		http.Error(w, "bad signature", http.StatusBadRequest)
		return
	}
	fmt.Println(event.ID, event.Type, event.Data.Run.ID)
	w.WriteHeader(http.StatusNoContent)
})
```

`VerifyWebhook` and `VerifyWebhookSignature` only check; both compare in constant time and refuse a timestamp more
than 5 minutes from now. `event.ID` is the same on every attempt: use it to ignore a repeat.

## Datasets

```go
page, err := client.Data.Postings.List(ctx, &enrichednow.PostingListParams{
	Title:   "data engineer",
	Country: "US",
	Status:  "open",
	Limit:   20,
})
for _, posting := range page.Postings {
	fmt.Println(posting.Title, posting.URL)
}

company, err := client.Data.Companies.Get(ctx, "figma.com")
detail, err := client.Data.Postings.Get(ctx, page.Postings[0].ID) // with description and history
```

The account must have been given Datasets; without it these calls answer `ErrPermission`. A data request costs one
data unit for each record it returns, so `ListAll` and `PostingsAll` are charged for every page they read.

In these types a value that is not known is `nil` — never an empty string and never `0`.

## Formulas

A formula is a ready-made way from a few parameters to a finished list of rows: companies hiring for your roles
with the person to contact, an account list researched, a directory turned into leads. Each formula has a
parameter type in this package, made from its parameters' schema (`JobOrderHunterParams`, `AccountResearchParams`
… `enrichednow.FormulaIDs` lists them); a map or a struct of yours works too, for any formula.

```go
run, rows, err := client.RunFormula(ctx, "job-order-hunter", enrichednow.JobOrderHunterParams{
	JobTitles:        []string{"Java Developer"},
	Locations:        []string{"London"},
	PostedWithinDays: enrichednow.Ptr(7), // an optional field is a pointer: nil leaves the formula's default
}, &enrichednow.RunFormulaParams{UnitCap: 2000})
var failed *enrichednow.FormulaRunFailedError
if errors.As(err, &failed) {
	log.Fatalf("the run ended %s: %s", failed.Status(), failed.Reason())
} else if err != nil {
	log.Fatal(err)
}
fmt.Println(run.Status, run.RowsTotal, "rows,", run.UnitsSpent, "units") // partial is not an error: run.Reasons say why

type Lead struct {
	CompanyName  string `json:"company_name"`
	ContactEmail string `json:"contact_email"`
}
for rows.Next() {
	row := rows.Current() // ID, Data (the columns), Score, Evidence (where each column came from), Status
	lead, _ := enrichednow.Decode[Lead](row.Data)
	fmt.Println(lead.CompanyName, lead.ContactEmail, *row.Score)
}
```

A parameter type knows its formula: giving `JobOrderHunterParams` to another formula is refused before anything is
sent. `RunFormula` waits up to 10 minutes (`RunFormulaParams.Timeout`); when the time is up the run goes on, and it
comes back with an error that matches `context.DeadlineExceeded`.

The steps, for when you want them:

```go
details, err := client.Formulas.Get(ctx, "job-order-hunter")   // what it does, details.Parameters (JSON Schema), Columns, Presets
list, err := client.Formulas.List(ctx, &enrichednow.FormulaListParams{Category: "recruiting"}) // the catalogue needs no key

estimate, err := client.Formulas.Estimate(ctx, "job-order-hunter", params, nil) // spends nothing
if !estimate.Valid {
	for _, problem := range estimate.Errors { // params that cannot run are an answer, not an error
		fmt.Println(problem.Path, problem.Message)
	}
}
fmt.Println(estimate.Estimate.ExpectedUnits, estimate.Estimate.MaxUnits)

run, err := client.Formulas.Run(ctx, "job-order-hunter", params, &enrichednow.FormulaRunParams{
	Preview: true,           // at most 10 rows and 50 units
	UnitCap: 2000,           // the most it may spend: it stops at the cap and ends partial
	Wait:    time.Minute,    // hold the request until it ends, up to 300 s
})
run, err = client.FormulaRuns.Wait(ctx, run.ID, &enrichednow.FormulaWaitParams{Timeout: 5 * time.Minute})
_, err = client.FormulaRuns.Cancel(ctx, run.ID)

it := client.FormulaRuns.Rows(ctx, run.ID, &enrichednow.FormulaRowListParams{
	Fields:   []string{"company_name", "contact_email"},
	MinScore: enrichednow.Ptr(30.0),
	View:     enrichednow.ViewCompany, // one row per company, for a formula with details.CompanyView
})

f, _ := os.Create("leads.csv")
err = client.FormulaRuns.ExportAll(ctx, run.ID, enrichednow.FormulaExportParams{Format: enrichednow.FormatCSV}, f) // the header once
```

Starting a run carries an `Idempotency-Key` and is retried as a scraper run is (`FormulaRunParams.IdempotencyKey` to
make a call of yours repeatable). A formula whose `NeedsAttestation` is true returns people: an owner or an admin
must accept its terms in the console first (`client.Attestations.List` shows them; a key cannot accept them).

### Saved formulas

A saved formula keeps its params and caps, runs by hand or on a schedule, and delivers only rows it did not
deliver before. Poll its new rows with `Since`:

```go
saved, err := client.SavedFormulas.Create(ctx, "job-order-hunter", params, &enrichednow.SavedFormulaCreateParams{
	Schedule:       &enrichednow.ScheduleParams{Cadence: enrichednow.CadenceWeekdays, LocalTime: "07:30", Timezone: "Europe/London"},
	MonthlyUnitCap: 20000,
	Notify:         enrichednow.NotifyDigest,
})

for row, err := range client.SavedFormulas.Rows(ctx, saved.ID, &enrichednow.SavedFormulaRowListParams{Since: lastLook}).All() {
	// …
}

run, err := client.SavedFormulas.Run(ctx, saved.ID, &enrichednow.SavedFormulaRunParams{Wait: time.Minute})
_, err = client.SavedFormulas.Pause(ctx, saved.ID) // and Resume
_, err = client.SavedFormulas.Update(ctx, saved.ID, enrichednow.SavedFormulaUpdateParams{
	RemoveSchedule: true, // the API's null: AutoName and RemoveMonthlyUnitCap take a setting away the same way
})
copied, err := client.SavedFormulas.Duplicate(ctx, saved.ID, &enrichednow.SavedFormulaDuplicateParams{CopySeen: true})
_, err = client.SavedFormulas.Archive(ctx, saved.ID)
limits, err := client.SavedFormulas.Limits(ctx) // schedules, the shortest interval, how long rows are kept
```

### The suppression list, and feedback on a row

Domains, addresses and company names on the suppression list never come back in a formula's rows:

```go
added, err := client.Suppressions.Add(ctx, []enrichednow.SuppressionEntry{
	{Kind: enrichednow.SuppressDomain, Value: "globex.example"},
	{Kind: enrichednow.SuppressEmail, Value: "ana@initech.example"},
}, nil)

f, _ := os.Open("customers.csv") // a header row with kind and value, or one value a line
upload, err := client.Suppressions.Upload(ctx, f, nil)
fmt.Println(upload.Added, upload.Already, upload.Invalid)

err = client.Suppressions.Delete(ctx, enrichednow.SuppressDomain, "globex.example")

_, err = client.FormulaRows.Feedback(ctx, row.ID, enrichednow.VerdictWrongCompany,
	&enrichednow.RowFeedbackParams{CorrectDomain: "corvid.example"}) // this account's runs match it right from now on
```

## Errors

```go
_, err := client.Scrapers.Run(ctx, "ats.board", input, nil)

var apiErr *enrichednow.APIError
switch {
case errors.Is(err, enrichednow.ErrValidation) && errors.As(err, &apiErr):
	for _, f := range apiErr.Fields() { // one per problem
		fmt.Println(f.Path, f.Message) // input.board_url  Field required
	}
case errors.Is(err, enrichednow.ErrLimitReached):  // your own spend limit (the account's or the key's): raise it in the console
case errors.Is(err, enrichednow.ErrQuotaExceeded): // this month's units are used up
case errors.Is(err, enrichednow.ErrTooManyRuns):   // wait for one of your runs to end
case errors.Is(err, enrichednow.ErrAuthentication):
}
```

An `*APIError` has `StatusCode`, the API's `Code` (match on this), `Message` (for a person), `Details`,
`RetryAfter`, and `RequestID` when the response carried one.

| `errors.Is(err, …)` | HTTP | API code |
| --- | --- | --- |
| `ErrAuthentication` | 401 | `unauthorized` |
| `ErrQuotaExceeded` | 402 | `quota_exceeded`, and `limit_reached` |
| `ErrLimitReached` | 402 | `limit_reached`: a spend limit of the account's own, of the key, or a saved formula's monthly cap (`Details["scope"]`: `account`, `key`, `saved_formula`, `formula_run`) |
| `ErrPermission` | 403 | `forbidden`, `plan_limit_reached` |
| `ErrPlanLimitReached` | 403 | `plan_limit_reached`: the plan does not allow it (`Details["limit"]`, `["allowed"]`) |
| `ErrNotFound` | 404 | `not_found`, `formula_not_found` |
| `ErrConflict` | 409 | `idempotency_conflict`, and the four below |
| `ErrAttestationRequired` | 409 | `attestation_required`: an owner or an admin must accept the formula's terms in the console |
| `ErrFormulaUnavailable` | 409 | `formula_unavailable`: e.g. a deprecated formula |
| `ErrRequirementMissing` | 409 | `requirement_missing`: `apiErr.Missing()` says what the account lacks |
| `ErrSavedFormulaArchived` | 409 | `saved_formula_archived`: duplicate it to use it again |
| `ErrValidation` | 422 | `invalid_input`, `unit_cap_too_low` (`apiErr.MinUnitCap()`) |
| `ErrRateLimit` | 429 | `too_many_runs`, `rate_limited` |
| `ErrTooManyRuns` | 429 | `too_many_runs` |
| `ErrServer` | 5xx | `internal`, `data_unavailable` |
| `ErrDataUnavailable` | 503 | `data_unavailable` |
| `ErrRunFailed` | — | a `*RunFailedError` or a `*FormulaRunFailedError`: the run ended `failed`, `unsupported` or `cancelled` |

A request that got no answer — a dropped connection, a timeout — is a `*ConnectionError`.

## Retries

Requests that are safe to repeat — every read, cancelling a run, and starting one — are sent again up to twice
(`WithMaxRetries`) after a connection error, a 5xx, or a 429 that is about the rate of requests. The pause is
what `Retry-After` asks for, or else 0.5 s doubling up to 8 s, with jitter.

**Starting a run is safe to repeat too.** `Scrapers.Run` and `Scrape` send an `Idempotency-Key` with the request, a
new one for each call, so a retry gets the run the first attempt started and never a second one — and so do
`Formulas.Run`, `RunFormula` and `SavedFormulas.Run`. Creating a webhook endpoint or a saved formula, and
duplicating a saved formula, are not sent again after a failure that may have reached the API. To make a call of
your own repeatable (a job that may be restarted, say), set `RunParams.IdempotencyKey`: for 24 hours the same key
with the same request gives the same run, and with a different request fails with `ErrConflict`.

## The raw JSON

The values a call returns keep the JSON they were decoded from, exactly as the API sent it, including fields this
version of the SDK does not know yet:

```go
run, _ := client.Runs.Get(ctx, id)
fmt.Println(string(run.RawJSON()))
```

## Development

```bash
go test ./...
go vet ./...
gofmt -l .
```

The tests run against `net/http/httptest` servers; nothing calls the real API.

## Licence

MIT. See [LICENSE](LICENSE). To report a security problem, see [SECURITY.md](SECURITY.md).
