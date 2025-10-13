package product

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID           int	`gorm:"primaryKey"`
	CategoryName string
	Products []*Product  `gorm:"foreignKey:CategoryID"`
}
