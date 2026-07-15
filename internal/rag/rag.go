package rag

import (
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
