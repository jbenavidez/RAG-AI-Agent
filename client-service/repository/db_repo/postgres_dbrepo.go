package dbrepo

import (
	"client/models"
	"client/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

type Document struct {
	Text        string
	ProjectName string
	Description string
}

func (m *PostgresDBRepo) GetTotalDocuments() (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var totalDocs int
	query := `
		select
			 COUNT(*) AS total_docs
		from
			documents
	`
	row := m.DB.QueryRowContext(ctx, query)
	err := row.Scan(&totalDocs)
	if err != nil {
		return totalDocs, nil
	}
	return totalDocs, nil
}

func toPGVector(v []float64) string {
	parts := make([]string, len(v))
	for i, x := range v {
		parts[i] = fmt.Sprintf("%g", x)
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ","))
}

func (m *PostgresDBRepo) InsertDocument(documents []models.Document) error {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if len(documents) == 0 {
		return errors.New("no document provided")
	}
	valueStrings := []string{}
	valueArgs := []interface{}{}
	argPos := 1

	for _, r := range documents {
		embedding := utils.SimpleEmbedding(r.Text)

		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", argPos, argPos+1, argPos+2, argPos+3))
		valueArgs = append(valueArgs, r.Text, toPGVector(embedding), r.ProjectName, r.Description)
		argPos += 4
	}

	stmt := fmt.Sprintf(
		"INSERT INTO documents (text, vector, project_name, description) VALUES %s",
		strings.Join(valueStrings, ","),
	)

	_, err := m.DB.ExecContext(ctx, stmt, valueArgs...)

	if err != nil {
		return err
	}
	fmt.Println("valinor data inserted succefully")
	return nil
}

func (m *PostgresDBRepo) GetEmbeddingDocument(queryText string, topK int) ([]models.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	queryVector := utils.SimpleEmbedding(queryText)

	stmt := `
		SELECT text,  vector <-> $1 AS distance
		FROM documents
		ORDER BY vector <-> $1
		LIMIT $2
	`
	rows, err := m.DB.QueryContext(ctx, stmt, toPGVector(queryVector), topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var documents []models.Document
	for rows.Next() {
		var doc models.Document
		err := rows.Scan(
			&doc.Text,
			&doc.Distance,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}
	fmt.Println("getting the result", len((documents)))
	return documents, nil
}
