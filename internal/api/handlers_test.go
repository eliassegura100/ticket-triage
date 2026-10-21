package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eliassegura100/ticket-triage/internal/ticket"
)

// newTestServer gives each test a fresh store and router.
func newTestServer() http.Handler {
	return NewHandler(ticket.NewStore()).Routes()
}

// do sends a request through the router and returns the recorded response.
func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCreateTicket(t *testing.T) {
	h := newTestServer()

	rec := do(t, h, http.MethodPost, "/tickets",
		`{"title":"Login fails","description":"500 on login","reporter":"support"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var got ticket.Ticket
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if got.ID != 1 || got.Title != "Login fails" || got.Status != ticket.StatusOpen {
		t.Errorf("got %+v, want ID 1, title %q, status open", got, "Login fails")
	}
}

func TestCreateTicketValidation(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"malformed JSON", `{"title":`},
		{"missing title", `{"description":"no title"}`},
		{"blank title", `{"title":"   "}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := do(t, newTestServer(), http.MethodPost, "/tickets", tc.body)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestGetTicket(t *testing.T) {
	h := newTestServer()
	do(t, h, http.MethodPost, "/tickets", `{"title":"Login fails"}`)

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"existing ticket", "/tickets/1", http.StatusOK},
		{"missing ticket", "/tickets/99", http.StatusNotFound},
		{"non-numeric id", "/tickets/abc", http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := do(t, h, http.MethodGet, tc.path, "")
			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}

func TestListTickets(t *testing.T) {
	h := newTestServer()
	do(t, h, http.MethodPost, "/tickets", `{"title":"first"}`)
	do(t, h, http.MethodPost, "/tickets", `{"title":"second"}`)

	rec := do(t, h, http.MethodGet, "/tickets", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got []ticket.Ticket
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("got %d tickets, want 2", len(got))
	}
}