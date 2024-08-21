package main

import (
	"context"
	"kara-bank/server"
	"kara-bank/utils"
	"log"
	"os"
)

func main() {
	log.Println("Starting kara-bank")
	port := os.Getenv("SERVER_PORT")

	log.Println("Initializing token maker")
	pasetoMaker := utils.NewPasetoMaker("") // TODO: get key for token generation

	log.Println("Connecting to database")
	connPool := server.ConnectToDb(context.Background())
	log.Println("Connected to databse")

	log.Println("Initializing http server")
	httpServer := server.InitHttpServer(port, connPool, pasetoMaker)

	log.Printf("Starting app on port %s", port)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
