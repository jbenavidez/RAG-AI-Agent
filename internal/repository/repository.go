package repository

import (
	"rag/internal/models"
	"rag/internal/storage"
)

type DatabaseRepo interface {
	SaveFileMetaData(storedFile *storage.StoredFile) (*models.UploadedFile, error)
	GetTotalDocs() (int, error)
	GetAllUploadFiles() ([]models.UploadedFile, error)
	GetFileByName(string) (*models.UploadedFile, bool, error)
	UpdateFile(uploadedfile *models.UploadedFile) error
	InsertDocuments(data []*models.CapitalProject) error
}
