package services

import (
	"rag/internal/rag"
	dbrepo "rag/internal/repository/db_repo"
)

type RagService struct {
	Rag  rag.Rag
	repo *dbrepo.WeaviateDBRepo
}
