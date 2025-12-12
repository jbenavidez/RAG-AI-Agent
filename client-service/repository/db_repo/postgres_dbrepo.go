package dbrepo

import (
	"client/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/pgvector/pgvector-go"
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

func toPGVector(v []float32) string {
	parts := make([]string, len(v))
	for i, x := range v {
		parts[i] = fmt.Sprintf("%f", x) // standard decimal notation
	}
	return "[" + strings.Join(parts, ",") + "]"
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
		if len(r.EmbeddingText) == 0 {
			return fmt.Errorf("document %q has empty embedding", r.Text)
		}

		// Normalize the embedding to unit length
		normalized := make([]float32, len(r.EmbeddingText))
		var norm float32
		for _, v := range r.EmbeddingText {
			norm += v * v
		}
		norm = float32(math.Sqrt(float64(norm)))
		for i, v := range r.EmbeddingText {
			normalized[i] = v / norm
		}

		// Wrap as pgvector.Vector
		pgVector := pgvector.NewVector(normalized)

		// Prepare SQL placeholders and values
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", argPos, argPos+1, argPos+2, argPos+3))
		valueArgs = append(valueArgs, r.Text, pgVector, r.ProjectName, r.Description)
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

	fmt.Println("Documents inserted successfully")
	return nil
}

func (m *PostgresDBRepo) GetEmbeddingDocument(queryVector []float32, topK int, keyword string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if len(queryVector) == 0 {
		return nil, fmt.Errorf("empty query vector")
	}

	// Normalize vector
	var sum float32
	for _, v := range queryVector {
		sum += v * v
	}
	norm := float32(math.Sqrt(float64(sum)))
	for i := range queryVector {
		queryVector[i] /= norm
	}

	vec := pgvector.NewVector(queryVector)

	stmt := `
        SELECT text
        FROM documents
        ORDER BY vector <#> $1::vector
        LIMIT $2
    `
	rows, err := m.DB.QueryContext(ctx, stmt, vec, topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []string{}
	existing := make(map[string]bool) //  avoid duplicates
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err != nil {
			return nil, err
		}
		results = append(results, text)
		existing[text] = true
	}

	// plan b if we dont get result
	if keyword != "" {
		kwStmt := `SELECT text FROM documents WHERE text ILIKE '%' || $1 || '%'`
		kwRows, _ := m.DB.QueryContext(ctx, kwStmt, keyword)
		if kwRows != nil {
			defer kwRows.Close()
			for kwRows.Next() {
				var text string
				_ = kwRows.Scan(&text)
				if !existing[text] {
					results = append(results, text)
				}
			}
		}
	}

	return results, nil
}
