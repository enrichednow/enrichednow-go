package enrichednow_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	enrichednow "github.com/enrichednow/enrichednow-go"
)

// Run a formula to the end with its parameter type, and read its rows.
func ExampleClient_RunFormula() {
	client, err := enrichednow.NewClient() // reads ENRICHEDNOW_API_KEY
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	run, rows, err := client.RunFormula(ctx, "job-order-hunter", enrichednow.JobOrderHunterParams{
		JobTitles:        []string{"Java Developer"},
		Locations:        []string{"London"},
		PostedWithinDays: enrichednow.Ptr(7),
	}, &enrichednow.RunFormulaParams{UnitCap: 2000})
	var failed *enrichednow.FormulaRunFailedError
	switch {
	case errors.As(err, &failed):
		log.Fatalf("the run ended %s: %s", failed.Status(), failed.Reason())
	case err != nil:
		log.Fatal(err)
	}
	fmt.Println(run.ID, run.Status, run.RowsTotal, "rows for", run.UnitsSpent, "units")

	type Lead struct {
		CompanyName  string `json:"company_name"`
		JobTitle     string `json:"job_title"`
		ContactEmail string `json:"contact_email"`
	}
	for rows.Next() {
		lead, err := enrichednow.Decode[Lead](rows.Current().Data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(lead.CompanyName, lead.JobTitle, lead.ContactEmail)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// The whole run as one CSV file, as the API writes it.
	f, err := os.Create("leads.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := client.FormulaRuns.ExportAll(ctx, run.ID, enrichednow.FormulaExportParams{Format: enrichednow.FormatCSV}, f); err != nil {
		log.Fatal(err)
	}
}

// Ask what a run would cost before starting it, then try a preview.
func ExampleFormulaService_Estimate() {
	client, err := enrichednow.NewClient(enrichednow.WithAPIKey("en_live_..."))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	params := map[string]any{"job_titles": []string{"Registered Nurse"}, "locations": []string{"Texas"}}

	estimate, err := client.Formulas.Estimate(ctx, "job-order-hunter", params, nil)
	if err != nil {
		log.Fatal(err)
	}
	if !estimate.Valid {
		for _, problem := range estimate.Errors {
			fmt.Println(problem.Path, problem.Message)
		}
		return
	}
	fmt.Println("about", estimate.Estimate.ExpectedUnits, "units, at most", estimate.Estimate.MaxUnits)

	preview, err := client.Formulas.Run(ctx, "job-order-hunter", params, &enrichednow.FormulaRunParams{Preview: true, Wait: time.Minute})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(preview.ID, preview.Status, preview.RowsTotal)
}

// Save a formula with a schedule, then poll its new rows.
func ExampleSavedFormulaService_Create() {
	client, err := enrichednow.NewClient(enrichednow.WithAPIKey("en_live_..."))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	saved, err := client.SavedFormulas.Create(ctx, "job-order-hunter", enrichednow.JobOrderHunterParams{
		JobTitles: []string{"Java Developer"},
		Locations: []string{"London"},
	}, &enrichednow.SavedFormulaCreateParams{
		Schedule:       &enrichednow.ScheduleParams{Cadence: enrichednow.CadenceWeekdays, LocalTime: "07:30", Timezone: "Europe/London"},
		MonthlyUnitCap: 20000,
		Notify:         enrichednow.NotifyDigest,
	})
	if errors.Is(err, enrichednow.ErrPlanLimitReached) {
		log.Fatal("the plan does not allow this schedule: ", err)
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Println(saved.ID, saved.Name, saved.Schedule.NextFireAt)

	// Later, and every day after: the rows written since the last look.
	since := time.Now().Add(-24 * time.Hour)
	rows := client.SavedFormulas.Rows(ctx, saved.ID, &enrichednow.SavedFormulaRowListParams{Since: since})
	for rows.Next() {
		fmt.Println(rows.Current().ID, string(rows.Current().Data))
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

// Keep your customers and those who opted out off every formula's rows.
func ExampleSuppressionService_Upload() {
	client, err := enrichednow.NewClient(enrichednow.WithAPIKey("en_live_..."))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	added, err := client.Suppressions.Add(ctx, []enrichednow.SuppressionEntry{
		{Kind: enrichednow.SuppressDomain, Value: "globex.example"},
		{Kind: enrichednow.SuppressEmail, Value: "ana@initech.example"},
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(added, "added")

	file := strings.NewReader("kind,value\ndomain,hooli.example\ncompany_name,Pied Piper\n")
	upload, err := client.Suppressions.Upload(ctx, file, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(upload.Added, "added,", upload.Already, "already there")
	for _, line := range upload.Invalid {
		fmt.Println("line", line.Index, line.Message)
	}
}

// Say a row matched the wrong company.
func ExampleFormulaRowService_Feedback() {
	client, err := enrichednow.NewClient(enrichednow.WithAPIKey("en_live_..."))
	if err != nil {
		log.Fatal(err)
	}
	feedback, err := client.FormulaRows.Feedback(context.Background(), "frow_000000000000000000000001", enrichednow.VerdictWrongCompany,
		&enrichednow.RowFeedbackParams{CorrectDomain: "corvid.example"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(feedback.Status, feedback.Correction.Domain)
}

// Tell the errors of the formulas apart.
func ExampleAPIError_formulas() {
	client, err := enrichednow.NewClient(enrichednow.WithAPIKey("en_live_..."))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Formulas.Run(context.Background(), "job-order-hunter", map[string]any{"job_titles": []string{"Nurse"}, "locations": []string{"Leeds"}}, nil)

	var apiErr *enrichednow.APIError
	errors.As(err, &apiErr)
	switch {
	case errors.Is(err, enrichednow.ErrAttestationRequired):
		fmt.Println("an owner or an admin must accept the terms in the console first")
	case errors.Is(err, enrichednow.ErrRequirementMissing):
		fmt.Println("the account lacks", apiErr.Missing())
	case errors.Is(err, enrichednow.ErrFormulaUnavailable):
		fmt.Println("the formula cannot be run now")
	case errors.Is(err, enrichednow.ErrValidation):
		if least, ok := apiErr.MinUnitCap(); ok {
			fmt.Println("give a unit cap of at least", least)
		}
		for _, field := range apiErr.Fields() {
			fmt.Println(field.Path, field.Message)
		}
	case errors.Is(err, enrichednow.ErrLimitReached):
		fmt.Println("a spend limit is reached:", apiErr.Details["scope"])
	case err != nil:
		log.Fatal(err)
	}
}
