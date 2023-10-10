package model

import (
	"time"
)

type Product struct {
	ID          int        `db:"id"  json:"id"`
	Name        string     `db:"name" json:"name"`
	Price       int        `db:"price" json:"price"`
	Description string     `db:"description" json:"description"`
	Quantity    int        `db:"quantity" json:"quantity"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type ProductList struct {
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
	TotalPages int       `json:"total_pages"`
	Products   []Product `json:"products"`
}
