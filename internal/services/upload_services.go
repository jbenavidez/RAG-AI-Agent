package services

import (
	"context"
	"mime/multipart"
	dbrepo "rag/internal/repository/db_repo"
	"rag/internal/storage"
)

type UploadService struct {
	repo    *dbrepo.WeaviateDBRepo
	storage storage.FileStorage
}

func NewUploadService(r *dbrepo.WeaviateDBRepo, s storage.FileStorage) *UploadService {
	return &UploadService{
		repo:    r,
		storage: s,
	}
}

func (s *UploadService) SaveUploadedFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, description string) (*storage.StoredFile, error) {
	//save file
	storedFile, err := s.storage.Save(ctx, file, header)
	if err != nil {
		return nil, err
	}
	return storedFile, nil
}
