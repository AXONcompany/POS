package postgres

import (
	"context"

	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
)

type OrderRepository struct {
	// mock implementation just to satisfy compiler
}

func NewOrderRepository(db *DB) *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) Create(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error) {
	return o, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int, restaurantID int) (*domainOrder.Order, error) {
	return nil, nil // mock
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int, restaurantID int, status string) error {
	return nil
}
