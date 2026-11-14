package ticket

import (
	"sync"
	"time"
	"context"
)

type MemoryStore struct { // MemoryStore struct manages the collection of tickets
	mu	  sync.RWMutex // mutex to ensure thread-safe access to the tickets map
	nextID int
	tickets map[int]Ticket // map to store tickets with their ID as the key
}

func NewMemoryStore() *MemoryStore { // NewMemoryStore initializes a new MemoryStore instance
	return &MemoryStore{ nextID: 1, tickets: make(map[int]Ticket) }
}

// Create method adds a new ticket to the store
func (s *MemoryStore) Create(ctx context.Context, t Ticket) (Ticket, error) {
	s.mu.Lock() // lock the mutex to ensure exclusive access to the tickets map
	defer s.mu.Unlock() // defer unlock to ensure the mutex is released after the function completes

	t.ID = s.nextID
	s.nextID++
	t.Status = StatusOpen
	t.CreatedAt = time.Now()
	s.tickets[t.ID] = t
	return t, nil
}

func (s *MemoryStore) Get(ctx context.Context, id int) (Ticket, error) { // Get method retrieves a ticket by its ID from the store
	s.mu.Lock() // lock the mutex to ensure exclusive access to the tickets map
	defer s.mu.Unlock() // defer unlock to ensure the mutex is released after the function completes

	t, ok := s.tickets[id]
	if !ok {
		return Ticket{}, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) List(ctx context.Context) ([]Ticket, error) { // List method returns a slice of all tickets in the store
	s.mu.Lock() // lock the mutex to ensure exclusive access to the tickets map
	defer s.mu.Unlock() // defer unlock to ensure the mutex is released after the function completes

	out := make([]Ticket, 0, len(s.tickets))
	for id := 1; id < s.nextID; id++ { // iterate through all possible ticket IDs
		if t, ok := s.tickets[id]; ok { // check if the ticket exists in the map
			out = append(out, t) // append the ticket to the output slice
		}
	}
	return out, nil
}