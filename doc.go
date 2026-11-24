// Package enrichednow is the official Go client for the EnrichedNow API.
//
// Create a client, then run a scraper and read its records:
//
//	client, err := enrichednow.NewClient() // reads ENRICHEDNOW_API_KEY
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	run, records, err := client.Scrape(ctx, "ats.board", map[string]any{
//		"board_url": "https://boards.greenhouse.io/acme",
//	})
//	if err != nil {
//		log.Fatal(err) // a *RunFailedError when the scrape itself did not work
//	}
//	fmt.Println(run.Status, run.RecordCount)
//
//	for records.Next() {
//		fmt.Println(string(records.Current()))
//	}
//	if err := records.Err(); err != nil {
//		log.Fatal(err)
//	}
//
// Or run a formula — a ready-made way from a few parameters to a finished list of rows — with its parameter type:
//
//	run, rows, err := client.RunFormula(ctx, "job-order-hunter", enrichednow.JobOrderHunterParams{
//		JobTitles: []string{"Java Developer"},
//		Locations: []string{"London"},
//	}, nil)
//
// The services mirror the API: [Client.Scrapers], [Client.Runs], [Client.Usage], [Client.Data], [Client.Webhooks],
// [Client.Formulas], [Client.FormulaRuns], [Client.SavedFormulas], [Client.Suppressions], [Client.Attestations]
// and [Client.FormulaRows]. Every method takes a [context.Context]. Requests that are safe to repeat are retried
// with exponential backoff; starting a run carries an idempotency key, so that a retry never starts a second run.
//
// Errors answered by the API are [*APIError] values and match the sentinel errors of this package with
// [errors.Is]: [ErrAuthentication], [ErrNotFound], [ErrValidation], [ErrRateLimit], [ErrServer] and the others.
package enrichednow

// Version is the version of this SDK. It is sent in the User-Agent header.
const Version = "0.4.0"
