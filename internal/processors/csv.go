package processors

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"rag/internal/models"
	dbrepo "rag/internal/repository/db_repo"
	"strings"
)

type CSVProcessor struct {
	repo      *dbrepo.WeaviateDBRepo
	chunkSize int
}

func NewCSVProcessor(repo *dbrepo.WeaviateDBRepo, chunkSize int) *CSVProcessor {
	return &CSVProcessor{
		repo:      repo,
		chunkSize: chunkSize,
	}
}

func (p *CSVProcessor) Process(ctx context.Context, file *os.File, uploadedFile *models.UploadedFile) error {
	fmt.Println("******** Init Processing file ********")
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read csv file: %w", err)
	}

	rows := records[1:] // skip header row
	capitalProjects := make([]*models.CapitalProject, len(rows))
	// read data
	for i := range rows {
		row := rows[i]
		// create text
		text := strings.TrimSpace(fmt.Sprintf("%s — %s", row[2], row[3]))
		capitalProjects[i] = &models.CapitalProject{
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

	}
	// store data
	fmt.Printf("******** Total capital projects to insert %v ********", len(capitalProjects))
	if err := p.repo.InsertDocuments(capitalProjects); err != nil {
		return err
	}

	fmt.Println("valinor", capitalProjects[0])

	return nil
}
