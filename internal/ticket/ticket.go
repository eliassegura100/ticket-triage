package ticket

import (
	"context"
	"errors"
	"time"
)

type Status string

const ( // three states for a ticket
	StatusOpen   Status = "open"
	StatusTriaged Status = "triaged"
	StatusResolved Status = "resolved"
)

type Ticket struct { // Ticket struct represents a ticket in the system
	ID 		int    `json:"id"` 		// struct tags for JSON serialization
	Title 	string `json:"title"`
	Description string `json:"description"`
	Reporter string `json:"reporter"`
	Status Status `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var ErrNotFound = errors.New("ticket not found")

type Store interface { // Store is anything that can save and load tickets. It is implemented by MemoryStore and PostgresStore.
	Create(ctx context.Context, t Ticket) (Ticket, error)
	Get(ctx context.Context, id int) (Ticket, error)
	List(ctx context.Context) ([]Ticket, error)
}