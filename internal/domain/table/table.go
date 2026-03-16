package table

import (
	"time"
)

const (
	StatusLibre     = "LIBRE"
	StatusOcupada   = "OCUPADA"
	StatusReservada = "RESERVADA"
)

// ValidStatus verifica que el estado sea uno de los permitidos.
func ValidStatus(s string) bool {
	switch s {
	case StatusLibre, StatusOcupada, StatusReservada:
		return true
	}
	return false
}

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

// Assignment representa una asignacion activa de mesero a mesa.
type Assignment struct {
	ID           int64      `json:"id"`
	TableID      int64      `json:"table_id"`
	UserID       int        `json:"user_id"`
	VenueID      int        `json:"venue_id"`
	AssignedAt   time.Time  `json:"assigned_at"`
	UnassignedAt *time.Time `json:"unassigned_at,omitempty"`
}

// AssignmentDetail incluye el nombre del mesero para las vistas de historial.
type AssignmentDetail struct {
	ID           int64      `json:"id"`
	TableID      int64      `json:"table_id"`
	UserID       int        `json:"user_id"`
	VenueID      int        `json:"venue_id"`
	WaiterName   string     `json:"waiter_name"`
	AssignedAt   time.Time  `json:"assigned_at"`
	UnassignedAt *time.Time `json:"unassigned_at,omitempty"`
}
