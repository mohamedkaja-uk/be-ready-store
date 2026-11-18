package models

import "time"

type CheckoutRequest struct {
	CustomerName    string     `json:"customer_name" binding:"required"`
	CustomerEmail   string     `json:"customer_email" binding:"required"`
	ShippingAddress string     `json:"shipping_address"`
	Items           []CartItem `json:"items" binding:"required"`
}

type CartItem struct {
	BundleID string `json:"bundle_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type Order struct {
	ID              string     `json:"id"`
	CustomerName    string     `json:"customer_name"`
	CustomerEmail   string     `json:"customer_email"`
	ShippingAddress string     `json:"shipping_address"`
	Items           []CartItem `json:"items"`
	TotalAmount     float64    `json:"total_amount"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
}