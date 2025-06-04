package models

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	InvoiceNumber string        `json:"invoice_number"`
	Date          time.Time     `json:"date"`
	CustomerName  string        `json:"customer_name"`
	UserID        uint          `json:"user_id"`
	Items         []InvoiceItem `json:"items" gorm:"foreignKey:InvoiceID" `
	TotalAmount   float64
}

type InvoiceItem struct {
	gorm.Model
	InvoiceID   uint    `json:"invoice_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Amount      float64 `json:"amount"`
	TotalPrice  float64 `json:"total_price"`
}
