package db

import (
	"context"
	db "kara-bank/db/testcontainer"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	testContainer, err := db.CreatePostgresContainer(context.Background())

	if err != nil {
		log.Fatalf("Could not start test container: %v", err)
	}

	testQueries = New(testContainer.Pool)

	os.Exit(m.Run())
}
