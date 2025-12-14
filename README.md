# RAG-AI-Agent(WIP)
A lightweight RAG-based AI agent that ingests content from the NYC Capital projects dataset and answers user questions with relevant, context-aware responses grounded in the ingested data.

 ## Stack
*   **Go** 
*   **Weaviate** 
 

 ## Micro-Services descriptions
*   **RAG-service**: A service built in Go responsible for handling Retrieval-Augmented Generation (RAG) operations. It manages incoming requests, communicates with the embedding service, and performs queries to the vector database. The service relies on Weaviate to retrieve relevant embeddings for generating context-aware responses.
*   **Database**: Is used as the vector database for storing and retrieving embeddings efficiently. It leverages the text2vec-openai vectorizer to transform textual data into high-dimensional vectors, enabling semantic search and similarity-based retrieval.
*   **Ollama**: A service that uses LLaMA 2 to generate responses based on the context provided by the RAG-service.
*   **CLI-service**: A command-line interface for interacting with the AI agent. Users can ask questions about NYC capital projects, and the agent will respond using the structured project data.

 ## Data source
The agent’s knowledge comes from  data.cityofnewyork.us, with data currently imported via a CSV file. This ingestion method is fully flexible—CSV loading can easily be replaced with direct API-based retrieval, since data.cityofnewyork.us provides both CSV exports and a public API for accessing capital projects content.

 ## Example
  ##### The following record is stored in our Weaviate DB
| Project Name                                     | Description                                                                                     |
|-------------------------------------------------|-------------------------------------------------------------------------------------------------|
| REPAIRS & REHAB OF INTERCEPTING SEWERS, MN & BX | Repair and Rehabilitation of Interceptor Sewers at Various Locations in the Boroughs of Manhattan and The Bronx |

##### Example Generated Response

**Question:** List all the boroughs mentioned for interceptor sewer projects.  

**Answer:**  
1. Manhattan  
2. The Bronx