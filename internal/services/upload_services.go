package services

import (
	dbrepo "rag/internal/repository/db_repo"
)

type UploadService struct {
	repo *dbrepo.WeaviateDBRepo
}

func NewUploadService(r *dbrepo.WeaviateDBRepo) *UploadService {

	return &UploadService{
		repo: r,
	}
}
