package main

import (
	"context"
	"kara-bank/pb"
	"kara-bank/server"
	"kara-bank/utils"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting kara-bank")
	port := os.Getenv("SERVER_PORT")

	log.Println("Initializing token maker")
	pasetoMaker := utils.NewPasetoMaker("") // TODO: get key for token generation

	log.Println("Connecting to database")
	connPool := server.ConnectToDb(context.Background())
	log.Println("Connected to databse")

	runRestServer(port, connPool, pasetoMaker)
}

func runGrpcServer(connPool *pgxpool.Pool, tokenMaker utils.TokenMaker) {
	log.Println("Initializing grpc server")
	handler := server.InitGrpcHandler(connPool, tokenMaker)

	grpc := grpc.NewServer()
	pb.RegisterKaraBankServer(grpc, handler)
}

func runRestServer(port string, connPool *pgxpool.Pool, tokenMaker utils.TokenMaker) {
	log.Println("Initializing rest server")
	httpServer := server.InitHttpServer(port, connPool, tokenMaker)

	log.Printf("Starting app on port %s", port)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
