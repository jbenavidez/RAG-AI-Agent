package rag

import (
	"fmt"
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
	fmt.Println("valinor_1")
	return r.repo.GetDocuments(question)
}
