import grpc
from concurrent import futures
from sentence_transformers import SentenceTransformer
import numpy as np

import embedding_pb2, embedding_pb2_grpc

# Load the sentence-transformers model
model = SentenceTransformer("sentence-transformers/all-mpnet-base-v2")


class EmbeddingService(embedding_pb2_grpc.EmbeddingServiceServicer):
    """ gRPC service for generating text embeddings. """

    def TextToEmbedding(self, request, context):
        text = request.text

        # Generate embeddings
        embeddings = model.encode(
            [text],
            convert_to_numpy=True,
            normalize_embeddings=False
        )[0]


        embeddings = embeddings.astype(np.float32)


        norm = np.linalg.norm(embeddings)
        if norm > 0:
            embeddings /= norm

        embedding_list = embeddings.tolist()

        return embedding_pb2.EmbeddingsMessageResponse(
            result=embedding_list
        )

def start_grpc_server():
    """ Starts the gRPC server on port 50001. """
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    embedding_pb2_grpc.add_EmbeddingServiceServicer_to_server(
        EmbeddingService(), server
    )

    server.add_insecure_port("[::]:50001")
    server.start()
    print("GRPC server running on port 50001")

    server.wait_for_termination()


if __name__ == "__main__":
    start_grpc_server()
