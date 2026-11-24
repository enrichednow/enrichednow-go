# Changelog

What changed in each version of the module. Versions before 0.4.0 are in the API's changelog,
<https://docs.enrichednow.com/changelog/>.

## 0.4.0 (not released yet)

Formulas: ready-made workflows that turn a few parameters into a finished list of rows.

- `client.Formulas`: `List` (`Query` in your own words, best matches first; each formula's `Basics`, the questions it
  asks), `Get`, `Estimate` (free; parameters that cannot run come back with `Valid` false and
  what to change) and `Run` (`FormulaRunParams`: preview, unit cap, row cap, wait, and an `Idempotency-Key` on every
  start, as on a scraper run).
- `client.FormulaRuns`: `List`/`ListAll`, `Get`, `GetWait`, `Wait` (each look held by the API for up to a minute),
  `RowsPage`/`Rows` across pages (`FormulaRowListParams`: fields, least score, one row per company), `Export` and
  `ExportAll` as CSV or NDJSON, `Cancel`.
- `client.RunFormula(ctx, id, params, opts)`: starts a run, waits for it and returns it with an iterator over its
  rows; a `*FormulaRunFailedError` (matching `ErrRunFailed`) when the run fails.
- `client.SavedFormulas`: `Create` (with a schedule), `List`/`ListAll`, `Get`, `Update` (`RemoveSchedule`,
  `RemoveMonthlyUnitCap`, `AutoName`), `Archive`, `Duplicate`, `Run`, `Pause`, `Resume`, `ResetSeen`, rows since a
  time, exports, `Snapshots`, `Limits`.
- `client.Suppressions` (`List`, `Add`, `Upload` a CSV file, `Delete`), `client.Attestations.List` and
  `client.FormulaRows.Feedback`.
- Each formula's parameters as a struct (`JobOrderHunterParams`, …, with `Ptr` for its optional fields); a struct
  for another formula than the one named is refused before anything is sent.
- New sentinels: `ErrPlanLimitReached` (also `ErrPermission`), `ErrAttestationRequired`, `ErrFormulaUnavailable`,
  `ErrRequirementMissing`, `ErrSavedFormulaArchived` (also `ErrConflict`), and the codes of each.
