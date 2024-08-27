package db

import (
	"context"
	"kara-bank/testcontainer"
	"path/filepath"

	"log"
	"os"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	initScriptPath := filepath.Join("..", "..", "testcontainer", "initScript.sql")
	dbHandler, err := testcontainer.CreatePostgresContainer(context.Background(), initScriptPath)

	if err != nil {
		log.Fatalf("Could not start test container: %v", err)
	}

	testStore = NewStore(dbHandler)

	os.Exit(m.Run())
}
