package services

import (
	"fmt"
	"rag/internal/rag"
	dbrepo "rag/internal/repository/db_repo"
)

type RagService struct {
	rag       *rag.Rag
	retriever *rag.Retriever
	repo      *dbrepo.WeaviateDBRepo
}

func NewRagService(r *rag.Rag, repo *dbrepo.WeaviateDBRepo) *RagService {
	return &RagService{
		rag:       r,
		repo:      repo,
		retriever: rag.NewRetriever(repo),
	}
}

func (s *RagService) AskQuestion(question string) (string, error) {
	docs, err := s.retriever.GetRelevantDocuments(question)
	if err != nil {
		return "", err
	}
	fmt.Println("valinopr", len(docs))
	if len(docs) == 0 {
		return "I could not find any relevant project information for that question.", nil
	}

	return "", nil
}
