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
// The services mirror the API: [Client.Scrapers], [Client.Runs], [Client.Usage] and [Client.Data]. Every method
// takes a [context.Context]. Requests that are safe to repeat are retried with exponential backoff; starting a run
// is never retried, because a repeated start is a second run.
//
// Errors answered by the API are [*APIError] values and match the sentinel errors of this package with
// [errors.Is]: [ErrAuthentication], [ErrNotFound], [ErrValidation], [ErrRateLimit], [ErrServer] and the others.
package enrichednow

// Version is the version of this SDK. It is sent in the User-Agent header.
const Version = "0.1.1"
