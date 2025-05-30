package models

type Category struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	Slug      string `gorm:"uniqueIndex;not null"`
	CreatedBy uint
}
