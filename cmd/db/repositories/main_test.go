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
	testContainer, err := testContainer.CreatePostgresContainer(context.Background())
	testDB := testContainer.Pool

	if err != nil {
		log.Fatalf("Could not start test container: %v", err)
	}

	testQueries = New(testDB)
	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
