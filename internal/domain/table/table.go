package table

import (
	"time"
)

type Table struct {
	ID          int64      `json:"id"`
	VenueID     int        `json:"venue_id"`
	Number      int        `json:"table_number"`
	Capacity    int        `json:"capacity"`
	Status      string     `json:"status"`
	ArrivalTime *time.Time `json:"arrival_time"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type TableUpdates struct {
	Number      *int       `json:"table_number"`
	Capacity    *int       `json:"capacity"`
	Status      *string    `json:"status"`
	ArrivalTime *time.Time `json:"arrival_time"`
}
