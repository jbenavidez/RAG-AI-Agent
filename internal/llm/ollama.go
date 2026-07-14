package llm

import (
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms/ollama"
)

func NewOllamaClient() (*ollama.LLM, error) {
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "llama3"
	}

	serverURL := os.Getenv("OLLAMA_URL")
	if serverURL == "" {
		serverURL = "http://localhost:11434"
	}

	llm, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(serverURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ollama: %w", err)
	}

	fmt.Println("Ollama connected")

	return llm, nil
}
