package main

import (
	"client/models"
	pb "client/proto/generated"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func (c *Config) LoadData() error {
	totalDocs, err := c.DB.GetTotalDocuments()
	if err != nil {
		return err
	}
	fmt.Println("Total documents on DB", totalDocs)
	// only inser tada if total_doc is equal to 0
	if totalDocs > 0 {
		return nil
	}
	f, err := os.Open("./cmd/data/data.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	rows, err := r.ReadAll()
	if err != nil {
		return err
	}

	const batchSize = 100
	var batch []models.Document
	for _, row := range rows[1:] {
		projectcName := row[1]
		description := row[2]
		text := strings.TrimSpace(fmt.Sprintf("%s — %s", projectcName, description))
		fmt.Println("gondor", text)
		doc := models.Document{
			Text:        text,
			ProjectName: projectcName,
			Description: description,
		}
		batch = append(batch, doc)

		if len(batch) >= batchSize {
			err := c.DB.InsertDocument(batch)
			if err != nil {
				return err
			}
			log.Printf("Inserted batch of %d", len(batch))
			batch = batch[:0] // reset batch
		}
	}

	if len(batch) > 0 {
		err := c.DB.InsertDocument(batch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *Config) TestEndpoint(w http.ResponseWriter, r *http.Request) {

	//test GRPC connection
	//set request
	req := &pb.EmbeddingsMessageRequest{
		Text: "hello there",
	}
	resp, err := app.GRPCClient.TextToEmbedding(r.Context(), req)

	if err != nil {
		fmt.Println("unable to calle GRP", err)
		return
	}

	jsonResponse := make(map[string]string)
	jsonResponse["status"] = "ok"
	jsonResponse["message"] = resp.Text
	//set response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)

}

func (app *Config) TextToEmbedding(ctx context.Context, text string) ([]float64, error) {
	// set req for grpc service
	if len(text) == 0 {
		return nil, errors.New("text cant be empty")
	}
	req := &pb.EmbeddingsMessageRequest{
		Text: text,
	}
	resp, err := app.GRPCClient.TextToEmbedding(ctx, req)
	if err != nil {
		fmt.Println("unable to calle GRP", err)
		return nil, err
	}
	_ = resp
	return nil, nil

}
