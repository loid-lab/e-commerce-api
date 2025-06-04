package models

type OrderItem struct {
	ID         uint    `gorm:"primaryKey"`
	OrderID    uint    `gorm:"not null"`
	ProductID  uint    `gorm:"not null"`
	Quantity   int     `gorm:"not null"`
	UnitPrice  float64 `gorm:"not null"` // snapshot of product price
	TotalPrice float64 `gorm:"not null"`
	Product    Product
}
