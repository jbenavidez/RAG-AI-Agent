package dbrepo

import (
	"context"
	"fmt"
	"rag/internal/models"
	"rag/internal/storage"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	mmodels "github.com/weaviate/weaviate/entities/models"
)

const (
	timeout                = time.Second * 3
	UploadedFilesClassName = "UploadedFile"
	DocumentClassName      = "Document"
)

type WeaviateDBRepo struct {
	DB *weaviate.Client
}

func NewWeaviateDBRepo(db *weaviate.Client) *WeaviateDBRepo {

	return &WeaviateDBRepo{
		DB: db,
	}

}
func (m *WeaviateDBRepo) SaveFileMetaData(storedFile *storage.StoredFile) (*models.UploadedFile, error) {

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
	resp, err := m.DB.Data().Creator().
		WithClassName(UploadedFilesClassName).
		WithProperties(fileProperties).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to save uploaded file metadata in weaviate: %w", err)
	}
	// get properties
	props := resp.Object.Properties.(map[string]interface{})
	uploadFiled := &models.UploadedFile{
		ID:               resp.Object.ID.String(),
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
	return uploadFiled, nil

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

func (m *WeaviateDBRepo) GetFileByName(fileName string) (*models.UploadedFile, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	where := filters.Where().
		WithPath([]string{"originalFileName"}).
		WithOperator(filters.Equal).
		WithValueString(fileName)

	res, err := m.DB.GraphQL().Get().
		WithClassName(UploadedFilesClassName).
		WithFields(
			graphql.Field{Name: "originalFileName"},
			graphql.Field{Name: "storedFileName"},
			graphql.Field{Name: "filePath"},
			graphql.Field{Name: "description"},
			graphql.Field{Name: "contentType"},
			graphql.Field{Name: "status"},
			graphql.Field{Name: "size"},
			graphql.Field{Name: "createdAt"},
			graphql.Field{Name: "updatedAt"},
			graphql.Field{Name: "errorMessage"},
			graphql.Field{
				Name: "_additional",
				Fields: []graphql.Field{
					{Name: "id"},
				},
			},
		).
		WithWhere(where).
		WithLimit(1).
		Do(ctx)

	if err != nil {
		return nil, false, err
	}

	data, ok := res.Data["Get"].(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("invalid weaviate response: missing Get")
	}

	items, ok := data[UploadedFilesClassName].([]interface{})
	if !ok || len(items) == 0 {
		return nil, false, nil
	}

	props, ok := items[0].(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("invalid uploaded file properties")
	}

	uploadedFile := models.UploadedFile{
		ID:               getGraphQLID(props),
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

	return &uploadedFile, true, nil
}

func (m *WeaviateDBRepo) UpdateFile(uploadedFile *models.UploadedFile) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fileProperties := map[string]interface{}{
		"originalFileName": uploadedFile.OriginalFileName,
		"storedFileName":   uploadedFile.StoredFileName,
		"filePath":         uploadedFile.FilePath,
		"description":      uploadedFile.Description,
		"contentType":      uploadedFile.ContentType,
		"status":           uploadedFile.Status,
		"size":             uploadedFile.Size,
		"updatedAt":        time.Now().UTC().Format(time.RFC3339),
		"errorMessage":     uploadedFile.ErrorMessage,
	}

	err := m.DB.Data().Updater().
		WithMerge().
		WithID(uploadedFile.ID).
		WithClassName(UploadedFilesClassName).
		WithProperties(fileProperties).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("failed to update uploaded file: %w", err)
	}

	return nil
}

