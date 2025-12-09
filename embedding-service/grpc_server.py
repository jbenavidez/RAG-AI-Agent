import grpc
from concurrent import futures
from sentence_transformers import SentenceTransformer

from proto.generated import embedding_pb2, embedding_pb2_grpc

# Load the model once
model = SentenceTransformer("all-MiniLM-L6-v2")


class EmbeddingService(embedding_pb2_grpc.EmbeddingServiceServicer):

    def TextToEmbedding(self, request, context):
        
        text = request.text
        print("the text is", text)
        embeddings = model.encode([text], normalize_embeddings=True).tolist()

        response_text = str(embeddings)

        return embedding_pb2.EmbeddingsMessageResponse(text=response_text)


class GRPCServer:

    @classmethod
    def start_server(cls):
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        embedding_pb2_grpc.add_EmbeddingServiceServicer_to_server(
            EmbeddingService(), server
        )

        # Listen on port 50051
        server.add_insecure_port("[::]:50051")
        server.start()
        print("gRPC server running on port 50051")

        # Keep server alive
        server.wait_for_termination()


 