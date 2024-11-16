package main

import (
	"context"
	dbserver "kara-bank/db"
	db "kara-bank/db/repositories"
	gapi "kara-bank/grpc_handler"
	"kara-bank/pb"
	"kara-bank/server"
	"kara-bank/services"
	"kara-bank/utils"
	"log"
	"net"
	"os"

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
	connPool := dbserver.ConnectToDb(context.Background())
	log.Println("Connected to databse")

	// init repository layer
	store := db.NewStore(connPool)

	// init service layer
	userService := services.NewUserService(store, pasetoMaker)
	accountService := services.NewAccountService(store)
	transferService := services.NewTransferService(store)

	//runRestServer(restPort, userService, accountService, transferService, pasetoMaker)
	runGrpcServer(grpcPort, userService, accountService, transferService)
}

func runGrpcServer(
	grpcPort string,
	userService services.UserServiceInterface,
	accountService services.AccountServiceInterface,
	transferService services.TransferServiceInterface,
) {
	log.Println("Initializing grpc server")
	handler := gapi.InitGrpcHandler(userService, accountService, transferService)

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

func runRestServer(
	port string,
	userService services.UserServiceInterface,
	accountService services.AccountServiceInterface,
	transferService services.TransferServiceInterface,
	tokenMaker utils.TokenMaker,
) {
	log.Println("Initializing rest server")
	httpServer := server.InitHttpServer(port, userService, accountService, transferService, tokenMaker)

	log.Printf("Starting app on port %s", port)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
