package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"ragAIAgent/models"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

func (c *RagConfig) LoadDataSet() error {
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

func (c *RagConfig) GetData() ([]*models.Doc, error) {
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

func (c *RagConfig) DocsToContext(docs []*models.Doc) string {
	//return tempty string is there are not docs
	if len(docs) == 0 {
		return ""
	}
	// Build context from relevant chunks
	context := "Context from documents:\n"
	for _, doc := range docs {
		context += fmt.Sprintf(`
					Reported Date %s\n\n Project Name: %s\n\n Description: %s\n\n Borough : %s\n\n Managing Agency : %s\n\n Client Agency : %s\n\n Current Phase : %s\n\n Design Start : %s\n\n Budget Forecast : %s\n\n Latest Budget Changes : %s\n\n Total Budget Changes : %s\n\n Forecast Completion : %s\n\nLatest Schedule Changes : %s\n\nTotal Schedule Changes : %s\n\n
					`,
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

func (c *RagConfig) GenerateAnswerFromSlides(ctx context.Context, question string, slides []*models.Doc) (string, error) {
	// marshal docks
	slidesJSON, err := json.Marshal(slides)
	if err != nil {
		return "", err
	}

	// prompt required to generate a answer
	prompt := fmt.Sprintf(
		"Answer the user's question concisely, professionally, and in a friendly manner using the project information below. "+
			"Do NOT mention the data source or any JSON/slides. "+
			"Number items if there are multiple. "+
			"Optionally suggest one or two helpful follow-up questions.\n\n"+
			"Project Data:\n%s\n\nQuestion: %s\nAnswer:",
		slidesJSON, question,
	)
	// let llm generate response
	response, err := c.Llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	})
	if err != nil {
		return "", err
	}

	// return error is now response
	if len(response.Choices) == 0 {
		return "", errors.New("no response from LLM")
	}
	return strings.TrimSpace(response.Choices[0].Content), nil
}
