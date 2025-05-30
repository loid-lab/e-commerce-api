package models

type Address struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	FullName  string `gorm:"not null"`
	Street    string
	City      string
	Country   string
	ZipCode   string
	IsDefault bool
}
