package product

import "gorm.io/gorm"


type Product struct{
	gorm.Model
    ID int64    `gorm:"primaryKey"`
    Name string
    Price float64
    Notes string   
    CategoryID int 
} 