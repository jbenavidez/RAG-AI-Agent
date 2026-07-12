package repository

import (
	"rag/internal/models"
	"rag/internal/storage"
)

type DatabaseRepo interface {
	SaveFileMetaData(storedFile *storage.StoredFile) error
	GetTotalDocs() (int, error)
	GetAllUploadFiles() ([]models.UploadedFile, error)
	GetFileByName(fileName string) (*models.UploadedFile, error)
}
