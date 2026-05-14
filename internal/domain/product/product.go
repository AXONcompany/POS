package product

import "time"

type Product struct {
	ID          int64
	VenueID     int
	Name        string
	Description string
	SalesPrice  float64
	ImageURL    string
	CategoryID  *int64
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}
