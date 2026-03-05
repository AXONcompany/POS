package order

import (
	"context"
	"fmt"

	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
)

type Repository interface {
	Create(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error)
	GetByID(ctx context.Context, id int64, venueID int) (*domainOrder.Order, error)
	UpdateStatus(ctx context.Context, id int64, venueID int, statusID int) error
	ListByTable(ctx context.Context, tableID int64, venueID int) ([]domainOrder.Order, error)
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) CreateOrder(ctx context.Context, venueID, userID int, tableID *int64, items []domainOrder.OrderItem) (*domainOrder.Order, error) {
	o, err := domainOrder.NewOrder(venueID, userID, tableID, items)
	if err != nil {
		return nil, err
	}

	return uc.repo.Create(ctx, o)
}

// CreateOrderWithoutItems crea una orden sin items.
func (uc *Usecase) CreateOrderWithoutItems(ctx context.Context, venueID, userID int, tableID *int64) (*domainOrder.Order, error) {
	o := &domainOrder.Order{
		VenueID:     venueID,
		UserID:      userID,
		TableID:     tableID,
		StatusID:    1, // PENDING
		TotalAmount: 0,
	}

	return uc.repo.Create(ctx, o)
}

func (uc *Usecase) AddProductToOrder(ctx context.Context, venueID int, orderID int64, items []domainOrder.OrderItem) error {
	// TODO: Implementar logica completa de agregar items a orden existente
	return nil
}

func (uc *Usecase) GetOrderByID(ctx context.Context, venueID int, orderID int64) (*domainOrder.Order, error) {
	return uc.repo.GetByID(ctx, orderID, venueID)
}

func (uc *Usecase) CheckoutOrder(ctx context.Context, venueID int, orderID int64) error {
	// 5 = PAID
	return uc.repo.UpdateStatus(ctx, orderID, venueID, 5)
}

func (uc *Usecase) ListOrdersByTable(ctx context.Context, venueID int, tableID int64) ([]domainOrder.Order, error) {
	return uc.repo.ListByTable(ctx, tableID, venueID)
}

func (uc *Usecase) UpdateOrderStatus(ctx context.Context, venueID int, orderID int64, statusID int) error {
	return uc.repo.UpdateStatus(ctx, orderID, venueID, statusID)
}

// CancelOrderItem cancela un item de una orden.
func (uc *Usecase) CancelOrderItem(ctx context.Context, venueID int, orderID, itemID int64) error {
	// TODO: Implementar logica completa de cancelar item
	_ = fmt.Sprintf("cancelando item %d de orden %d", itemID, orderID)
	return nil
}
