package ticket

import (
	"context"
	"errors"
	"sync"
	"testing"
)

// testStore runs the same checks against any Store implementation.
// newStore must return an empty store each time it is called.
func testStore(t *testing.T, newStore func(t *testing.T) Store) {
	ctx := context.Background()

	t.Run("create assigns ID, status, and time", func(t *testing.T) {
		s := newStore(t)
		got, err := s.Create(ctx, Ticket{Title: "Login fails"})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if got.ID != 1 {
			t.Errorf("ID = %d, want 1", got.ID)
		}
		if got.Status != StatusOpen {
			t.Errorf("Status = %q, want %q", got.Status, StatusOpen)
		}
		if got.CreatedAt.IsZero() {
			t.Error("CreatedAt was not set")
		}
	})

	t.Run("get existing and missing", func(t *testing.T) {
		s := newStore(t)
		created, err := s.Create(ctx, Ticket{Title: "Login fails"})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}

		got, err := s.Get(ctx, created.ID)
		if err != nil {
			t.Fatalf("Get(%d): %v", created.ID, err)
		}
		if got.Title != "Login fails" {
			t.Errorf("Title = %q, want %q", got.Title, "Login fails")
		}

		if _, err := s.Get(ctx, 99); !errors.Is(err, ErrNotFound) {
			t.Errorf("Get(99) error = %v, want ErrNotFound", err)
		}
	})

	t.Run("list returns tickets in order", func(t *testing.T) {
		s := newStore(t)
		got, err := s.List(ctx)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if got == nil || len(got) != 0 {
			t.Fatalf("empty List() = %v, want empty non-nil slice", got)
		}

		for _, title := range []string{"a", "b", "c"} {
			if _, err := s.Create(ctx, Ticket{Title: title}); err != nil {
				t.Fatalf("Create: %v", err)
			}
		}
		got, err = s.List(ctx)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(got) != 3 {
			t.Fatalf("List() has %d tickets, want 3", len(got))
		}
		for i, tk := range got {
			if tk.ID != i+1 {
				t.Errorf("List()[%d].ID = %d, want %d", i, tk.ID, i+1)
			}
		}
	})

	t.Run("concurrent creates are all saved", func(t *testing.T) {
		s := newStore(t)
		const n = 50

		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if _, err := s.Create(ctx, Ticket{Title: "concurrent"}); err != nil {
					t.Errorf("Create: %v", err)
				}
			}()
		}
		wg.Wait()

		got, err := s.List(ctx)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(got) != n {
			t.Errorf("List() has %d tickets, want %d", len(got), n)
		}
	})
}

func TestMemoryStore(t *testing.T) {
	testStore(t, func(t *testing.T) Store { return NewMemoryStore() })
}