import threading
from flask import Flask, request, jsonify

from fastembed import TextEmbedding
from grpc_server import start_grpc_server 

# Load FastEmbed 
model = TextEmbedding("sentence-transformers/all-MiniLM-L6-v2")

app = Flask(__name__)


@app.route("/", methods=["GET"])
def root():
    return jsonify({
        "status": "running",
        "model": "Welcome to gondor3",
        "grpc_port": 50051
    })


@app.route("/embed", methods=["POST"])
def embed():
    data = request.get_json()
    if not data or "text" not in data:
        return jsonify({"error": "Missing 'text' field"}), 400

    text = data["text"]
    normalize = data.get("normalize", True)
    inputs = [text] if isinstance(text, str) else text

    vectors = model.encode(inputs, normalize_embeddings=normalize)
    embeddings = vectors.tolist() if hasattr(vectors, "tolist") else vectors

    return jsonify({
        "model": "sentence-transformerads2",
        "count": len(embeddings),
        "embeddings": embeddings
    })


def start_services():
    # Start gRPC in background
    grpc_thread = threading.Thread(target=start_grpc_server, daemon=True)
    grpc_thread.start()

    # Start Flask (main thread)
    print("Flask running on port 8000")
    app.run(host="0.0.0.0", port=8000)


if __name__ == "__main__":
    start_services()
