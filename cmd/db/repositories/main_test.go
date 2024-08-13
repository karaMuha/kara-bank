package db

import (
	"context"
	testContainer "kara-bank/db/testcontainer"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testStore Store

func TestMain(m *testing.M) {
	testDB, err := testContainer.CreatePostgresContainer(context.Background())

	if err != nil {
		log.Fatalf("Could not start test container: %v", err)
	}

	testQueries = New(testDB)
	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
