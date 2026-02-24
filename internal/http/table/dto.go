package table

import (
	"time"
)

// Request para crear
type CreateRequest struct {
	Number   int    `json:"table_number" binding:"required,gt=0"`
	Capacity int    `json:"capacity" binding:"required,gt=0"`
	Status   string `json:"status" binding:"required"`
}

// Request para actualizar
type UpdateRequest struct {
	Number      *int       `json:"table_number" binding:"omitempty,gt=0"`
	Capacity    *int       `json:"capacity" binding:"omitempty,gt=0"`
	Status      *string    `json:"status"`
	ArrivalTime *time.Time `json:"arrival_time"`
}

// Request para asignar mesero
type AssignWaitressRequest struct {
	WaitressID int64 `json:"waitress_id" binding:"required"`
}

// Response
type Response struct {
	ID          int64      `json:"id"`
	Number      int        `json:"table_number"`
	Capacity    int        `json:"capacity"`
	Status      string     `json:"status"`
	ArrivalTime *time.Time `json:"arrival_time"`
	CreatedAt   time.Time  `json:"created_at"`
}
