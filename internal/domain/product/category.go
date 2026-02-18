package product

import "time"

type Category struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
