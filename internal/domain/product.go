package domain

import "time"

type Product struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Price        int        `json:"price"`
	Stock        int        `json:"stock"`
	CategoryID   int        `json:"category_id"`
	CategoryName string     `json:"category_name,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-"` // Hidden from JSON
}
