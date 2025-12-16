package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "ragAIAgent/proto/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedAIAgentServiceServer
}

func (app *RagConfig) gRPCListenAndServe() {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	pb.RegisterAIAgentServiceServer(srv, &server{})

	reflection.Register(srv)
	log.Printf("gRPC server connected on port %s ", gRpcPort)
	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}

func (s *server) GetAIAgentAnswerFromUserQuestion(ctx context.Context, request *pb.AIAgentRequest) (*pb.AIAgentResponse, error) {
	question := request.GetQuestion()

	if len(question) == 0 {
		return &pb.AIAgentResponse{Answer: "plase ask question related to NYC capital project"}, nil
	}
	// get data from db
	result, err := RAGConfig.WDBRepo.GetDocuments(question)
	if err != nil {
		fmt.Println("unable to get data", err)
		return nil, err
	}
	// generate response
	resp, err := RAGConfig.GenerateAnswerFromSlides(ctx, question, result)
	if err != nil {
		fmt.Println("unable to get data", err)
		return nil, err
	}
	//return rsponse
	return &pb.AIAgentResponse{Answer: resp}, nil
}
