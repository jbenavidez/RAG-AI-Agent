# RAG-AI-Agent(WIP)
A lightweight RAG-based AI agent that ingests content from the NYC Capital projects dataset and answers user questions with relevant, context-aware responses grounded in the ingested data.

### Stack
<ul>
<li>Go — High-performance backend powering core microservices </li>
<li>Python — processing and embedding generation</li>
<li>gRPC — Fast, type-safe communication between services</li>
<li>Postgres + pgvector — Vector-enabled database for similarity search </li>
</ul>

 ## micro-services descriptions
*   **Client-service**:  Service built on Go  responsible for handling Retrieval-Augmented Generation (RAG) operations. It coordinates requests, communicates with the embedding service, and queries the vector database
*   **Embending-service**: Service built on Python dedicated to generating vector embeddings. The model is easily replaceable—swap in any embedding provider such as OpenAI, HuggingFace, Nomic, etc., without changing the system architecture.
*   **Database**: Primary storage using PostgreSQL enhanced with the pgvector extension for similarity search and vector operations. This component is fully pluggable and can be replaced with modern vector databases such as Pinecone, Chroma, Weaviate, or Milvus depending on
   performance needs.

### Dataset
The agent’s knowledge comes from  data.cityofnewyork.us, with data currently imported via a CSV file. This ingestion method is fully flexible—CSV loading can easily be replaced with direct API-based retrieval, since data.cityofnewyork.us provides both CSV exports and a public API for accessing article content.

