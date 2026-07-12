package processors

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"rag/internal/models"
	dbrepo "rag/internal/repository/db_repo"
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
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read csv file: %w", err)
	}

	fmt.Println("csv records:", records)

	return nil
}
