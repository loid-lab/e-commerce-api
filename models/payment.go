package models

import "time"

type Payment struct {
	ID        uint   `gorm:"primaryKey"`
	OrderID   uint   `gorm:"not null"`
	Method    string // e.g. "card"
	Provider  string // e.g. "Stripe"
	RefID     string // Stripe PaymentIntent or CheckoutSession ID
	Status    string // "pending", "paid", "failed"
	PaidAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
