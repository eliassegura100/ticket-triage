package ticket

import (
	"sync"
	"time"
)

type Store struct { // Store struct manages the collection of tickets
	mu	  sync.RWMutex // mutex to ensure thread-safe access to the tickets map
	nextID int
	tickets map[int]Ticket // map to store tickets with their ID as the key
}

func NewStore() *Store { // NewStore initializes a new Store instance
	return &Store{ nextID: 1, tickets: make(map[int]Ticket) }
}

func (s *Store) Create(t Ticket) Ticket { // Create method adds a new ticket to the store
	s.mu.Lock() // lock the mutex to ensure exclusive access to the tickets map
	defer s.mu.Unlock() // defer unlock to ensure the mutex is released after the function completes

	t.ID = s.nextID
	s.nextID++
	t.Status = StatusOpen
	t.CreatedAt = time.Now()
	s.tickets[t.ID] = t
	return t
}

func (s *Store) Get(id int) (Ticket, bool) { // Get method retrieves a ticket by its ID
	s.mu.RLock() // read lock to allow concurrent reads
	defer s.mu.RUnlock() // defer unlock to ensure the mutex is released after the function completes

	t, ok := s.tickets[id]
	return t, ok
}

func (s *Store) List() []Ticket { // List method returns a slice of all tickets in the store
	s.mu.Lock() // lock the mutex to ensure exclusive access to the tickets map
	defer s.mu.Unlock() // defer unlock to ensure the mutex is released after the function completes

	out := make([]Ticket, 0, len(s.tickets))
	for id := 1; id < s.nextID; id++ { // iterate through all possible ticket IDs
		if t, ok := s.tickets[id]; ok { // check if the ticket exists in the map
			out = append(out, t) // append the ticket to the output slice
		}
	}
	return out
}