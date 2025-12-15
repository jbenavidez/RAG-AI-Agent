# RAG-AI-Agent(WIP)
A lightweight **Retrieval-Augmented Generation (RAG)** AI agent that ingests data from the **NYC Capital Projects dataset** and answers user questions with relevant, context-aware responses grounded in the ingested data.The agent is designed to operate independently and can be integrated into other applications or services. A **web application** is provided as a testing and interaction tool, allowing users to easily query and validate the agent’s capabilities.


 ## Stack
- **Go** — Core language for building the **RAG AI agent**
- **Weaviate** — Vector database used by the **RAG AI agent** for storing and retrieving embeddings
- **gRPC** — High-performance communication between microservices
- **Ollama** — Local LLM runtime for serving LLaMA 2
- **LLaMA 2** — Large Language Model for context-aware response generation
- **OpenAI Embeddings** — Used for text vectorization via `text2vec-openai`
- **Docker & Docker Compose** — Containerization and service orchestration

 ## Features
- **Retrieval-Augmented Generation (RAG)** for accurate, data-grounded responses  
- **Semantic search** powered by vector embeddings  
- **Context-aware responses** generated using **LLaMA 2**  
- **Microservices-based architecture** for modularity and scalability  
- **Web application** for testing and interacting with the AI agent

 ## Microservices Details
The system is built using a **microservices architecture**:

- **RAG-IA-Agent-Service (Go)** — Handles queries, coordinates embedding retrieval, and orchestrates the AI agent  
- **Weaviate Database** — Stores and retrieves vector embeddings for semantic search  
- **Ollama (LLaMA 2)** — Generates natural language responses based on context  
- **Web-Application-Service** — Small Web-App for testing and interacting with the AI agent

 ## Data source
The agent’s knowledge comes from   **data.cityofnewyork.us**, with data currently imported via a CSV file. This ingestion method is fully flexible—CSV loading can easily be replaced with direct API-based retrieval, since  **data.cityofnewyork.us** provides both CSV exports and a public API for accessing the capital projects information.



 ## Setup and Running the Project
### 1. Start services

Start all services using Docker Compose:

```bash
make up_build
```

 ### 2. Accessing the Web Application
 Once all services are running, access the frontend in your browser:


```bash
http://localhost:4000\
```