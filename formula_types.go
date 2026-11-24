package enrichednow

import (
	"encoding/json"
	"time"
)

// The statuses of a formula in the catalogue.
const (
	FormulaDraft      = "draft"
	FormulaBeta       = "beta"
	FormulaLive       = "live"
	FormulaDeprecated = "deprecated"
)

// The modes of a formula run. A preview is at most 10 rows and 50 units, and nothing it returns is remembered as
// delivered.
const (
	ModeFull    = "full"
	ModePreview = "preview"
)

// The views of a formula run's rows. [ViewCompany] is one row per company, for a formula whose
// [FormulaDetails.CompanyView] is set.
const (
	ViewRows    = "rows"
	ViewCompany = "company"
)

// The cadences of a saved formula's schedule. A formula's [Formula.Schedules] says which it allows, and the plan
// sets the shortest ([PlanLimits.MinIntervalMinutes]).
const (
	CadenceHourly   = "hourly"
	CadenceDaily    = "daily"
	CadenceWeekdays = "weekdays"
	CadenceWeekly   = "weekly"
	CadenceMonthly  = "monthly"
)

// When a saved formula's scheduled runs send mail about new rows.
const (
	NotifyNone    = "none"
	NotifyEachRun = "each_run"
	NotifyDigest  = "digest"
)

// The statuses of a saved formula.
const (
	SavedFormulaActive   = "active"
	SavedFormulaPaused   = "paused"
	SavedFormulaArchived = "archived"
)

// The kinds of entry on the suppression list.
const (
	SuppressDomain      = "domain"
	SuppressEmail       = "email"
	SuppressCompanyName = "company_name"
)

// The verdicts of [FormulaRowService.Feedback].
const (
	VerdictGood         = "good"
	VerdictWrongCompany = "wrong_company"
	VerdictWrongContact = "wrong_contact"
	VerdictNotRelevant  = "not_relevant"
)

