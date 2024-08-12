package db

import (
	"context"
	testContainer "kara-bank/db/testcontainer"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	testContainer, err := testContainer.CreatePostgresContainer(context.Background())
	testDB = testContainer.Pool

	if err != nil {
		log.Fatalf("Could not start test container: %v", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
