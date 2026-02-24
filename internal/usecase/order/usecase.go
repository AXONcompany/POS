package order

import (
	"context"

	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
)

type Repository interface {
	Create(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error)
	GetByID(ctx context.Context, id int, restaurantID int) (*domainOrder.Order, error)
	UpdateStatus(ctx context.Context, id int, restaurantID int, status string) error
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) CreateOrder(ctx context.Context, restaurantID, userID int, tableID *int) (*domainOrder.Order, error) {
	o := &domainOrder.Order{
		RestaurantID: restaurantID,
		UserID:       userID,
		TableID:      tableID,
		Status:       "OPEN",
		TotalAmount:  0,
	}

	return uc.repo.Create(ctx, o)
}

func (uc *Usecase) CheckoutOrder(ctx context.Context, restaurantID, orderID int) error {
	// Add business logic for payments, totals etc
	return uc.repo.UpdateStatus(ctx, orderID, restaurantID, "PAID")
}
