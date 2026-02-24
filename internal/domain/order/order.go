package order

import "time"

type Order struct {
	ID           int       `json:"id" db:"id"`
	RestaurantID int       `json:"restaurant_id" db:"restaurant_id"`
	UserID       int       `json:"user_id" db:"user_id"` // Mesero who created it
	TableID      *int      `json:"table_id,omitempty" db:"table_id"`
	Status       string    `json:"status" db:"status"` // E.g., 'OPEN', 'PAID', 'CANCELLED'
	TotalAmount  float64   `json:"total_amount" db:"total_amount"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
