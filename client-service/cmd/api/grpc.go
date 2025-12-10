package main

import (
	pb "client/proto/generated"
	"context"
)

var app *Config

// NewGrpcHelper make app config avaible
func NewGrpcHelper(a *Config) {
	app = a
}

type server struct {
	pb.UnimplementedEmbeddingServiceServer
}

func (s *server) TextToEmbedding(ctx context.Context, request *pb.EmbeddingsMessageRequest) (*pb.EmbeddingsMessageResponse, error) {

	return &pb.EmbeddingsMessageResponse{Text: "Hello from Valinor"}, nil
}
