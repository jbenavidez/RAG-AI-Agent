FROM python:3.11-slim


WORKDIR /app

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

RUN python -m pip install --upgrade pip --no-cache-dir --disable-pip-version-check

RUN python -m pip install --no-cache-dir \
        fastembed==0.6.0 \
        flask \
        grpcio \
        grpcio-tools \
        protobuf \
    --disable-pip-version-check

COPY . .

ENV PYTHONPATH=/app


EXPOSE 8000 50051 50052

CMD ["python", "app.py"]
