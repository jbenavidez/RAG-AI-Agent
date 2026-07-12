package dbrepo

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

type WeaviateDBRepo struct {
	DB *weaviate.Client
}

func NewWeaviateDBRepo(db *weaviate.Client) *WeaviateDBRepo {

	return &WeaviateDBRepo{
		DB: db,
	}

}

func (m *WeaviateDBRepo) GetTotalDocs() (int, error) {

	// TODO: total docs
	return 0, nil
}

func (m *WeaviateDBRepo) Savefile() error {

	//save files
	return nil
}
