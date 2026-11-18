package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	// Fallback to a default local connection string if DATABASE_URL is not set
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
	dbURL = "postgres://postgres:postgres@localhost:5432/beready?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v\n", err)
	}

	// Disable named prepared statement caching to work smoothly with Neon's pooler
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	ctx := context.Background()
	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if err := DB.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	fmt.Println(" Connected to PostgreSQL database successfully!")
}