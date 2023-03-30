package models

import "time"

type Sale struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	SalePrice float64   `json:"sale_price"`
	SaleDate  time.Time `json:"sale_date"`
}
