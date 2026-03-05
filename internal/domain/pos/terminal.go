package pos

import "time"

// Terminal representa un terminal POS (caja registradora) dentro de una sede.
type Terminal struct {
	ID           int       `json:"id" db:"id"`
	VenueID      int       `json:"venue_id" db:"venue_id"`
	TerminalName string    `json:"terminal_name" db:"terminal_name"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
