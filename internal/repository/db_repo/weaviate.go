package dbrepo

import (
	"context"
	"fmt"
	"rag/internal/models"
	"rag/internal/storage"
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
func (m *WeaviateDBRepo) SaveFileMetaData(storedFile *storage.StoredFile) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	now := time.Now().UTC().Format(time.RFC3339)
	fileProperties := map[string]interface{}{
		"originalFileName": storedFile.OriginalFileName,
		"storedFileName":   storedFile.StoredFileName,
		"filePath":         storedFile.FilePath,
		"description":      storedFile.Description,
		"contentType":      storedFile.ContentType,
		"status":           "pending",
		"size":             storedFile.Size,
		"createdAt":        now,
		"updatedAt":        now,
	}
	_, err := m.DB.Data().Creator().
		WithClassName(UploadedFilesClassName).
		WithProperties(fileProperties).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("failed to save uploaded file metadata in weaviate: %w", err)
	}
	return nil

}
func (m *WeaviateDBRepo) GetAllUploadFiles() ([]models.UploadedFile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := m.DB.Data().ObjectsGetter().
		WithClassName(UploadedFilesClassName).
		Do(ctx)

	if err != nil {
		return nil, err
	}
	uploadedFiles := make([]models.UploadedFile, 0, len(res))
	for i := range res {
		obj := res[i]
		props, ok := obj.Properties.(map[string]interface{})
		if !ok {
			continue
		}
		uploadfile := models.UploadedFile{
			ID:               obj.ID.String(),
			OriginalFileName: getStringProperty(props, "originalFileName"),
			StoredFileName:   getStringProperty(props, "storedFileName"),
			FilePath:         getStringProperty(props, "filePath"),
			Description:      getStringProperty(props, "description"),
			ContentType:      getStringProperty(props, "contentType"),
			Status:           getStringProperty(props, "status"),
			Size:             getInt64Property(props, "size"),
			CreatedAt:        getTimeProperty(props, "createdAt"),
			UpdatedAt:        getTimeProperty(props, "updatedAt"),
			ErrorMessage:     getStringProperty(props, "errorMessage"),
		}
		uploadedFiles = append(uploadedFiles, uploadfile)
	}

	return uploadedFiles, nil
}

func (m *WeaviateDBRepo) GetTotalDocs() (int, error) {

	// TODO: total docs
	return 0, nil
}

func getStringProperty(props map[string]interface{}, key string) string {
	value, ok := props[key]
	if !ok || value == nil {
		return ""
	}

	strValue, ok := value.(string)
	if !ok {
		return ""
	}

	return strValue
}

func getInt64Property(props map[string]interface{}, key string) int64 {
	value, ok := props[key]
	if !ok || value == nil {
		return 0
	}

	switch v := value.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	default:
		return 0
	}
}

func getTimeProperty(props map[string]interface{}, key string) time.Time {
	value, ok := props[key]
	if !ok || value == nil {
		return time.Time{}
	}

	strValue, ok := value.(string)
	if !ok {
		return time.Time{}
	}

	parsedTime, err := time.Parse(time.RFC3339, strValue)
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}
