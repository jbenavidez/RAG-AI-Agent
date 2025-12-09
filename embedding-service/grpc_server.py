# grpc_server.py

import grpc
from concurrent import futures
from fastembed import TextEmbedding

import embedding_pb2, embedding_pb2_grpc


# Load Model for embedding 
model = TextEmbedding("sentence-transformers/all-MiniLM-L6-v2")


class EmbeddingService(embedding_pb2_grpc.EmbeddingServiceServicer):

    def TextToEmbedding(self, request, context):
        text = request.text
        
        print("Valinor is calling", text)

        vectors = model.encode([text], normalize_embeddings=True)
        embeddings = vectors.tolist() if hasattr(vectors, "tolist") else vectors

        return embedding_pb2.EmbeddingsMessageResponse(
            text=str(embeddings)
        )


def start_grpc_server():
    """Start the gRPC server on port 50051."""
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    embedding_pb2_grpc.add_EmbeddingServiceServicer_to_server(
        EmbeddingService(), server
    )

    server.add_insecure_port("[::]:50051")
    server.start()
    print("gRPC server running on port 50051")

    server.wait_for_termination()
