package services

import (
	"context"
	"rag/internal/rag"
	dbrepo "rag/internal/repository/db_repo"
)

type RagService struct {
	rag  *rag.Rag
	repo *dbrepo.WeaviateDBRepo
}

func NewRagService(r *rag.Rag, repo *dbrepo.WeaviateDBRepo) *RagService {
	return &RagService{
		rag:  r,
		repo: repo,
	}
}

func (s *RagService) AskQuestion(ctx context.Context, question string) (string, error) {

	return "nil", nil
}