// Formula is one entry of the catalogue of formulas: a ready-made way from a few parameters to a finished list of
// rows.
type Formula struct {
	ID string `json:"id"`
	// Version is the formula's version, such as "1.1.0".
	Version  string `json:"version"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	UseWhen  string `json:"use_when"`
	Category string `json:"category"`
	// Personas are who it is made for, such as "recruiter".
	Personas []string `json:"personas"`
	// Status is draft, beta, live or deprecated. A deprecated formula can be read, not run.
	Status string `json:"status"`
	// Available says whether it can run now, for this caller.
	Available bool `json:"available"`
	// NeedsProvider: its contact columns stay empty until a provider is connected; the rest of the row is filled.
	NeedsProvider bool `json:"needs_provider"`
	// RecommendedCadence is how often a saved formula of it is best run, such as "weekdays", or "once".
	RecommendedCadence string `json:"recommended_cadence"`
	// Schedules are the cadences a saved formula of it may have.
	Schedules []string `json:"schedules"`
	// NeedsAttestation: it returns people, and an owner or an admin of the account must accept the terms for it,
	// in the console, before it runs. See [AttestationService.List].
	NeedsAttestation bool `json:"needs_attestation"`
	// Health is how it has been running, as last measured; nil before any run.
	Health *FormulaHealth `json:"health"`
	// Basics are the questions it asks before it runs, the required ones first: everything else has a default.
	Basics []FormulaInput `json:"basics"`
}

// FormulaInput is one question a formula asks before it runs: a parameter of its Basics.
type FormulaInput struct {
	// Name is the key in the params.
	Name  string `json:"name"`
	Title string `json:"title"`
	// Type is "string", "integer", "boolean", "object", or a list: "list of strings".
	Type string `json:"type"`
	// Required: a run is refused without it. Otherwise it has a default.
	Required bool `json:"required"`
	// Help is one sentence on what to give; empty when there is none.
	Help string `json:"help"`
}

// FormulaHealth is how a formula has been running.
type FormulaHealth struct {
	// Status is "ok", or "degraded" when a source it reads is failing or returning much less than usual.
	Status string `json:"status"`
	// DegradedSources are the sources behind a degraded status.
	DegradedSources []string `json:"degraded_sources"`
	// SuccessRate7d is the share of its runs in the last 7 days that ended ok or partial; nil when there were too few.
	SuccessRate7d *float64 `json:"success_rate_7d"`
	// MedianSeconds7d is how long its runs took, the median of the last 7 days; nil when there were too few.
	MedianSeconds7d *float64 `json:"median_seconds_7d"`
}

// FormulaDetails is one formula with everything needed to run it: its parameters, presets, steps, columns, what
// it needs, its limits and sample rows.
type FormulaDetails struct {
	rawBody
	Formula
	Description string `json:"description"`
	// Parameters is the JSON Schema of a run's params. Its x-group, x-widget and x-help keywords describe the
	// setup form. The parameter types of this package (such as [JobOrderHunterParams]) are made from it.
	Parameters json.RawMessage `json:"parameters"`
	// Presets are sets of params for common cases.
	Presets []FormulaPreset `json:"presets"`
	// Steps are what a run does, in order.
	Steps []FormulaStep `json:"steps"`
	// Columns are the columns of its rows, in order.
	Columns []FormulaColumn `json:"columns"`
	// CompanyView, when set, says its rows can also be read one row per company ([ViewCompany]).
	CompanyView *FormulaCompanyView `json:"company_view"`
	// Requirements is what it needs, and what of that this account has.
	Requirements FormulaRequirements `json:"requirements"`
	// Limits are its bounds: "default_rows", "max_rows", "default_unit_cap", "max_unit_cap", "min_unit_cap" and others.
	Limits map[string]int `json:"limits"`
	// SampleRows show what a row looks like, free to read. Decode one with [Decode].
	SampleRows []json.RawMessage `json:"sample_rows"`
	// SampleIsFixture is true while the sample is made-up data of the right shape, not rows from a real run.
	SampleIsFixture bool     `json:"sample_is_fixture"`
	Changelog       []string `json:"changelog"`
}

// FormulaPreset is a named set of params.
type FormulaPreset struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	// Params are the preset's parameters: give them to [FormulaService.Run] as they are, or change some first.
	Params json.RawMessage `json:"params"`
}

// FormulaStep is one step of what a formula does.
type FormulaStep struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	// Label is what the step is called on the formula's page; empty when it has none.
	Label string `json:"label"`
}

// FormulaColumn is one column of a formula's rows.
type FormulaColumn struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Type  string `json:"type"`
	// DefaultVisible says whether the console shows it without being asked.
	DefaultVisible bool `json:"default_visible"`
}

// FormulaCompanyView is how a formula's rows are read one row per company.
type FormulaCompanyView struct {
	// Key is the column that names the company.
	Key string `json:"key"`
	// Lists are the columns gathered into a list for each company.
	Lists []string `json:"lists"`
	Label string   `json:"label"`
}

// FormulaRequirements is what a formula needs, and what of that the account has.
type FormulaRequirements struct {
	// Scrapers are what it reads from.
	Scrapers []string `json:"scrapers"`
	// Capabilities are the outside services it can use: people_search, contact_find, email_verify.
	Capabilities []string `json:"capabilities"`
	// CapabilitiesConnected are those this account has connected; empty for a caller without a key.
	CapabilitiesConnected []string `json:"capabilities_connected"`
	// Products are the products the account needs, such as "scrapers" or "data".
	Products []string `json:"products"`
	// Missing is what the account lacks for the formula to run at all, such as "product:data".
	Missing []string `json:"missing"`
	// DisabledSources are sources switched off for now: runs go on without them, and say so.
	DisabledSources []string `json:"disabled_sources"`
	// NotOffered are sources not offered here: a run goes on without them and names them in its reasons.
	NotOffered []string `json:"not_offered"`
}

// Estimate is about what a formula run costs, and how many rows it gives. Only MaxUnits is a promise: a run
// cannot pass its unit cap.
type Estimate struct {
	// MinUnits is the least the run can cost: with a unit cap below it the run is refused.
	MinUnits int `json:"min_units"`
	// ExpectedUnits is about what it will cost.
	ExpectedUnits int `json:"expected_units"`
	// MaxUnits is the most it can cost: never above its unit cap.
	MaxUnits     int `json:"max_units"`
	ExpectedRows int `json:"expected_rows"`
	// Basis says what drives the numbers, in one sentence.
	Basis string `json:"basis"`
}

// FormulaEstimate is the answer of [FormulaService.Estimate]. Params that cannot run are not an error: Valid is
// false, Errors says what to change, and Estimate is nil.
type FormulaEstimate struct {
	rawBody
	Valid bool `json:"valid"`
	// Errors are the problems of the params, one per field, such as {"params.job_titles", "This is required."}.
	Errors   []FieldError `json:"errors"`
	Estimate *Estimate    `json:"estimate"`
	// Missing is what the account lacks for these params.
	Missing []string `json:"missing"`
	// NotOffered are sources these params ask for that are not offered here: the run goes on without them.
	NotOffered []string `json:"not_offered"`
}

// FormulaRun is one run of a formula with one set of params. Its Status is a [RunStatus], as a scraper run's is.
type FormulaRun struct {
	rawBody
	ID        string `json:"id"`
	FormulaID string `json:"formula_id"`
	// Version is the version of the formula the run used.
	Version string `json:"version"`
	// SavedFormulaID is the saved formula the run is of; empty for a run started with params.
	SavedFormulaID string    `json:"saved_formula_id"`
	Status         RunStatus `json:"status"`
	// Reason says why it is not plainly ok, as a stable code: "unit_cap_reached", "contacts_not_connected",
	// "source_failed:<scraper>" and others. Empty when the API sent null.
	Reason string `json:"reason"`
	// Reasons are every reason, when there is more than one, with the step or source each concerns.
	Reasons []FormulaRunReason `json:"reasons"`
	// Mode is [ModeFull] or [ModePreview].
	Mode string `json:"mode"`
	// Trigger is "api", "console", "mcp", or "schedule" for a saved formula's scheduled run.
	Trigger string `json:"trigger"`
	// Params are the params the run used, the formula's defaults filled in. Decode them with [Decode], into the
	// formula's parameter type or a map.
	Params json.RawMessage `json:"params"`
	// Estimate is what it was estimated to cost when it started; nil when there was none.
	Estimate *Estimate `json:"estimate"`
	// UnitCap is the most the run may spend. It cannot be overspent: the run stops at the cap and ends partial.
	UnitCap int `json:"unit_cap"`
	// RowCap is the most rows it may produce.
	RowCap int `json:"row_cap"`
	// UnitsSpent is the units charged so far: the sum of its steps' units.
	UnitsSpent int `json:"units_spent"`
	RowsTotal  int `json:"rows_total"`
	// IncludeSeen: a saved formula's run that was asked to deliver rows already delivered too.
	IncludeSeen bool `json:"include_seen"`
	// ScheduledFor is the slot of the schedule a scheduled run was for; nil otherwise.
	ScheduledFor *time.Time `json:"scheduled_for"`
	// Steps say what each step did. They are on one run ([FormulaRunService.Get], starting a run), not in a list.
	Steps      []FormulaRunStep `json:"steps"`
	CreatedAt  time.Time        `json:"created_at"`
	StartedAt  *time.Time       `json:"started_at"`
	FinishedAt *time.Time       `json:"finished_at"`
}

// Err returns a [*FormulaRunFailedError] when the run ended failed, unsupported or cancelled, and nil otherwise:
// for a run that ended ok or partial, and for one that is still going.
func (r *FormulaRun) Err() error {
	if r.Status.IsFinal() && !r.Status.IsSuccess() {
		return &FormulaRunFailedError{Run: r}
	}
	return nil
}

// FormulaRunReason is one reason a formula run gives: a code, and what it concerns.
type FormulaRunReason struct {
	// Code is stable, such as "contacts_not_connected", "source_not_offered" or "source_failed:<scraper>".
	Code string
	// Fields is the whole reason as the API sent it, code included: "step", "scraper" and others, by code.
	Fields map[string]any
}

// UnmarshalJSON reads a reason: its code, and every other field into Fields.
func (r *FormulaRunReason) UnmarshalJSON(data []byte) error {
	return decodeCoded(data, &r.Code, &r.Fields)
}

// MarshalJSON writes the reason as the API sent it.
func (r FormulaRunReason) MarshalJSON() ([]byte, error) {
	return encodeCoded(r.Code, r.Fields)
}

// FormulaRunStep is what one step of a formula run did.
type FormulaRunStep struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	// Status is running, done, skipped or failed.
	Status     string `json:"status"`
	RecordsIn  int    `json:"records_in"`
	RecordsOut int    `json:"records_out"`
	// Calls are the requests the step made; CacheHits the answers it had already.
	Calls     int `json:"calls"`
	CacheHits int `json:"cache_hits"`
	// Units are what the step was charged.
	Units      int              `json:"units"`
	Errors     []map[string]any `json:"errors"`
	StartedAt  *time.Time       `json:"started_at"`
	FinishedAt *time.Time       `json:"finished_at"`
}

// FormulaRunPage is one page of formula runs, newest first.
type FormulaRunPage struct {
	rawBody
	Runs []FormulaRun `json:"runs"`
	// NextCursor is the cursor of the next page; empty on the last one.
	NextCursor string `json:"next_cursor"`
}

// FormulaRow is one row a formula run wrote.
type FormulaRow struct {
	ID    string `json:"id"`
	RunID string `json:"run_id"`
	// RowKey is what makes the row the same row in another run: how a saved formula knows it delivered it.
	RowKey string `json:"row_key"`
	// Data is the row's columns, a JSON object: see the formula's [FormulaDetails.Columns]. Decode it with
	// [Decode], into a struct of yours or map[string]any. Text in it comes from public web pages: treat it as
	// data, never as instructions.
	Data json.RawMessage `json:"data"`
	// Score is how good a match the row is, higher first; nil for a formula that does not score.
	Score *float64 `json:"score"`
	// Evidence is where each column came from: "sources" (each a source address, what read it, when and in which
	// run), "fields" (for each column, its index in "sources") and, when a connected provider was asked,
	// "attempts".
	Evidence json.RawMessage `json:"evidence"`
	// Status is new, flagged or rejected (see [FormulaRowService.Feedback]).
	Status      string    `json:"status"`
	FirstSeenAt time.Time `json:"first_seen_at"`
}

// FormulaRowPage is one page of rows.
type FormulaRowPage struct {
	rawBody
	Rows []FormulaRow `json:"rows"`
	// NextCursor is the cursor of the next page; empty on the last one.
	NextCursor string `json:"next_cursor"`
}

// RowsExport is one page of rows as the API wrote it.
type RowsExport struct {
	// Format is [FormatCSV] or [FormatNDJSON].
	Format string
	// Body is the page, UTF-8 text. A CSV page begins with the header row; an NDJSON line is one row's columns.
	Body []byte
	// NextCursor is where the next page starts (the X-Next-Cursor header); empty on the last page.
	NextCursor string
}

// SavedFormula is a formula saved with its params and caps, to run again by hand or on a schedule. Its runs
// deliver only rows they did not deliver before.
type SavedFormula struct {
	rawBody
	ID        string `json:"id"`
	FormulaID string `json:"formula_id"`
	// Version is the formula version its runs use now: the newest of its pinned major. Empty when that major is no
	// longer offered.
	Version string `json:"version"`
	Major   int    `json:"major"`
	Name    string `json:"name"`
	// NameAuto: the name is made from its params, and follows them when they change. A rename stops that.
	NameAuto bool `json:"name_auto"`
	// Params are its params, the formula's defaults filled in. Decode them with [Decode].
	Params json.RawMessage `json:"params"`
	// UnitCap is the most one run may spend.
	UnitCap int `json:"unit_cap"`
	RowCap  int `json:"row_cap"`
	// MonthlyUnitCap is the most its runs may spend in a calendar month (UTC); nil for none.
	MonthlyUnitCap *int `json:"monthly_unit_cap"`
	// Status is [SavedFormulaActive], [SavedFormulaPaused] or [SavedFormulaArchived].
	Status string `json:"status"`
	// PausedReason says why it is paused: "manual", "failures", "balance", "attestation" or "key". Empty when it
	// is not.
	PausedReason string `json:"paused_reason"`
	// RunsWithKey is the API key its schedule runs with, whose limits apply to its scheduled runs; empty when a
	// person set the schedule.
	RunsWithKey string `json:"runs_with_key"`
	// Notify is when scheduled runs send mail about new rows: [NotifyEachRun], [NotifyDigest] or [NotifyNone].
	Notify string `json:"notify"`
	// Schedule is nil for a saved formula that runs only by hand.
	Schedule            *Schedule `json:"schedule"`
	ConsecutiveFailures int       `json:"consecutive_failures"`
	// EmptyRuns are the scheduled runs in a row that found no new rows.
	EmptyRuns int                `json:"empty_runs"`
	Notes     []SavedFormulaNote `json:"notes"`
	// Estimate is about what one run costs, as estimated when it was last changed.
	Estimate *Estimate            `json:"estimate"`
	LastRun  *SavedFormulaLastRun `json:"last_run"`
	// Units30d are the units its runs were charged in the last 30 days.
	Units30d   int        `json:"units_30d"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at"`
}

