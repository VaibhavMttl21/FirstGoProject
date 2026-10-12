package model

import "time"

type Product struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PriceCents    int64     `json:"price_cents"`
	Currency      string    `json:"currency"`
	StockQuantity int       `json:"stock_quantiy"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
