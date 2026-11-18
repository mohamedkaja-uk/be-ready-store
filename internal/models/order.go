package models

import "time"

type OrderItem struct {
BundleID string `json:"bundle_id"`
Quantity int    `json:"quantity"`
}

type CheckoutRequest struct {
CustomerName  string      `json:"customer_name" binding:"required"`
CustomerEmail string      `json:"customer_email" binding:"required"`
Items         []OrderItem `json:"items" binding:"required,min=1"`
}

type Order struct {
ID            string      `json:"id"`
CustomerName  string      `json:"customer_name"`
CustomerEmail string      `json:"customer_email"`
Items         []OrderItem `json:"items"`
TotalAmount   float64     `json:"total_amount"`
Status        string      `json:"status"`
CreatedAt     time.Time   `json:"created_at"`
}
