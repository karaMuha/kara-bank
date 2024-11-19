package rest

import (
	"context"
	db "kara-bank/db/repositories"
	"kara-bank/testcontainer"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var testStore db.Store

func TestMain(m *testing.M) {
	initScriptPath := filepath.Join("..", "testcontainer", "initScript.sql")
	testDb, err := testcontainer.CreatePostgresContainer(context.Background(), initScriptPath)

	if err != nil {
		log.Fatalf("Could not start test container: %v", err)
	}

	testStore = db.NewStore(testDb)

	os.Exit(m.Run())
}
