package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func NewWeaviateClient() (*weaviate.Client, error) {
	ctx := context.Background()
	weaviateURL := os.Getenv("WEAVIATE_URL")
	if weaviateURL == "" {
		weaviateURL = "http://localhost:8081"
	}

	parsedURL, err := url.Parse(weaviateURL)
	if err != nil {
		return nil, fmt.Errorf("invalid WEAVIATE_URL: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid WEAVIATE_URL: %s", weaviateURL)
	}

	client := weaviate.New(weaviate.Config{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
	})

	// Wait until Weaviate is ready
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		_, err := client.Schema().Getter().Do(ctx)
		if err == nil {
			//break since Weaviate is ready
			break
		}
		fmt.Println("Waiting for Weaviate to be ready...")
		time.Sleep(2 * time.Second)
	}

	// Check one last time
	_, err = client.Schema().Getter().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("weaviate not ready after retries: %w", err)
	}

	// check is class exist
	className := "Document"
	exists, err := client.Schema().ClassExistenceChecker().
		WithClassName(className).
		Do(context.Background())

	if err != nil {
		log.Fatalf("Failed to check class existence: %v", err)
	}
	if !exists {
		fmt.Println("creating class for documents exist")
		// create class
		documentClass := &models.Class{
			Class:       className,
			Description: "Collection of documents",
			Vectorizer:  "text2vec-openai",
			Properties: []*models.Property{
				{Name: "text", DataType: schema.DataTypeText.PropString()},
				{Name: "dateReported", DataType: schema.DataTypeText.PropString()},
				{Name: "projectName", DataType: schema.DataTypeText.PropString()},
				{Name: "description", DataType: schema.DataTypeText.PropString()},
				{Name: "category", DataType: schema.DataTypeText.PropString()},
				{Name: "borough", DataType: schema.DataTypeText.PropString()},
				{Name: "managingAgency", DataType: schema.DataTypeText.PropString()},
				{Name: "clientAgency", DataType: schema.DataTypeText.PropString()},
				{Name: "currentPhase", DataType: schema.DataTypeText.PropString()},
				{Name: "designStart", DataType: schema.DataTypeText.PropString()},
				{Name: "budgetForecast", DataType: schema.DataTypeText.PropString()},
				{Name: "latestBudgetChanges", DataType: schema.DataTypeText.PropString()},
				{Name: "totalBudgetChanges", DataType: schema.DataTypeText.PropString()},
				{Name: "forecastCompletion", DataType: schema.DataTypeText.PropString()},
				{Name: "latestScheduleChanges", DataType: schema.DataTypeText.PropString()},
				{Name: "totalScheduleChanges", DataType: schema.DataTypeText.PropString()},
			},
		}

		if err := client.Schema().ClassCreator().WithClass(documentClass).Do(ctx); err != nil {
			return nil, err
		}
		fmt.Println("class for documents created")
	}

	// check if uploaded file class exists
	uploadedFileClassName := "UploadedFile"
	uploadedFileExists, err := client.Schema().ClassExistenceChecker().
		WithClassName(uploadedFileClassName).
		Do(context.Background())

	if err != nil {
		log.Fatalf("Failed to check uploaded file class existence: %v", err)
	}

	if !uploadedFileExists {
		fmt.Println("creating class for uploaded files")

		uploadedFileClass := &models.Class{
			Class:       uploadedFileClassName,
			Description: "Uploaded file metadata with local file path and processing status",
			Vectorizer:  "none",
			Properties: []*models.Property{
				{Name: "originalFileName", DataType: schema.DataTypeText.PropString()},
				{Name: "storedFileName", DataType: schema.DataTypeText.PropString()},
				{Name: "filePath", DataType: schema.DataTypeText.PropString()},
				{Name: "description", DataType: schema.DataTypeText.PropString()},
				{Name: "contentType", DataType: schema.DataTypeText.PropString()},
				{Name: "status", DataType: schema.DataTypeText.PropString()},
				{Name: "size", DataType: schema.DataTypeInt.PropString()},
				{Name: "createdAt", DataType: schema.DataTypeDate.PropString()},
				{Name: "updatedAt", DataType: schema.DataTypeDate.PropString()},
				{Name: "errorMessage", DataType: schema.DataTypeText.PropString()},
			},
		}

		if err := client.Schema().ClassCreator().WithClass(uploadedFileClass).Do(ctx); err != nil {
			return nil, err
		}
		fmt.Println("class for uploaded files created")
	}

	fmt.Println("Weaviate DB is ready")
	return client, nil
}
