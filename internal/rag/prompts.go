package rag

import (
	"github.com/tmc/langchaingo/prompts"
)

func BuildRAGPrompt(projectData []byte, question string) (string, error) {
	systemTemplateStr := `
You are an AI assistant for NYC Capital Projects.

Use only the project data provided below to answer the user's question.

Rules:
- Do not make up project names, budget values, dates, agencies, or status information.
- If the answer is not available in the project data, say that the available project data does not contain enough information.
- Keep the answer concise and easy to read.
- When listing projects, use bullet points or a numbered list.
- If the user asks for specific fields, only include those fields.
- Do not mention JSON, embeddings, vector search, Weaviate, CSV files, or internal system details.
- Do not say "based on the provided data" in every sentence.
- Preserv

Project Data:
{{.projectData}}
`

	systemTemplate := prompts.NewSystemMessagePromptTemplate(
		systemTemplateStr,
		[]string{"projectData"},
	)

	humanTemplate := prompts.NewHumanMessagePromptTemplate(
		"Question: {{.question}}",
		[]string{"question"},
	)

	chatTemplate := prompts.NewChatPromptTemplate(
		[]prompts.MessageFormatter{
			systemTemplate,
			humanTemplate,
		},
	)

	data := map[string]any{
		"projectData": string(projectData),
		"question":    question,
	}

	formattedChatPrompt, err := chatTemplate.Format(data)
	if err != nil {
		return "", err
	}

	return formattedChatPrompt, nil
}
