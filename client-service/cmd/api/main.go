package main

import (
	pb "client/proto/generated"
	"client/repository"
	dbrepo "client/repository/db_repo"
	"fmt"
	"log"
	"net/http"
)

const (
	gRpcPort = "50001"
	port     = 4000
)

type Config struct {
	DB         repository.DatabaseRepo
	GRPCClient pb.EmbeddingServiceClient
}

func main() {
	fmt.Println("starting url-shortener service...")
	app := Config{}

	conn := app.connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	fmt.Println("client up and running")

	err := app.LoadData()
	if err != nil {
		fmt.Println("sotmhthing break", err)
	}
	log.Println("Starting agent on port", port)

	_, err = app.DB.GetEmbeddingDocument("valinor", 3)

	if err != nil {
		fmt.Println("valinor_faild", err)
		panic(err)
	}

	//
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
