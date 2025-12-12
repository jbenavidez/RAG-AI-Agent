package repository

import (
	"client/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	GetTotalDocuments() (int, error)
	InsertDocument(documents []models.Document) error
	GetEmbeddingDocument(queryText []float32, topK int, keyword string) ([]string, error)
}
