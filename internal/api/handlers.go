package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/eliassegura100/ticket-triage/internal/ticket"
)

type Handler struct { // Handler struct manages HTTP requests related to tickets
	store *ticket.Store
}

func NewHandler(store *ticket.Store) *Handler { // NewHandler initializes a new Handler instance with the provided ticket store
	return &Handler{store: store}
}

// Routes wires URL patterns to handler methods.
func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tickets", h.createTicket)
	mux.HandleFunc("GET /tickets", h.listTickets)
	mux.HandleFunc("GET /tickets/{id}", h.getTicket)
	return mux
}

type createTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Reporter    string `json:"reporter"`
}

func (h *Handler) createTicket(w http.ResponseWriter, r *http.Request) { // createTicket handles the creation of a new ticket
	var req createTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // decode the JSON request body into the createTicketRequest struct
		writeError(w, http.StatusBadRequest, "invalid JSON body") // return a 400 Bad Request error if the JSON is invalid
		return
	}
	if strings.TrimSpace(req.Title) == "" { // check if the title is empty or only whitespace
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	t := h.store.Create(ticket.Ticket{ // create a new Ticket struct with the provided title, description, and reporter
		Title:       req.Title,
		Description: req.Description,
		Reporter:    req.Reporter,
	})
	writeJSON(w, http.StatusCreated, t)
}

func (h *Handler) listTickets(w http.ResponseWriter, r *http.Request) { // listTickets handles the retrieval of all tickets
	writeJSON(w, http.StatusOK, h.store.List())
}

func (h *Handler) getTicket(w http.ResponseWriter, r *http.Request) { // getTicket handles the retrieval of a specific ticket by its ID
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be a number")
		return
	}

	t, ok := h.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "ticket not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func writeJSON(w http.ResponseWriter, status int, v any) { // writeJSON writes a JSON response with the given status code and value
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) { // writeError writes a JSON error response with the given status code and error message
	writeJSON(w, status, map[string]string{"error": msg})
}