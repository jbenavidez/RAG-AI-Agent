FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential \
        git \
        curl \
        wget \
        ca-certificates \
        libffi-dev \
        libssl-dev \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Upgrade pip
RUN python -m pip install --upgrade pip --no-cache-dir --disable-pip-version-check

# Install Python dependencies
RUN python -m pip install --no-cache-dir \
        torch \
        transformers \
        sentence-transformers \
        grpcio \
        grpcio-tools \
        protobuf \
    --disable-pip-version-check

# Copy source code
COPY . .

# Set Python path
ENV PYTHONPATH=/app

# Expose gRPC port
EXPOSE 50001

# Start gRPC server
CMD ["python", "app.py"]
