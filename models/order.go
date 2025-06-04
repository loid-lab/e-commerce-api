package models

import "time"

type Order struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`
	Status    string  `gorm:"type:varchar(20);default:'pending'"` // pending, paid, shipped
	Total     float64 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []OrderItem
	Payment   Payment
}