func (m *WeaviateDBRepo) InsertDocuments(docs []*models.CapitalProject) error {
	const batchSize = 50

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if len(docs) == 0 {
		return nil
	}

	batcher := m.DB.Batch().ObjectsBatcher()
	count := 0

	for i, doc := range docs {
		if doc == nil {
			continue
		}

		obj := &mmodels.Object{
			Class: DocumentClassName,
			Properties: map[string]interface{}{
				"text":                  doc.Text,
				"dateReported":          doc.DateReported,
				"projectName":           doc.ProjectName,
				"description":           doc.Description,
				"category":              doc.Category,
				"borough":               doc.Borough,
				"managingAgency":        doc.ManagingAgency,
				"clientAgency":          doc.ClientAgency,
				"currentPhase":          doc.CurrentPhase,
				"designStart":           doc.DesignStart,
				"budgetForecast":        doc.BudgetForecast,
				"latestBudgetChanges":   doc.LatestBudgetChanges,
				"totalBudgetChanges":    doc.TotalBudgetChanges,
				"forecastCompletion":    doc.ForecastCompletion,
				"latestScheduleChanges": doc.LatestScheduleChanges,
				"totalScheduleChanges":  doc.TotalScheduleChanges,
			},
		}

		batcher = batcher.WithObjects(obj)
		count++

		if count >= batchSize || i == len(docs)-1 {
			_, err := batcher.Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to insert document batch: %w", err)
			}

			fmt.Printf("Inserted %d documents\n", count)

			batcher = m.DB.Batch().ObjectsBatcher()
			count = 0

			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil
}

func (m *WeaviateDBRepo) GetDocuments(question string) ([]models.CapitalProject, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nearText := m.DB.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{question})

	res, err := m.DB.GraphQL().Get().
		WithClassName(DocumentClassName).
		WithFields(
			graphql.Field{Name: "text"},
			graphql.Field{Name: "dateReported"},
			graphql.Field{Name: "projectName"},
			graphql.Field{Name: "description"},
			graphql.Field{Name: "category"},
			graphql.Field{Name: "borough"},
			graphql.Field{Name: "managingAgency"},
			graphql.Field{Name: "clientAgency"},
			graphql.Field{Name: "currentPhase"},
			graphql.Field{Name: "designStart"},
			graphql.Field{Name: "budgetForecast"},
			graphql.Field{Name: "latestBudgetChanges"},
			graphql.Field{Name: "totalBudgetChanges"},
			graphql.Field{Name: "forecastCompletion"},
			graphql.Field{Name: "latestScheduleChanges"},
			graphql.Field{Name: "totalScheduleChanges"},
		).
		WithNearText(nearText).
		WithLimit(20).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documents from weaviate: %w", err)
	}

	if len(res.Errors) > 0 {
		for i, gqlErr := range res.Errors {
			if gqlErr == nil {
				continue
			}

			fmt.Printf("weaviate graphql error %d: %s\n", i, gqlErr.Message)
		}

		return nil, fmt.Errorf("weaviate graphql error: %s", res.Errors[0].Message)
	}

	var results []models.CapitalProject

	getData, ok := res.Data["Get"].(map[string]interface{})
	if !ok {
		return results, fmt.Errorf("invalid weaviate response: missing Get data")
	}

	docs, ok := getData[DocumentClassName].([]interface{})
	if !ok {
		return results, fmt.Errorf("invalid weaviate response: missing %s documents", DocumentClassName)
	}

	for _, doc := range docs {
		d, ok := doc.(map[string]interface{})
		if !ok {
			continue
		}

		results = append(results, models.CapitalProject{
			Text:                  getStringProperty(d, "text"),
			DateReported:          getStringProperty(d, "dateReported"),
			ProjectName:           getStringProperty(d, "projectName"),
			Description:           getStringProperty(d, "description"),
			Category:              getStringProperty(d, "category"),
			Borough:               getStringProperty(d, "borough"),
			ManagingAgency:        getStringProperty(d, "managingAgency"),
			ClientAgency:          getStringProperty(d, "clientAgency"),
			CurrentPhase:          getStringProperty(d, "currentPhase"),
			DesignStart:           getStringProperty(d, "designStart"),
			BudgetForecast:        getStringProperty(d, "budgetForecast"),
			LatestBudgetChanges:   getStringProperty(d, "latestBudgetChanges"),
			TotalBudgetChanges:    getStringProperty(d, "totalBudgetChanges"),
			ForecastCompletion:    getStringProperty(d, "forecastCompletion"),
			LatestScheduleChanges: getStringProperty(d, "latestScheduleChanges"),
			TotalScheduleChanges:  getStringProperty(d, "totalScheduleChanges"),
		})
	}

	fmt.Printf("retrieved documents from weaviate: %d\n", len(results))

	return results, nil
}

// HELPERS
func getGraphQLID(props map[string]interface{}) string {
	additional, ok := props["_additional"].(map[string]interface{})
	if !ok {
		return ""
	}

	return getStringProperty(additional, "id")
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
