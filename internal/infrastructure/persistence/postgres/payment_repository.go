package postgres

import (
	"context"
	"fmt"

	"github.com/AXONcompany/POS/internal/domain/payment"
)

type PaymentRepository struct {
	db *DB
}

func NewPaymentRepository(db *DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, p *payment.Payment) (*payment.Payment, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO payments (order_id, division_id, payment_method, amount, tip, total, status, reference, venue_id, pos_terminal_id, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at`

	err = tx.QueryRow(ctx, query,
		p.OrderID, p.DivisionID, p.PaymentMethod, p.Amount, p.Tip, p.Total,
		p.Status, p.Reference, p.VenueID, p.POSTerminalID, p.UserID,
	).Scan(&p.ID, &p.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	if p.DivisionID != nil && *p.DivisionID != "" {
		_, err = tx.Exec(ctx, "UPDATE order_divisions SET is_paid = true WHERE id = $1", *p.DivisionID)
		if err != nil {
			return nil, fmt.Errorf("update division status: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return p, nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id int64) (*payment.Payment, error) {
	query := `
		SELECT id, order_id, division_id, payment_method, amount, tip, total, status, reference, venue_id, pos_terminal_id, user_id, created_at
		FROM payments WHERE id = $1`

	p := &payment.Payment{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.OrderID, &p.DivisionID, &p.PaymentMethod, &p.Amount, &p.Tip,
		&p.Total, &p.Status, &p.Reference, &p.VenueID, &p.POSTerminalID, &p.UserID, &p.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}

	return p, nil
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID int64) ([]*payment.Payment, error) {
	query := `
		SELECT id, order_id, division_id, payment_method, amount, tip, total, status, reference, venue_id, pos_terminal_id, user_id, created_at
		FROM payments WHERE order_id = $1 ORDER BY created_at`

	rows, err := r.db.Pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("get payments by order: %w", err)
	}
	defer rows.Close()

	payments := make([]*payment.Payment, 0)
	for rows.Next() {
		p := &payment.Payment{}
		if err := rows.Scan(
			&p.ID, &p.OrderID, &p.DivisionID, &p.PaymentMethod, &p.Amount, &p.Tip,
			&p.Total, &p.Status, &p.Reference, &p.VenueID, &p.POSTerminalID, &p.UserID, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan payment: %w", err)
		}
		payments = append(payments, p)
	}

	return payments, nil
}
