package model

import "time"

type Post_new struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
