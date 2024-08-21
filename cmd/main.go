package main

import (
	"context"
	"kara-bank/server"
	"kara-bank/util"
	"log"
	"os"
)

func main() {
	log.Println("Starting kara-bank")

	log.Println("Initializing token maker")
	pasetoMaker := util.NewPasetoMaker("") // TODO: get key for token generation

	log.Println("Connecting to database")
	connPool := server.ConnectToDb(context.Background())
	log.Println("Connected to databse")

	log.Println("Initializing http server")
	httpServer := server.InitHttpServer(connPool, pasetoMaker)

	log.Printf("Starting app on port %s", os.Getenv("SERVER_PORT"))
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
