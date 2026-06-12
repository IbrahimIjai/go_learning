package config

import (
	db "authentication/internal/db"
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool is the shared pgx connection pool.
// Queries is the sqlc-generated query interface bound to that pool.
var (
	Pool    *pgxpool.Pool
	Queries *db.Queries
)

// ConnectDB establishes the Postgres connection pool and initializes the sqlc
// Queries. It reads DATABASE_URL, falling back to the local docker-compose DSN.
func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/usersdb?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to create Postgres pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Postgres is not reachable: %v", err)
	}

	Pool = pool
	Queries = db.New(pool)
	log.Println("Successfully connected to Postgres!")
}
