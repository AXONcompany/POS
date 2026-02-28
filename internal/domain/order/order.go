package order

import "time"

type OrderStatus struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type OrderItem struct {
	ID        int64     `json:"id" db:"id"`
	OrderID   int64     `json:"order_id" db:"order_id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	UnitPrice float64   `json:"unit_price" db:"unit_price"`
	Notes     string    `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Order struct {
	ID           int64      `json:"id" db:"id"`
	RestaurantID int        `json:"restaurant_id" db:"restaurant_id"`
	TableID      *int64     `json:"table_id,omitempty" db:"table_id"`
	UserID       int        `json:"user_id" db:"user_id"` // Mesero who created it
	StatusID     int        `json:"status_id" db:"status_id"`
	TotalAmount  float64    `json:"total_amount" db:"total_amount"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

	Items  []OrderItem `json:"items,omitempty" db:"-"`
	Status string      `json:"status,omitempty" db:"-"` // Joined from order_statuses
}
