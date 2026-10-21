package ticket

import "time"

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