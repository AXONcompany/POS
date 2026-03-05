package table

import (
	"time"
)

// Request para crear (swagger: numero y capacidad)
type CreateRequest struct {
	Number   int    `json:"numero" binding:"required,gt=0"`
	Capacity int    `json:"capacidad" binding:"required,gt=0"`
	Status   string `json:"status"`
}

// Request para actualizar
type UpdateRequest struct {
	Number      *int       `json:"table_number" binding:"omitempty,gt=0"`
	Capacity    *int       `json:"capacity" binding:"omitempty,gt=0"`
	Status      *string    `json:"status"`
	ArrivalTime *time.Time `json:"arrival_time"`
}

// Request para actualizar solo el estado (swagger: PATCH /mesas/:id/estado)
type UpdateEstadoRequest struct {
	Estado string `json:"estado" binding:"required"`
}

// Request para asignar mesero
type AssignWaitressRequest struct {
	WaitressID int64 `json:"waitress_id" binding:"required"`
}

// Response
type Response struct {
	ID          int64      `json:"id"`
	Number      int        `json:"number"`
	Capacity    int        `json:"capacity"`
	Status      string     `json:"state"`
	ArrivalTime *time.Time `json:"arrival_time,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
