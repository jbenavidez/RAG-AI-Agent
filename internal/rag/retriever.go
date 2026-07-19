package rag

import (
	"rag/internal/models"
	"rag/internal/repository"
)

type Retriever struct {
	repo repository.DatabaseRepo
}

func NewRetriever(repo repository.DatabaseRepo) *Retriever {
	return &Retriever{
		repo: repo,
	}
}

func (r *Retriever) GetRelevantDocuments(question string) ([]models.CapitalProject, error) {
	return r.repo.GetDocuments(question)
}
