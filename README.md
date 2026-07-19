# RAG-AI-Agent

A lightweight RAG AI agent that ingests data from the **NYC Capital Projects dataset** and answers user questions with relevant, context-aware responses grounded in the ingested data.

The project includes a Go web application for uploading source files, indexing documents, and chatting with the AI agent through a browser-based interface.

## Stack

- **Go** — Core backend language for the web application, upload flow, WebSocket chat, and RAG orchestration
- **Chi Router** — HTTP routing for pages, upload endpoints, static files, and WebSocket routes
- **Gorilla WebSocket** — Real-time chat communication between the browser and the Go backend
- **Weaviate** — Vector database used for storing and retrieving NYC Capital Project documents
- **OpenAI Embeddings** — Used by Weaviate to vectorize project documents and user questions for semantic search
- **Ollama** — Local LLM runtime used to run the language model
- **Llama 3** — Local language model used to generate answers from retrieved project context
- **Docker Compose** — Runs the Go app, Weaviate, and Ollama services locally
- **HTML, Bootstrap, and JavaScript** — Browser UI for uploading files and interacting with the AI agent

## Features

- **File upload flow** for ingesting NYC Capital Projects CSV data
- **Document processing pipeline** that reads uploaded CSV rows and stores them in Weaviate
- **Retrieval-Augmented Generation (RAG)** for grounded answers using retrieved project data
- **Semantic search** powered by OpenAI embeddings through Weaviate
- **Local LLM response generation** using Ollama and Llama 3
- **WebSocket-based chat UI** for real-time question and answer interaction
- **Docker-based local development setup**

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

Create a `.env` file in the project root:

```env
APP_PORT=8080
WEAVIATE_PORT=8081
APP_CMD=./cmd/web

OPENAI_APIKEY=your-openai-api-key-here
OPENAI_EMBEDDING_MODEL=text-embedding-3-small

OLLAMA_URL=http://ollama-service:11434
OLLAMA_MODEL=llama3
CHUNK_SIZE=1000
```



### 2. Start the Docker services

From the project folder, run:

```bash
make up_build
```

This will build and start the local services:

- Go web application
- Weaviate vector database
- Ollama local LLM service

### 3. Pull the Ollama model

If the model is not already available inside the Ollama container, pull it manually:

```bash
docker exec -it ollama-service ollama pull llama3
```

You can confirm the model is installed with:

```bash
docker exec -it ollama-service ollama list
```

### 4. Access the web application

Once the services are running, open the app in your browser:

```bash
http://localhost:8080
```

From the web app, you can:

- Upload NYC Capital Projects CSV files
- View uploaded documents
- Ask questions about the ingested project data
- Receive AI-generated answers through the chat interface

### 5. Upload project data

Go to the upload page:

```bash
http://localhost:8080/upload
```

Upload the NYC Capital Projects CSV file. The app will:

- Save the uploaded file locally
- Store file metadata in Weaviate
- Process the CSV rows
- Insert project documents into Weaviate for retrieval

### 6. View uploaded documents

Go to:

```bash
http://localhost:8080/docs
```

This page displays uploaded files and their processing status.

### 7. Ask questions

Go back to the home page:

```bash
http://localhost:8080
```

Example question:

```txt
List the project names along with their budget forecast and total budget changes for each interceptor sewer project.
```

The app will retrieve relevant project records from Weaviate and generate an answer using Ollama and Llama 3.

### 8. Stop the services

To stop the project:

```bash
make down
```

### 9. Reset local Weaviate data

For local development, if you need to clear the Weaviate database and re-upload data:

```bash
make down
mv weaviate_data weaviate_data_backup_$(date +%Y%m%d_%H%M%S)
mkdir weaviate_data
make up_build
```

After resetting Weaviate, upload the CSV again from the web app.

## Example Questions 

**Question # 1:**

List the project names along with their budget forecast and total budget changes for each interceptor sewer project.

**Answer:**

The AI agent retrieves relevant NYC Capital Project records from Weaviate and generates a concise answer using the project names, budget forecasts, and total budget changes found in the retrieved data.

<img width="720" height="806" alt="Screenshot 2026-07-19 at 1 27 27 PM" src="https://github.com/user-attachments/assets/9046b324-055d-439b-9530-17ee139ff240" />

---

### Question #2

**Question:**

What are the main objectives and expected benefits of the emergency technology infrastructure upgrade to improve public safety response times?

**Answer:**

The main objective of the emergency technology infrastructure upgrade is to enhance public safety response times by providing faster and more reliable communication systems, improved data analytics, and enhanced situational awareness. This helps emergency responders respond to incidents more quickly and effectively.
 
<img width="646" height="752" alt="Screenshot 2026-07-19 at 1 31 35 PM" src="https://github.com/user-attachments/assets/4ce11266-fa0a-41d9-a221-d1139995326c" />