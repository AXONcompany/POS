package order

import (
	"context"

	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
)

type Repository interface {
	Create(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error)
	GetByID(ctx context.Context, id int64, restaurantID int) (*domainOrder.Order, error)
	UpdateStatus(ctx context.Context, id int64, restaurantID int, statusID int) error
	ListByTable(ctx context.Context, tableID int64, restaurantID int) ([]domainOrder.Order, error)
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) CreateOrder(ctx context.Context, restaurantID, userID int, tableID *int64, items []domainOrder.OrderItem) (*domainOrder.Order, error) {
	o, err := domainOrder.NewOrder(restaurantID, userID, tableID, items)
	if err != nil {
		return nil, err
	}

	return uc.repo.Create(ctx, o)
}

func (uc *Usecase) AddProductToOrder(ctx context.Context, restaurantID int, orderID int64, items []domainOrder.OrderItem) error {
	// Not fully implemented yet, but keeping structure for future tasks
	return nil
}

func (uc *Usecase) GetOrderByID(ctx context.Context, restaurantID int, orderID int64) (*domainOrder.Order, error) {
	return uc.repo.GetByID(ctx, orderID, restaurantID)
}

func (uc *Usecase) CheckoutOrder(ctx context.Context, restaurantID int, orderID int64) error {
	// Add business logic for payments, totals etc
	// 5 = PAID
	return uc.repo.UpdateStatus(ctx, orderID, restaurantID, 5)
}

func (uc *Usecase) ListOrdersByTable(ctx context.Context, restaurantID int, tableID int64) ([]domainOrder.Order, error) {
	return uc.repo.ListByTable(ctx, tableID, restaurantID)
}

func (uc *Usecase) UpdateOrderStatus(ctx context.Context, restaurantID int, orderID int64, statusID int) error {
	return uc.repo.UpdateStatus(ctx, orderID, restaurantID, statusID)
}
