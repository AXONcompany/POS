package product

import "time"

type Category struct {
	ID        int64
	VenueID   int
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
