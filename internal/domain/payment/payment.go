package payment

import "time"

// Payment representa un pago procesado.
type Payment struct {
	ID            int64     `db:"id"`
	OrderID       int64     `db:"order_id"`
	DivisionID    *string   `db:"division_id"`
	PaymentMethod string    `db:"payment_method"` // efectivo, tarjeta, multiple
	Amount        float64   `db:"amount"`
	Tip           float64   `db:"tip"`
	Total         float64   `db:"total"`
	Status        string    `db:"status"` // pendiente, aprobado, rechazado
	Reference     string    `db:"reference"`
	VenueID       int       `db:"venue_id"`
	POSTerminalID *int      `db:"pos_terminal_id"`
	UserID        int       `db:"user_id"`
	CreatedAt     time.Time `db:"created_at"`
}
