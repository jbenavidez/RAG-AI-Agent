package models

type Document struct {
	Text          string
	EmbeddingText []float32
	ProjectName   string
	Description   string
	Distance      float64
}
