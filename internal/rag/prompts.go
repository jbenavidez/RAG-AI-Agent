package rag

import (
	"github.com/tmc/langchaingo/prompts"
)

func BuildRAGPrompt(projectData []byte, question string) (string, error) {
	systemTemplateStr := `
Answer the user's question concisely, professionally, and in a friendly manner using the project information below.

Rules:
- Use only the project information provided.
- Do not mention the data source.
- Do not mention JSON, CSV, slides, or internal storage.
- Number items if there are multiple.

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
