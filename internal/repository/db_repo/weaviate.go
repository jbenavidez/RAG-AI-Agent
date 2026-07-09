package dbrepo

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

type WeaviateDBRepo struct {
	DB *weaviate.Client
}

func (m *WeaviateDBRepo) GetTotalDocs() (int, error) {

	// TODO: total docs
	return 0, nil
}
