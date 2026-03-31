package order

import (
	"errors"
	"time"
)

var (
	ErrInvalidOrderItems       = errors.New("invalid order items: must have at least one item")
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
	ErrInsufficientStock       = errors.New("insufficient stock for ingredient")
)

type RecipeLine struct {
	IngredientID     int64
	QuantityRequired float64
}

type StockDeduction struct {
	IngredientID int64
	VenueID      int
	Quantity     float64
}

// Estado IDs: 1=PENDING, 2=SENT, 3=PREPARING, 4=READY, 5=PAID, 6=CANCELLED
//
// Se permite saltar estados intermedios para reducir fricción operativa (MVP sin KDS).
// Con KDS, las transiciones SENT→PREPARING→READY serán disparadas automáticamente.
var validTransitions = map[int][]int{
	1: {2, 3, 4, 5, 6}, // PENDING → cualquier estado adelante
	2: {3, 4, 5, 6},    // SENT    → cualquier estado adelante
	3: {4, 5, 6},       // PREPARING → READY | PAID | CANCELLED
	4: {5, 6},          // READY → PAID | CANCELLED
	5: {},              // PAID (terminal)
	6: {},              // CANCELLED (terminal)
}

func CanTransitionTo(current, next int) bool {
	allowed, ok := validTransitions[current]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}

func NewOrder(venueID, userID int, tableID *int64, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrInvalidOrderItems
	}

	o := &Order{
		VenueID:     venueID,
		UserID:      userID,
		TableID:     tableID,
		StatusID:    1, // 1 = PENDING Assuming this is the default
		TotalAmount: 0,
		Items:       items,
	}

	for _, item := range items {
		o.TotalAmount += item.UnitPrice * float64(item.Quantity)
	}

	return o, nil
}

type OrderStatus struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

var ErrItemAlreadyCancelled = errors.New("order item is already cancelled")

type OrderItem struct {
	ID          int64      `json:"id" db:"id"`
	OrderID     int64      `json:"order_id" db:"order_id"`
	ProductID   int64      `json:"product_id" db:"product_id"`
	Quantity    int        `json:"quantity" db:"quantity"`
	UnitPrice   float64    `json:"unit_price" db:"unit_price"`
	Notes       string     `json:"notes" db:"notes"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
}

type Order struct {
	ID            int64      `json:"id" db:"id"`
	VenueID       int        `json:"venue_id" db:"venue_id"`
	TableID       *int64     `json:"table_id,omitempty" db:"table_id"`
	UserID        int        `json:"user_id" db:"user_id"` // Mesero who created it
	POSTerminalID *int       `json:"pos_terminal_id,omitempty" db:"pos_terminal_id"`
	StatusID      int        `json:"status_id" db:"status_id"`
	TotalAmount   float64    `json:"total_amount" db:"total_amount"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

	Items  []OrderItem `json:"items,omitempty" db:"-"`
	Status string      `json:"status,omitempty" db:"-"` // Joined from order_statuses
}
