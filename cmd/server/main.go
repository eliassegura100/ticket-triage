package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eliassegura100/ticket-triage/internal/api"
	"github.com/eliassegura100/ticket-triage/internal/ticket"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	store, closeStore, err := openStore(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer closeStore()

	handler := api.NewHandler(store)

	addr := ":8080"
	log.Printf("triage service listening on %s", addr)
	if err := http.ListenAndServe(addr, handler.Routes()); err != nil {
		log.Print(err)
	}
}

// openStore uses Postgres when DATABASE_URL is set, otherwise memory.
func openStore(ctx context.Context) (ticket.Store, func(), error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Print("DATABASE_URL not set; using in-memory store")
		return ticket.NewMemoryStore(), func() {}, nil
	}

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, nil, fmt.Errorf("create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("connect to postgres: %w", err)
	}

	log.Print("connected to PostgreSQL")
	return ticket.NewPostgresStore(pool), pool.Close, nil
}