// Schedule is when a saved formula runs.
type Schedule struct {
	// Cadence is hourly, daily, weekdays, weekly or monthly.
	Cadence string `json:"cadence"`
	// LocalTime is HH:MM, 24-hour, in Timezone. For hourly, only the minute counts.
	LocalTime string `json:"local_time"`
	// Timezone is an IANA time zone, such as Europe/London.
	Timezone string `json:"timezone"`
	// DayOfWeek is the day of a weekly schedule, 0 for Monday; nil for other cadences.
	DayOfWeek *int `json:"day_of_week"`
	// DayOfMonth is the day of a monthly schedule; nil for other cadences.
	DayOfMonth *int `json:"day_of_month"`
	// NextFireAt is the next run, UTC; nil while the saved formula is paused or archived.
	NextFireAt *time.Time `json:"next_fire_at"`
}

// SavedFormulaNote is something a saved formula wants you to know.
type SavedFormulaNote struct {
	// Code is "source_disabled" (a source it reads is switched off for now, and its runs go on without it) or
	// "widen" (its last 14 scheduled runs found nothing new).
	Code    string
	Message string
	At      time.Time
	// Fields is the whole note as the API sent it, such as "scraper" for source_disabled.
	Fields map[string]any
}

// UnmarshalJSON reads a note: its code, message and time, and every field into Fields.
func (n *SavedFormulaNote) UnmarshalJSON(data []byte) error {
	if err := decodeCoded(data, &n.Code, &n.Fields); err != nil {
		return err
	}
	n.Message, _ = n.Fields["message"].(string)
	if at, ok := n.Fields["at"].(string); ok {
		n.At, _ = time.Parse(time.RFC3339Nano, at)
	}
	return nil
}

