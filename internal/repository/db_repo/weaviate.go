package dbrepo

import (
	"context"
	"fmt"
	"rag/internal/models"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

const (
	timeout                = time.Second * 3
	UploadedFilesClassName = "UploadedFile"
)

type WeaviateDBRepo struct {
	DB *weaviate.Client
}

func NewWeaviateDBRepo(db *weaviate.Client) *WeaviateDBRepo {

	return &WeaviateDBRepo{
		DB: db,
	}

}

func (m *WeaviateDBRepo) GetAllUploadFiles() ([]models.UploadFile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var UploadedFiles []models.UploadFile
	//get files
	res, err := m.DB.Data().ObjectsGetter().WithClassName(UploadedFilesClassName).Do(ctx)
	if err != nil {
		return UploadedFiles, err
	}
	fmt.Println("we were able to queried!!!", res)
	return UploadedFiles, nil
}

func (m *WeaviateDBRepo) GetTotalDocs() (int, error) {

	// TODO: total docs
	return 0, nil
}
