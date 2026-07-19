package rag

import (
	"context"

	"github.com/tmc/langchaingo/llms"
)

type Rag struct {
	LLM llms.Model
}

func New(llm llms.Model) *Rag {

	return &Rag{
		LLM: llm,
	}
}

func (r *Rag) Generate(ctx context.Context, prompt string) (string, error) {
	answer, err := llms.GenerateFromSinglePrompt(ctx, r.LLM, prompt)
	if err != nil {
		return "", err
	}

	return answer, nil
}
