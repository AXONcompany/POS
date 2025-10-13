package table

import "time"

type Table struct{
    ID int64    `gorm:"primaryKey"`
    Number int
    Capacity int
    IsAvailable bool
    OcuppiedAt time.Time
    ReleasedAt time.Time
    OrderID int64
}