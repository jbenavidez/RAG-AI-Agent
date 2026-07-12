package repository

import "rag/internal/models"

type DatabaseRepo interface {
	GetTotalDocs() (int, error)
	GetAllUploadFiles() ([]models.UploadFile, error)
}
