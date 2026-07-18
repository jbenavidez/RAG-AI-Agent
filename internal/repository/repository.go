package repository

import (
	"rag/internal/models"
	"rag/internal/storage"
)

type DatabaseRepo interface {
	SaveFileMetaData(*storage.StoredFile) (*models.UploadedFile, error)
	GetTotalDocs() (int, error)
	GetAllUploadFiles() ([]models.UploadedFile, error)
	GetFileByName(string) (*models.UploadedFile, bool, error)
	UpdateFile(*models.UploadedFile) error
	InsertDocuments([]*models.CapitalProject) error
	GetDocuments(string) ([]models.CapitalProject, error)
}
