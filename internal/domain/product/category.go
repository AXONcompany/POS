package product

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	CategoryName string
	Products []*Product  `gorm:"foreignKey:CategoryID"`
}
