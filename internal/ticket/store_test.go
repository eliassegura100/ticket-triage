package ticket

import (
	"sync"
	"testing"
)

func TestCreateAssignsIDStatusAndTime(t *testing.T) { // Test that Create assigns an ID, sets the status to open, and sets CreatedAt
	s := NewStore()

	got := s.Create(Ticket{Title: "Login fails"})

	if got.ID != 1 {
		t.Errorf("ID = %d; want 1", got.ID)
	}
	if got.Status != StatusOpen {
		t.Errorf("Status = %q, want %q", got.Status, StatusOpen)
	}
	if got.CreatedAt.IsZero() {
		t.Errorf("CreatedAt was not set")
	}
}

func TestCreateIncrementIDs(t *testing.T) { // Test that IDs are incremented correctly
	s := NewStore()
	first := s.Create(Ticket{Title: "first"})
	second := s.Create(Ticket{Title: "second"})

	if first.ID != 1 || second.ID != 2 {
		t.Errorf("IDs = %d, %d, want 1, 2", first.ID, second.ID)
	}
}

func TestGet(t *testing.T) {
	s := NewStore()
	created := s.Create(Ticket{Title: "Login fails"})
	got, ok := s.Get(created.ID)
	if !ok {
		t.Fatalf("Get(%d) found nothing, want a ticket", created.ID)
	}
	if got.Title != "Login fails" {
		t.Errorf("Title = %q, want %q", got.Title, "Login fails")
	}

	if _, ok := s.Get(99); ok {
		t.Error("Get(99) found a ticket, want none")
	}
}
func TestListReturnsTicketsInOrder(t *testing.T) {
	s := NewStore()
	if got := s.List(); len(got) != 0 {
		t.Fatalf("empty store List() has %d tickets, want 0", len(got))
	}

	s.Create(Ticket{Title: "a"})
	s.Create(Ticket{Title: "b"})
	s.Create(Ticket{Title: "c"})

	got := s.List()
	if len(got) != 3 {
		t.Fatalf("List() has %d tickets, want 3", len(got))
	}
	for i, tk := range got {
		if tk.ID != i+1 {
			t.Errorf("List()[%d].ID = %d, want %d", i, tk.ID, i+1)
		}
	}
}

func TestCreateIsSafeForConcurrentUse(t *testing.T) { 
	s := NewStore()
	const n = 100 // Test that Create is safe for concurrent use by creating 100 tickets concurrently

	var wg sync.WaitGroup // Use a WaitGroup to wait for all goroutines to finish
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done() // Decrement the WaitGroup counter when the goroutine completes
			s.Create(Ticket{Title: "concurrent"})
		}()
	}
	wg.Wait()

	if got := len(s.List()); got != n { 
		t.Errorf("List() has %d tickets, want %d", got, n)
	}
}