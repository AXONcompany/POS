package order

import (
	"github.com/AXONcompany/POS/internal/domain/product"
	"github.com/AXONcompany/POS/internal/domain/table"
	"time"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID        int64	`gorm:"primaryKey"`
	OrderDate time.Time
	Total     float64
	Client    string
	Products  []*product.Product `gorm:"many2many: order_products;"`
	Tables 	  []*table.Table  `gorm:"foreignKey:OrderID;references:ID"`
}
