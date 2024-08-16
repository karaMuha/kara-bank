package db

import (
	"context"
	"kara-bank/testcontainer"

	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testStore Store

func TestMain(m *testing.M) {
	testDB, err := testcontainer.CreatePostgresContainer(context.Background())

	if err != nil {
		log.Fatalf("Could not start test container: %v", err)
	}

	testQueries = New(testDB)
	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
