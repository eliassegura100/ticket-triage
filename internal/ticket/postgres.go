package ticket

import ( 
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) Create(ctx context.Context, t Ticket) (Ticket, error) { // Create method adds a new ticket to the Postgres database
	err := s.pool.QueryRow(ctx,
		`INSERT INTO tickets (title, description, reporter)
		 VALUES ($1, $2, $3)
		 RETURNING id, status, created_at`,
		t.Title, t.Description, t.Reporter,
	).Scan(&t.ID, &t.Status, &t.CreatedAt) // Scan the returned values into the Ticket struct
	if err != nil {
		return Ticket{}, fmt.Errorf("insert ticket: %w", err)
	}
	return t, nil
}

func (s *PostgresStore) Get(ctx context.Context, id int) (Ticket, error) { // Get method retrieves a ticket by its ID from the Postgres database
	var t Ticket
	err := s.pool.QueryRow(ctx,
		`SELECT id, title, description, reporter, status, created_at
		 FROM tickets WHERE id = $1`,
		id,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Reporter, &t.Status, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) { // If no rows are returned, return the ErrNotFound error
		return Ticket{}, ErrNotFound
	}
	if err != nil {
		return Ticket{}, fmt.Errorf("get ticket %d: %w", id, err)
	}
	return t, nil
}

func (s *PostgresStore) List(ctx context.Context) ([]Ticket, error) { // List method returns a slice of all tickets in the Postgres database
	rows, err := s.pool.Query(ctx,
		`SELECT id, title, description, reporter, status, created_at
		 FROM tickets ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list tickets: %w", err)
	}
	defer rows.Close()

	out := []Ticket{} // empty slice, so the API returns [] instead of null
	for rows.Next() {
		var t Ticket
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Reporter, &t.Status, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan ticket: %w", err)
		}
		out = append(out, t)
	}
	return out, rows.Err()
}