# RAG-AI-Agent

A lightweight RAG AI agent that ingests data from the **NYC Capital Projects dataset** and answers user questions with relevant, context-aware responses grounded in the ingested data.

The project includes a web application for uploading source files, indexing documents, and chatting with the AI agent through a browser-based interface.

The agent uses Redis as a memory system to remember recent chat history during an active session. This allows it to understand follow-up questions, keep the conversation going, and provide answers with better context from the previous messages.

## Stack

* **Go** — Core backend language for the web application, upload flow, WebSocket chat, memory handling, and RAG orchestration
* **Weaviate** — Vector database used for storing and retrieving NYC Capital Project documents
* **OpenAI Embeddings** — Used by Weaviate to vectorize project documents and user questions for semantic search
* **Redis** — Stores recent chat history by session so the agent can understand follow-up questions
* **Ollama** — Local LLM runtime used to run the language model
* **Llama 3** — Local language model used to generate answers from retrieved project context and chat history
* **Docker Compose** — Runs the Go app, Weaviate, Redis, and Ollama services locally
* **HTML, Bootstrap, and JavaScript** — APP for uploading files and interacting with the AI agent

## Features

* **File upload flow** for ingesting NYC Capital Projects CSV data
* **Document processing pipeline** that reads uploaded CSV rows and stores them in Weaviate
* **Retrieval-Augmented Generation (RAG)** for grounded answers using retrieved project data
* **Semantic search** powered by OpenAI embeddings through Weaviate
* **Redis-backed memory system** for storing recent chat history by session
* **Context-aware follow-up questions** using previous questions and answers from the active chat session
* **Local LLM response generation** using Ollama and Llama 3
* **WebSocket-based chat UI** for real-time question and answer interaction
* **Docker-based local development setup**

## Service Details

The project currently runs as a local Docker Compose setup:

* **Go RAG Web App** — Handles routing, uploads, WebSocket chat, document retrieval, memory retrieval, prompt building, and LLM orchestration
* **Weaviate Database** — Stores uploaded project documents and performs semantic retrieval
* **OpenAI Embeddings** — Used by Weaviate for vector search
* **Redis Memory Store** — Stores recent chat turns for each active session so the agent can answer follow-up questions with conversation context
* **Ollama / Llama 3** — Generates natural language answers using the retrieved project context and recent chat history

## Memory System

The agent uses Redis to store recent chat history for each WebSocket session.

Each successful chat turn stores:

* The user question
* The assistant answer
* The created timestamp

This allows the agent to handle follow-up questions such as:

```txt
Does that project have any total budget changes?
```

after a previous question like:

```txt
List the project names along with their budget forecast and total budget changes for each interceptor sewer project.
```

The memory is session-based and temporary, so it helps the agent maintain short-term conversation context without permanently storing every interaction.

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

REDIS_URL=redis://redis-service:6379
```


### 2. Start the Docker services

From the project folder, run:

```bash
make up_build
```

This will build and start the local services:

* web application
* Weaviate vector database
* Redis memory store
* Ollama local LLM service

### 3. Pull the Ollama model

If the model is not already available inside the Ollama container, pull it manually:

```bash
docker exec -it ollama-service ollama pull llama3
```

You can confirm the model is installed with:

```bash
docker exec -it ollama-service ollama list
```

### 4. Confirm Redis is running

You can confirm Redis is available with:

```bash
docker exec -it redis-service redis-cli ping
```

Expected response:

```txt
PONG
```

### 5. Access the web application

Once the services are running, open the app in your browser:

```bash
http://localhost:8080
```

From the web app, you can:

* Upload NYC Capital Projects CSV files
* View uploaded documents
* Ask questions about the ingested project data
* Ask follow-up questions using the active chat memory
* Receive AI-generated answers through the chat interface

### 6. Upload project data

Go to the upload page:

```bash
http://localhost:8080/upload
```

Upload the NYC Capital Projects CSV file. The app will:

* Save the uploaded file locally
* Store file metadata in Weaviate
* Process the CSV rows
* Insert project documents into Weaviate for retrieval

### 7. View uploaded documents

Go to:

```bash
http://localhost:8080/docs
```

This page displays uploaded files and their processing status.

### 8. Ask questions

Go back to the home page:

```bash
http://localhost:8080
```

Example question:

```txt
List the project names along with their budget forecast and total budget changes for each interceptor sewer project.
```

The app will retrieve relevant project records from Weaviate and generate an answer using Ollama and Llama 3.

You can then ask a follow-up question like:

```txt
Does that project have any total budget changes?
```

The app will use the Redis-backed chat history to understand that “that project” refers to the project from the previous answer.

### 9. Stop the services

To stop the project:

```bash
make down
```

### 10. Reset local Weaviate data

For local development, if you need to clear the Weaviate database and re-upload data:

```bash
make down
mv weaviate_data weaviate_data_backup_$(date +%Y%m%d_%H%M%S)
mkdir weaviate_data
make up_build
```

After resetting Weaviate, upload the CSV again from the web app.

## Example Questions

### Question #1

**Question:**

List the project names along with their budget forecast and total budget changes for each interceptor sewer project.

**Answer:**

Here are the projects that match your request:

* **REPAIRS & REHAB OF INTERCEPTING SEWERS, MN & BX** — **$30,611,000**
  Total budget changes are not available in the retrieved project data.

Note: There is only one interceptor sewer project in the provided data. The other retrieved projects are not interceptor sewer projects.

 <img width="773" height="602" alt="Screenshot 2026-07-22 at 8 53 54 PM" src="https://github.com/user-attachments/assets/8bbb0f9e-9e7e-499d-ba1f-be4918710f9a" />


---

### Follow-up Question

**Question:**

Does that project have any total budget changes?

**Answer:**

No, the available project data does not include total budget changes for **REPAIRS & REHAB OF INTERCEPTING SEWERS, MN & BX**. The available budget forecast is **$30,611,000**.

This shows that the agent can use the active chat history to understand that “that project” refers to the interceptor sewer project from the previous question.

<img width="702" height="281" alt="Screenshot 2026-07-22 at 8 54 57 PM" src="https://github.com/user-attachments/assets/ad841121-b1c6-4b85-bb5b-a4a2dcefa39a" />


---

### Question #2

**Question:**

What are the main objectives and expected benefits of the emergency technology infrastructure upgrade to improve public safety response times?

**Answer:**

The main objective of the emergency technology infrastructure upgrade is to enhance public safety response times by providing faster and more reliable communication systems, improved data analytics, and enhanced situational awareness. This helps emergency responders respond to incidents more quickly and effectively.

<img width="646" height="752" alt="Screenshot 2026-07-19 at 1 31 35 PM" src="https://github.com/user-attachments/assets/4ce11266-fa0a-41d9-a221-d1139995326c" />


