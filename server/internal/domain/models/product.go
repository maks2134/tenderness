package models

import (
	"time"
)

type Product struct {
	ID          int       `db:"id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price" json:"price"`
	ImageURL    string    `db:"image_url" json:"image_url"`
	Category    string    `db:"category" json:"category"`
	InStock     bool      `db:"in_stock" json:"in_stock"`
	Rating      float64   `db:"rating" json:"rating"`
	Views       int       `db:"views" json:"views"`
}

type Category struct {
	ID          int       `db:"id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	ImageURL    string    `db:"image_url" json:"image_url"`
}
