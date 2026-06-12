package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"tutorial.sqlc.dev/app/tutorial"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://tutorial:tutorial@localhost:5432/tutorial?sslmode=disable"
	}

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer conn.Close(ctx)

	queries := tutorial.New(conn)

	author, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  pgtype.Text{String: "Co-author of The C Programming Language", Valid: true},
	})
	if err != nil {
		log.Fatalf("create author: %v", err)
	}
	fmt.Printf("created author %d: %s\n", author.ID, author.Name)

	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		log.Fatalf("list authors: %v", err)
	}
	for _, a := range authors {
		fmt.Printf("- %d %s (%s)\n", a.ID, a.Name, a.Bio.String)
	}
}
