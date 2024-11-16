package main

import (
	"context"
	"kara-bank/pb"
	"kara-bank/server"
	"kara-bank/utils"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting kara-bank")
	//restPort := os.Getenv("SERVER_PORT")
	grpcPort := os.Getenv("GRPC_SERVER_PORT")

	log.Println("Initializing token maker")
	pasetoMaker := utils.NewPasetoMaker("") // TODO: get key for token generation

	log.Println("Connecting to database")
	connPool := server.ConnectToDb(context.Background())
	log.Println("Connected to databse")

	// runRestServer(restPort, connPool, pasetoMaker)
	runGrpcServer(grpcPort, connPool, pasetoMaker)
}

func runGrpcServer(grpcPort string, connPool *pgxpool.Pool, tokenMaker utils.TokenMaker) {
	log.Println("Initializing grpc server")
	handler := server.InitGrpcHandler(connPool, tokenMaker)

	server := grpc.NewServer()
	pb.RegisterKaraBankServer(server, handler)
	reflection.Register(server)
	port := "0.0.0.0" + grpcPort

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	err = server.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server")
	}

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
