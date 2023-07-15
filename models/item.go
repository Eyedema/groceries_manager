package models

import "time"

// Item represents an item with an ID, name, and creation timestamp.
type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
