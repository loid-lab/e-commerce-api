package models

import "time"

type Cart struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []CartItem
}
