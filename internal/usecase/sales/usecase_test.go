package sale_test

import (
	"context"
	"errors"
	"testing"

	domainSale "github.com/AXONcompany/POS/internal/domain/sale"
	uc "github.com/AXONcompany/POS/internal/usecase/sales"
)

// Mock SaleRepository
type mockSaleRepository struct {
	createSaleFunc func(ctx context.Context, s domainSale.Sale) (*domainSale.Sale, error)
	getByIDFunc    func(ctx context.Context, id int64) (*domainSale.Sale, error)
}

func (m *mockSaleRepository) CreateSale(ctx context.Context, s domainSale.Sale) (*domainSale.Sale, error) {
	if m.createSaleFunc != nil {
		return m.createSaleFunc(ctx, s)
	}
	return &s, nil
}

func (m *mockSaleRepository) GetByID(ctx context.Context, id int64) (*domainSale.Sale, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

// Mock OrderUsecase
type mockOrderUsecase struct {
	checkoutOrderFunc func(ctx context.Context, restaurantID int, orderID int64) error
}

func (m *mockOrderUsecase) CheckoutOrder(ctx context.Context, restaurantID int, orderID int64) error {
	if m.checkoutOrderFunc != nil {
		return m.checkoutOrderFunc(ctx, restaurantID, orderID)
	}
	return nil
}

// Tests ProcessPayment
func TestProcessPayment_Success(t *testing.T) {
	mockRepo := &mockSaleRepository{
		createSaleFunc: func(ctx context.Context, s domainSale.Sale) (*domainSale.Sale, error) {
			s.ID = 1
			return &s, nil
		},
	}
	mockOrder := &mockOrderUsecase{}
	usecase := uc.NewUsecase(mockRepo, mockOrder)

	sale, err := usecase.ProcessPayment(context.Background(), 1, 1, 100.0, "cash")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sale.Total != 100.0 {
		t.Errorf("expected total 100.0, got %f", sale.Total)
	}
	if sale.PaymentMethod != "cash" {
		t.Errorf("expected payment method cash, got %s", sale.PaymentMethod)
	}
}

func TestProcessPayment_InvalidOrderID(t *testing.T) {
	mockRepo := &mockSaleRepository{}
	mockOrder := &mockOrderUsecase{}
	usecase := uc.NewUsecase(mockRepo, mockOrder)

	_, err := usecase.ProcessPayment(context.Background(), 0, 1, 100.0, "cash")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domainSale.ErrInvalidOrderID) {
		t.Errorf("expected ErrInvalidOrderID, got %v", err)
	}
}

func TestProcessPayment_EmptyPaymentMethod(t *testing.T) {
	mockRepo := &mockSaleRepository{}
	mockOrder := &mockOrderUsecase{}
	usecase := uc.NewUsecase(mockRepo, mockOrder)

	_, err := usecase.ProcessPayment(context.Background(), 1, 1, 100.0, "")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domainSale.ErrPaymentMethodEmpty) {
		t.Errorf("expected ErrPaymentMethodEmpty, got %v", err)
	}
}

func TestProcessPayment_CheckoutError(t *testing.T) {
	mockRepo := &mockSaleRepository{}
	mockOrder := &mockOrderUsecase{
		checkoutOrderFunc: func(ctx context.Context, restaurantID int, orderID int64) error {
			return errors.New("checkout error")
		},
	}
	usecase := uc.NewUsecase(mockRepo, mockOrder)

	_, err := usecase.ProcessPayment(context.Background(), 1, 1, 100.0, "cash")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// Tests GetInvoice
func TestGetInvoice_Success(t *testing.T) {
	mockRepo := &mockSaleRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*domainSale.Sale, error) {
			return &domainSale.Sale{ID: id, Total: 100.0}, nil
		},
	}
	mockOrder := &mockOrderUsecase{}
	usecase := uc.NewUsecase(mockRepo, mockOrder)

	sale, err := usecase.GetInvoice(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sale.ID != 1 {
		t.Errorf("expected ID 1, got %d", sale.ID)
	}
}

func TestGetInvoice_InvalidID(t *testing.T) {
	mockRepo := &mockSaleRepository{}
	mockOrder := &mockOrderUsecase{}
	usecase := uc.NewUsecase(mockRepo, mockOrder)

	_, err := usecase.GetInvoice(context.Background(), 0)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domainSale.ErrInvalidID) {
		t.Errorf("expected ErrInvalidID, got %v", err)
	}
}

func TestGetInvoice_NotFound(t *testing.T) {
	mockRepo := &mockSaleRepository{
		getByIDFunc: func(ctx context.Context, id int64) (*domainSale.Sale, error) {
			return nil, errors.New("not found")
		},
	}
	mockOrder := &mockOrderUsecase{}
	usecase := uc.NewUsecase(mockRepo, mockOrder)

	_, err := usecase.GetInvoice(context.Background(), 1)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// Tests SplitOrder
func TestSplitOrder_Success(t *testing.T) {
	mockRepo := &mockSaleRepository{}
	mockOrder := &mockOrderUsecase{}
	usecase := uc.NewUsecase(mockRepo, mockOrder)

	result, err := usecase.SplitOrder(context.Background(), 100.0, 4)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.AmountPerPerson != 25.0 {
		t.Errorf("expected 25.0 per person, got %f", result.AmountPerPerson)
	}
}

func TestSplitOrder_InvalidPeople(t *testing.T) {
	mockRepo := &mockSaleRepository{}
	mockOrder := &mockOrderUsecase{}
	usecase := uc.NewUsecase(mockRepo, mockOrder)

	_, err := usecase.SplitOrder(context.Background(), 100.0, 0)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestSplitOrder_InvalidTotal(t *testing.T) {
	mockRepo := &mockSaleRepository{}
	mockOrder := &mockOrderUsecase{}
	usecase := uc.NewUsecase(mockRepo, mockOrder)

	_, err := usecase.SplitOrder(context.Background(), 0, 4)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}