package enrichednow_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	enrichednow "github.com/enrichednow/enrichednow-go"
)

// Run a scraper to the end and read its records.
func ExampleClient_Scrape() {
	client, err := enrichednow.NewClient() // reads ENRICHEDNOW_API_KEY
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	run, records, err := client.Scrape(ctx, "ats.board", map[string]any{
		"board_url": "https://boards.greenhouse.io/acme",
	})
	var failed *enrichednow.RunFailedError
	switch {
	case errors.As(err, &failed):
		log.Fatalf("the run ended %s: %s", failed.Status(), failed.Reason())
	case err != nil:
		log.Fatal(err)
	}
	fmt.Println(run.ID, run.Status, run.RecordCount)

	type Job struct {
		Title    string `json:"title"`
		ApplyURL string `json:"apply_url"`
	}
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

// Start a run, do something else, and come back for it.
func ExampleRunService_Wait() {
	client, err := enrichednow.NewClient(enrichednow.WithAPIKey("en_live_..."))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	run, err := client.Scrapers.Run(ctx, "careers.company", map[string]any{"domain": "acme.com"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	run, err = client.Runs.Wait(ctx, run.ID, &enrichednow.WaitParams{
		OnPoll: func(r *enrichednow.Run) { log.Println(r.ID, r.Status) },
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := run.Err(); err != nil {
		log.Fatal(err) // failed, unsupported or cancelled
	}
	fmt.Println(run.RecordCount, "records for", run.Units, "units")
}

// Tell the kinds of error apart.
func ExampleAPIError() {
	client, err := enrichednow.NewClient(enrichednow.WithAPIKey("en_live_..."))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Scrapers.Run(context.Background(), "ats.board", map[string]any{"board": "acme"}, nil)

	var apiErr *enrichednow.APIError
	switch {
	case errors.Is(err, enrichednow.ErrValidation) && errors.As(err, &apiErr):
		for _, field := range apiErr.Fields() {
			fmt.Println(field.Path, field.Message)
		}
	case errors.Is(err, enrichednow.ErrQuotaExceeded):
		fmt.Println("this month's units are used up")
	case errors.Is(err, enrichednow.ErrTooManyRuns):
		fmt.Println("wait for one of your runs to end")
	case errors.Is(err, enrichednow.ErrAuthentication):
		fmt.Println("check the API key")
	case err != nil:
		log.Fatal(err)
	}
}

// Search the job postings dataset.
func ExamplePostingService_List() {
	client, err := enrichednow.NewClient(enrichednow.WithAPIKey("en_live_..."))
	if err != nil {
		log.Fatal(err)
	}
	page, err := client.Data.Postings.List(context.Background(), &enrichednow.PostingListParams{
		Title:   "data engineer",
		Country: "US",
		Status:  "open",
		Limit:   20,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, posting := range page.Postings {
		fmt.Println(posting.Title, posting.URL)
	}
	fmt.Println("next page:", page.NextCursor)
}
