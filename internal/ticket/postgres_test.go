package ticket

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TestPostgresStore runs only when TEST_DATABASE_URL is set,
// so `go test ./...` still works without a database.
func TestPostgresStore(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping Postgres tests")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)

	testStore(t, func(t *testing.T) Store {
		// Wipe the table and reset IDs so every subtest starts clean.
		if _, err := pool.Exec(ctx, "TRUNCATE tickets RESTART IDENTITY"); err != nil {
			t.Fatalf("truncate: %v", err)
		}
		return NewPostgresStore(pool)
	})
}