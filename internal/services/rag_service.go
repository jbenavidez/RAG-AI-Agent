package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rag/internal/models"
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
	//Get relevant documents from Weaviate.
	docs, err := s.retriever.GetRelevantDocuments(question)
	if err != nil {
		return "", err
	}
	fmt.Println("valinopr", len(docs))
	if len(docs) == 0 {
		return "I could not find any relevant project information for that question.", nil
	}
	// genarate response
	resp, err := s.GenerateAnswerFromSlides(context.Background(), question, docs)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (s *RagService) GenerateAnswerFromSlides(ctx context.Context, question string, slides []models.CapitalProject) (string, error) {
	// marshal slides
	slidesJson, err := json.Marshal(slides)
	if err != nil {
		return "", err
	}
	// Get promts
	promts, err := rag.BuildRAGPrompt(slidesJson, question)
	if err != nil {
		log.Fatal(err)
	}

	// Get answer
	answer, err := s.rag.Generate(ctx, promts)
	if err != nil {
		return "", err
	}

	return answer, nil
}
