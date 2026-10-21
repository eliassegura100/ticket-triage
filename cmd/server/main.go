package main

import (
	"log"
	"net/http"

	"github.com/eliassegura100/ticket-triage/internal/api"
	"github.com/eliassegura100/ticket-triage/internal/ticket"
)

func main() {
	store := ticket.NewStore()
	handler := api.NewHandler(store)

	addr := ":8080"
	log.Printf("triage service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler.Routes()))
}