package models

import "time"

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null"`
	ImageURL    string
	CategoryID  uint
	Category    Category
	CreatedBy   uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
