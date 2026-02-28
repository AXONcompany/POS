package order_test

import (
	"context"
	"errors"
	"testing"

	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
	uc "github.com/AXONcompany/POS/internal/usecase/order"
)

type mockOrderRepository struct {
	createFunc       func(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error)
	getByIDFunc      func(ctx context.Context, id int64, restaurantID int) (*domainOrder.Order, error)
	updateStatusFunc func(ctx context.Context, id int64, restaurantID int, statusID int) error
	listByTableFunc  func(ctx context.Context, tableID int64, restaurantID int) ([]domainOrder.Order, error)
}

func (m *mockOrderRepository) Create(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, o)
	}
	return o, nil
}

func (m *mockOrderRepository) GetByID(ctx context.Context, id int64, restaurantID int) (*domainOrder.Order, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id, restaurantID)
	}
	return nil, nil
}

func (m *mockOrderRepository) UpdateStatus(ctx context.Context, id int64, restaurantID int, statusID int) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, id, restaurantID, statusID)
	}
	return nil
}

func (m *mockOrderRepository) ListByTable(ctx context.Context, tableID int64, restaurantID int) ([]domainOrder.Order, error) {
	if m.listByTableFunc != nil {
		return m.listByTableFunc(ctx, tableID, restaurantID)
	}
	return []domainOrder.Order{}, nil
}

func TestCreateOrder_Success(t *testing.T) {
	mockRepo := &mockOrderRepository{
		createFunc: func(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error) {
			o.ID = 1
			return o, nil
		},
	}
	usecase := uc.NewUsecase(mockRepo)

	tableID := int64(5)
	items := []domainOrder.OrderItem{
		{ProductID: 1, Quantity: 2, UnitPrice: 10.0},
		{ProductID: 2, Quantity: 1, UnitPrice: 5.0},
	}

	order, err := usecase.CreateOrder(context.Background(), 1, 2, &tableID, items)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.TotalAmount != 25.0 {
		t.Errorf("expected total amount 25.0, got %f", order.TotalAmount)
	}
	if order.StatusID != 1 {
		t.Errorf("expected status PENDING (1), got %d", order.StatusID)
	}
	if order.ID != 1 {
		t.Errorf("expected ID 1, got %d", order.ID)
	}
}

func TestCreateOrder_Error(t *testing.T) {
	mockRepo := &mockOrderRepository{
		createFunc: func(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error) {
			return nil, errors.New("db error")
		},
	}
	usecase := uc.NewUsecase(mockRepo)

	tableID := int64(5)
	_, err := usecase.CreateOrder(context.Background(), 1, 2, &tableID, nil)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestCheckoutOrder_Success(t *testing.T) {
	mockRepo := &mockOrderRepository{
		updateStatusFunc: func(ctx context.Context, id int64, restaurantID int, statusID int) error {
			if statusID != 5 {
				return errors.New("expected status PAID")
			}
			return nil
		},
	}
	usecase := uc.NewUsecase(mockRepo)

	err := usecase.CheckoutOrder(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateOrderStatus_Success(t *testing.T) {
	mockRepo := &mockOrderRepository{
		updateStatusFunc: func(ctx context.Context, id int64, restaurantID int, statusID int) error {
			if statusID != 3 {
				return errors.New("expected status READY")
			}
			return nil
		},
	}
	usecase := uc.NewUsecase(mockRepo)

	err := usecase.UpdateOrderStatus(context.Background(), 1, 100, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestListOrdersByTable_Success(t *testing.T) {
	mockRepo := &mockOrderRepository{
		listByTableFunc: func(ctx context.Context, tableID int64, restaurantID int) ([]domainOrder.Order, error) {
			if tableID != 5 {
				return nil, errors.New("wrong table ID")
			}
			return []domainOrder.Order{
				{ID: 1, TotalAmount: 100.0},
				{ID: 2, TotalAmount: 50.0},
			}, nil
		},
	}
	usecase := uc.NewUsecase(mockRepo)

	orders, err := usecase.ListOrdersByTable(context.Background(), 1, 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}
	if orders[0].TotalAmount != 100.0 {
		t.Errorf("expected first order amount 100, got %f", orders[0].TotalAmount)
	}
}
