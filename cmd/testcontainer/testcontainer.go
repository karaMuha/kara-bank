package testcontainer

import (
	"context"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreatePostgresContainer(ctx context.Context) (*pgxpool.Pool, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16.3-alpine",
		postgres.WithInitScripts(filepath.Join("..", "migration", "000001_init_schema.up.sql")),
		postgres.WithDatabase("kara_bank_db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(10*time.Second),
		),
	)

	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		return nil, err
	}

	dbHandler, err := pgxpool.New(ctx, connStr)

	if err != nil {
		return nil, err
	}

	err = dbHandler.Ping(ctx)

	if err != nil {
		return nil, err
	}

	return dbHandler, nil
}
