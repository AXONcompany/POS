package table

import (
	"time"
)

type Table struct {
	ID          int64      `json:"id"`
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

type TableWaitress struct {
	ID         int64 `json:"id"`
	TableID    int64 `json:"table_id"`
	WaitressID int64 `json:"waitress_id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
