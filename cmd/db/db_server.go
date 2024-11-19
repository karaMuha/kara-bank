package dbserver

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initDatabase(ctx context.Context, dbConnection string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbConnection)
	if err != nil {
		log.Printf("Error while connecting to database %v", err)
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Printf("Error while validating database connection: %v", err)
		return nil, err
	}

	return pool, nil
}

func ConnectToDb(ctx context.Context) *pgxpool.Pool {
	var count int
	dbConnection := os.Getenv("DB_CONNECTION")

	for {
		dbHandler, err := initDatabase(ctx, dbConnection)

		if err == nil {
			return dbHandler
		}

		log.Println("Postgres container not yet ready...")
		count++
		log.Println(count)

		if count > 10 {
			log.Fatalf("Failed to connect to database %v", err)
			return nil
		}

		log.Println("Backing off for five seconds...")
		time.Sleep(5 * time.Second)
	}
}
