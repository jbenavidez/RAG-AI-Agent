package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"ragService/models"
	"strings"
)

func (c *Config) LoadDataSet() error {
	totalDocs, err := c.WDBRepo.GetTotalDocs()
	if err != nil {
		return err
	}
	fmt.Println("******** total docs ********", totalDocs)
	if totalDocs > 0 {
		fmt.Printf(" the total docs %v", totalDocs)
		return nil
	}
	fmt.Println("******** Getting  data from csv ********")
	docs, err := c.GetData()
	if err != nil {
		return err
	}
	fmt.Printf("******** Total rows to insert %v ********", len(docs))

	err = c.WDBRepo.InsertDocuments(docs)
	if err != nil {
		return err
	}
	fmt.Println("******** Data load completed ********")
	return nil

}

func (c *Config) GetData() ([]*models.Doc, error) {
	f, err := os.Open("./cmd/data/data.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var docs []*models.Doc
	// skip ehader
	for _, row := range rows[1:] {
		// Combine key fields for  search
		text := strings.TrimSpace(fmt.Sprintf("%s — %s", row[2], row[3])) // ProjectName — Description

		d := models.Doc{
			ID:                    row[1],
			DateReported:          row[0],
			ProjectName:           row[2],
			Description:           row[3],
			Category:              row[4],
			Borough:               row[5],
			ManagingAgency:        row[6],
			ClientAgency:          row[7],
			CurrentPhase:          row[8],
			DesignStart:           row[9],
			BudgetForecast:        row[10],
			LatestBudgetChanges:   row[11],
			TotalBudgetChanges:    row[12],
			ForecastCompletion:    row[13],
			LatestScheduleChanges: row[14],
			TotalScheduleChanges:  row[15],
			Text:                  text,
		}

		docs = append(docs, &d)
	}

	return docs, nil
}

func (c *Config) DocsToContext(docs []*models.Doc) string {
	//return tempty string is there are not docs
	if len(docs) == 0 {
		return ""
	}
	// Build context from relevant chunks
	context := "Context from documents:\n"
	for i, doc := range docs {
		context += fmt.Sprintf(`
					%d
					Reported Date  "%s"
					Project Name: "%s"
					Description: "%s"
					Borough : "%s"
					Managing Agency : "%s"
					Client Agency : "%s"
					Current Phase : "%s"
					Design Start : "%s"
					Budget Forecast : "%s"
					Latest Budget Changes : "%s"
					Total Budget Changes : "%s"
					Forecast Completion : "%s"
					Latest Schedule Changes : "%s"
					Total Schedule Changes : "%s"
					%s
					\n\n
					`,
			i+1,
			doc.DateReported,
			doc.ProjectName,
			doc.Description,
			doc.Borough,
			doc.ManagingAgency,
			doc.ClientAgency,
			doc.CurrentPhase,
			doc.DesignStart,
			doc.BudgetForecast,
			doc.LatestBudgetChanges,
			doc.TotalBudgetChanges,
			doc.ForecastCompletion,
			doc.LatestScheduleChanges,
			doc.TotalScheduleChanges,
		)
	}
	return context
}
