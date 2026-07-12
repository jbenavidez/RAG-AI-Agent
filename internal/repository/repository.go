package repository

import (
	"rag/internal/models"
	"rag/internal/storage"
)

type DatabaseRepo interface {
	SaveFileMetaData(storedFile *storage.StoredFile) (*models.UploadedFile, error)
	GetTotalDocs() (int, error)
	GetAllUploadFiles() ([]models.UploadedFile, error)
	GetFileByName(fileName string) (*models.UploadedFile, error)
	UpdateFile(uploadedfile *models.UploadedFile) error
}
