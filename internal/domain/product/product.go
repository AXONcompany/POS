package product

import "time"

type Product struct {
	ID         int64
	Name       string
	SalesPrice float64
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
}
