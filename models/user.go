package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`                            // hashed
	Role      string `gorm:"type:varchar(20);default:'customer'"` // "customer", "admin"
	CreatedAt time.Time
	UpdatedAt time.Time
	Orders    []Order
	Addresses []Address
}
