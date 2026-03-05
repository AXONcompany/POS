package user

import "time"

type User struct {
	ID           int        `json:"id" db:"id"`
	VenueID      int        `json:"venue_id" db:"venue_id"`
	RoleID       int        `json:"role_id" db:"role_id"`
	Name         string     `json:"name" db:"name"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"` // Never expose in JSON
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	Phone        *string    `json:"phone" db:"phone"`
	LastAccess   *time.Time `json:"last_access" db:"last_access"`
}
