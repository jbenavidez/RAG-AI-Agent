package main

import (
	"fmt"
	"log"
	"ragService/repository"
	dbrepo "ragService/repository/db_repo"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

const (
	gRpcPort = "50001"
)

type RagConfig struct {
	Llm            llms.Model
	WeaviateClient *weaviate.Client
	WDBRepo        repository.DatabaseRepo
}

func main() {
	app := RagConfig{}
	fmt.Println("************* connecting to Weaviate *************")
	client, err := app.ConnectWeaviateDB()
	if err != nil {
		fmt.Printf("unable to connected %v", err)
		panic(err)
	}
	app.WDBRepo = &dbrepo.WeaviateDBRepo{DB: client}
	app.WeaviateClient = client
	fmt.Println("************* Loading Data *************")
	err = app.LoadDataSet()
	if err != nil {
		fmt.Println("somethings break", err)
	}
	fmt.Println("*************  Init Ollama *************")
	// Initialize Ollama LLMs
	llm, err := ollama.New(
		ollama.WithModel("llama2"),
		ollama.WithServerURL("http://ollama-service:11434"),
	)
	if err != nil {
		fmt.Println("failed to Initialize Ollama: ", err)
		panic(err)
	}
	fmt.Println("*************  Ollama Connected *************")
	app.Llm = llm
	log.Println("Starting GRPC server on port", gRpcPort)
	app.gRPCListenAndServe()
}
