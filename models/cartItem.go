package models

type CartItem struct {
	ID        uint `gorm:"primaryKey"`
	CartID    uint `gorm:"not null"`
	ProductID uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
	Product   Product
	UserID    uint
}
