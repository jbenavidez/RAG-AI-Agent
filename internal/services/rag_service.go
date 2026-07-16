package services

import (
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

func (s *RagService) AskQuestion(question string) (string, error) {

	return "hello from Agent", nil
}
