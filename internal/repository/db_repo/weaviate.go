package dbrepo

import "github.com/weaviate/weaviate-go-client/v5/weaviate"

type WeaviateDBRepo struct {
	DB *weaviate.Client
}

func (m *WeaviateDBRepo) Connection() *weaviate.Client {
	return m.DB
}
