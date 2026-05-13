package table

import (
	"time"
)

type Table struct {
	ID      int64  `json:"id"`
	VenueID int    `json:"venue_id"`
	Number  int    `json:"table_number"`
	Name    string `json:"name"`

	Capacity         int        `json:"capacity"`
	Status           string     `json:"status"`
	Guests           int        `json:"guests"`
	AssignedWaiterID *int       `json:"assigned_waiter_id,omitempty"`
	ArrivalTime      *time.Time `json:"arrival_time,omitempty"`

	// Canvas layout
	X          int     `json:"x"`
	Y          int     `json:"y"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Shape      string  `json:"shape"`
	Rotation   int     `json:"rotation"`
	Color      *string `json:"color,omitempty"`
	Floor      int     `json:"floor"`
	IsMerged   bool    `json:"is_merged"`
	MergedFrom []int64 `json:"merged_from,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type TableUpdates struct {
	Number      *int       `json:"table_number"`
	Capacity    *int       `json:"capacity"`
	Status      *string    `json:"status"`
	ArrivalTime *time.Time `json:"arrival_time"`
}
