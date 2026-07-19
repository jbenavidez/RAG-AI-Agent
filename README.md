# RAG-AI-Agent (WIP)

A lightweight RAG AI agent that ingests data from the **NYC Capital Projects dataset** and answers user questions with relevant, context-aware responses grounded in the ingested data.

The project includes a Go web application for uploading source files, indexing documents, and chatting with the AI agent through a browser-based interface.

## Stack

- **Go** — Core backend language for the web application, upload flow, WebSocket chat, and RAG orchestration
- **Weaviate** — Vector database used for storing and retrieving NYC Capital Project documents
- **OpenAI Embeddings** — Used by Weaviate to vectorize project documents and user questions for semantic search
- **Ollama** — Local LLM runtime used to run the language model
- **Llama 3** — Local language model used to generate answers from retrieved project context
- **Docker Compose** — Runs the Go app, Weaviate, and Ollama services locally

## Features

- **File upload flow** for ingesting NYC Capital Projects CSV data
- **Document processing pipeline** that reads uploaded CSV rows and stores them in Weaviate
- **RAG** for grounded answers using retrieved project data
- **Semantic search** powered by OpenAI embeddings through Weaviate
- **Local LLM response generation** using Ollama and Llama 3
- **WebSocket-based chat UI** for real-time question and answer interaction

## Service Details

The project currently runs as a local Docker Compose setup:

- **Go RAG Web App** — Handles routing, uploads, WebSocket chat, document retrieval, prompt building, and LLM orchestration
- **Weaviate Database** — Stores uploaded project documents and performs semantic retrieval
- **OpenAI Embeddings** — Used by Weaviate for vector search
- **Ollama / Llama 3** — Generates natural language answers using the retrieved project context

## Data Source

The agent’s knowledge comes from **data.cityofnewyork.us**, with data currently imported through a CSV file.

The ingestion flow is flexible and can later be extended to support direct API-based retrieval from the NYC Open Data API.

## Setup and Running the Project

### 1. Add environment variables

Create a `.env` file with the required local settings:

```env
APP_PORT=8080
WEAVIATE_PORT=8081
APP_CMD=./cmd/web

OPENAI_APIKEY=your-openai-api-key-here
OPENAI_EMBEDDING_MODEL=text-embedding-3-small

OLLAMA_URL=http://ollama-service:11434
OLLAMA_MODEL=llama3
CHUNK_SIZE=1000