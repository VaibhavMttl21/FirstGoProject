package model

import "time"

type Product struct {
    ID            string    `json:"id" gorm:"type:uuid;primaryKey"`
    Name          string    `json:"name" gorm:"not null"`
    Description   string    `json:"description" gorm:"not null;default:''"`
    PriceCents    int64     `json:"price_cents" gorm:"not null"`
    Currency      string    `json:"currency" gorm:"type:varchar(3);not null"`
    StockQuantity int       `json:"stock_quantity" gorm:"not null"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}