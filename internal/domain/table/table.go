package table

import (
	"time"

	"gorm.io/gorm"
)

type Table struct{
    gorm.Model
    Number int
    Capacity int
    IsAvailable bool
    OcuppiedAt time.Time
    ReleasedAt time.Time
    OrderID uint
}