// MarshalJSON writes the note as the API sent it.
func (n SavedFormulaNote) MarshalJSON() ([]byte, error) {
	fields := map[string]any{"message": n.Message, "at": n.At}
	for k, v := range n.Fields {
		fields[k] = v
	}
	return encodeCoded(n.Code, fields)
}

// SavedFormulaLastRun is the latest run of a saved formula, in short.
type SavedFormulaLastRun struct {
	ID     string    `json:"id"`
	Status RunStatus `json:"status"`
	// Reason is empty when the API sent null.
	Reason string `json:"reason"`
	// RowsNew are the rows it delivered: rows the saved formula had not delivered before.
	RowsNew    int        `json:"rows_new"`
	CreatedAt  time.Time  `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

// SavedFormulaPage is one page of saved formulas.
type SavedFormulaPage struct {
	rawBody
	SavedFormulas []SavedFormula `json:"saved_formulas"`
	// NextCursor is the cursor of the next page; empty on the last one.
	NextCursor string `json:"next_cursor"`
}

// SeenReset is the answer of [SavedFormulaService.ResetSeen].
type SeenReset struct {
	rawBody
	// Cleared is how many delivered rows were forgotten.
	Cleared      int          `json:"cleared"`
	SavedFormula SavedFormula `json:"saved_formula"`
}

// Snapshot is what a saved formula observed of one entity (a price, a stock level, a title), each time it changed.
type Snapshot struct {
	Entity    string         `json:"entity"`
	EntityKey string         `json:"entity_key"`
	Fields    map[string]any `json:"fields"`
	// RunID is the run that observed it; empty when the API sent null.
	RunID      string    `json:"run_id"`
	ObservedAt time.Time `json:"observed_at"`
}

// PlanLimits is what the account's plan allows of saved formulas.
type PlanLimits struct {
	rawBody
	// Schedules is how many saved formulas may have an active schedule at once; SchedulesActive how many do.
	Schedules       int `json:"schedules"`
	SchedulesActive int `json:"schedules_active"`
	// MinIntervalMinutes is the shortest interval between scheduled runs: 60 is hourly, 1440 daily.
	MinIntervalMinutes int `json:"min_interval_minutes"`
	// RowRetentionDays is how long rows are kept.
	RowRetentionDays int `json:"row_retention_days"`
	// Destinations is how many destinations one saved formula may have.
	Destinations int `json:"destinations"`
}

// Suppression is one entry of the suppression list: a domain, an address or a company name that no formula run
// returns or sends.
type Suppression struct {
	// Kind is [SuppressDomain], [SuppressEmail] or [SuppressCompanyName].
	Kind string `json:"kind"`
	// Value is the entry in the form runs match it: a bare domain, a lower-case address, a company name without
	// its legal form ("Acme, Inc." is "acme").
	Value string `json:"value"`
	// Source is "manual" or "upload".
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

// SuppressionEntry is one entry to add to the suppression list.
type SuppressionEntry struct {
	// Kind is [SuppressDomain], [SuppressEmail] or [SuppressCompanyName].
	Kind string `json:"kind"`
	// Value is a domain (any address of the site will do), an email address, or a company name.
	Value string `json:"value"`
}

// SuppressionPage is one page of the suppression list.
type SuppressionPage struct {
	rawBody
	Suppressions []Suppression `json:"suppressions"`
	// NextCursor is the cursor of the next page; empty on the last one.
	NextCursor string `json:"next_cursor"`
}

// SuppressionUpload is the answer of [SuppressionService.Upload].
type SuppressionUpload struct {
	rawBody
	// Added are the entries new to the list.
	Added int `json:"added"`
	// Already are the entries the list already held, left as they are.
	Already int `json:"already"`
	// Invalid are the lines that could not be read; the others were added.
	Invalid []EntryProblem `json:"invalid"`
}

// EntryProblem is one line of an uploaded file that could not be read.
type EntryProblem struct {
	// Index is the line of the file, from 1.
	Index   int    `json:"index"`
	Message string `json:"message"`
}

// Attestation is the terms some formulas need accepted, and whether the account has accepted this version.
type Attestation struct {
	// Kind is what the terms are about, such as "personal_data".
	Kind  string `json:"kind"`
	Title string `json:"title"`
	Text  string `json:"text"`
	// TextVersion is the version of the text; accepting an older one does not count.
	TextVersion string `json:"text_version"`
	// Accepted says whether an owner or an admin accepted this version.
	Accepted   bool       `json:"accepted"`
	AcceptedAt *time.Time `json:"accepted_at"`
	// AcceptedBy is who accepted it; empty when nobody has.
	AcceptedBy string `json:"accepted_by"`
}

// RowFeedback is the answer of [FormulaRowService.Feedback].
type RowFeedback struct {
	rawBody
	RowID   string `json:"row_id"`
	Verdict string `json:"verdict"`
	// Status is the row's status now: new, flagged or rejected.
	Status string `json:"status"`
	// Correction is, for wrong_company, the match this account's runs now use; nil otherwise.
	Correction *RowCorrection `json:"correction"`
}

// RowCorrection is a company match corrected by [VerdictWrongCompany] feedback.
type RowCorrection struct {
	CompanyName string `json:"company_name"`
	// Country is empty when the API sent null.
	Country string `json:"country"`
	// Domain is what this account's runs now match the company to; empty for no domain.
	Domain string `json:"domain"`
}

// decodeCoded reads a JSON object that has a "code": the code, and the whole object as a map.
func decodeCoded(data []byte, code *string, fields *map[string]any) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*fields = m
	*code, _ = m["code"].(string)
	return nil
}

// encodeCoded writes fields with code.
func encodeCoded(code string, fields map[string]any) ([]byte, error) {
	out := make(map[string]any, len(fields)+1)
	for k, v := range fields {
		out[k] = v
	}
	out["code"] = code
	return json.Marshal(out)
}
