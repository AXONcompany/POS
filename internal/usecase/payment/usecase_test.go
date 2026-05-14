package payment_test

import (
	"context"
	"errors"
	"testing"

	domainPayment "github.com/AXONcompany/POS/internal/domain/payment"
	uc "github.com/AXONcompany/POS/internal/usecase/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) Create(ctx context.Context, p *domainPayment.Payment) (*domainPayment.Payment, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainPayment.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetByID(ctx context.Context, id int64) (*domainPayment.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainPayment.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetByOrderID(ctx context.Context, orderID int64) ([]*domainPayment.Payment, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).([]*domainPayment.Payment), args.Error(1)
}

// --- ProcessPayment ---

func TestProcessPayment_TotalCalculation(t *testing.T) {
	repo := new(MockPaymentRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	amount := 100.0
	tip := 15.0
	expectedTotal := 115.0

	repo.On("Create", ctx, mock.MatchedBy(func(p *domainPayment.Payment) bool {
		return p.Amount == amount && p.Tip == tip && p.Total == expectedTotal
	})).Return(&domainPayment.Payment{
		ID: 1, Amount: amount, Tip: tip, Total: expectedTotal,
	}, nil)

	result, err := service.ProcessPayment(ctx, 1, nil, "efectivo", amount, tip, "", 1, 1)
	assert.NoError(t, err)
	assert.InDelta(t, expectedTotal, result.Total, 0.001)
	repo.AssertExpectations(t)
}

func TestProcessPayment_ZeroTip(t *testing.T) {
	repo := new(MockPaymentRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	amount := 50.0
	tip := 0.0

	repo.On("Create", ctx, mock.MatchedBy(func(p *domainPayment.Payment) bool {
		return p.Total == 50.0 && p.Tip == 0
	})).Return(&domainPayment.Payment{ID: 1, Amount: amount, Tip: tip, Total: 50.0}, nil)

	result, err := service.ProcessPayment(ctx, 1, nil, "tarjeta", amount, tip, "", 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, float64(50), result.Total)
	repo.AssertExpectations(t)
}

func TestProcessPayment_WithDivisionID(t *testing.T) {
	repo := new(MockPaymentRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	divID := "div-abc-123"
	repo.On("Create", ctx, mock.MatchedBy(func(p *domainPayment.Payment) bool {
		return p.DivisionID != nil && *p.DivisionID == divID
	})).Return(&domainPayment.Payment{ID: 1, DivisionID: &divID}, nil)

	_, err := service.ProcessPayment(ctx, 1, &divID, "efectivo", 30.0, 0, "", 1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProcessPayment_StatusIsAprobado(t *testing.T) {
	repo := new(MockPaymentRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	repo.On("Create", ctx, mock.MatchedBy(func(p *domainPayment.Payment) bool {
		return p.Status == "aprobado"
	})).Return(&domainPayment.Payment{ID: 1, Status: "aprobado"}, nil)

	result, err := service.ProcessPayment(ctx, 1, nil, "efectivo", 20.0, 0, "", 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "aprobado", result.Status)
	repo.AssertExpectations(t)
}

func TestProcessPayment_RepoError(t *testing.T) {
	repo := new(MockPaymentRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	repo.On("Create", ctx, mock.Anything).Return(nil, errors.New("db error"))

	_, err := service.ProcessPayment(ctx, 1, nil, "efectivo", 50.0, 5.0, "", 1, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// --- GenerateInvoice ---

func TestGenerateInvoice_NotFound(t *testing.T) {
	repo := new(MockPaymentRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, int64(99)).Return(nil, errors.New("not found"))

	_, err := service.GenerateInvoice(ctx, 99)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payment not found")
	repo.AssertExpectations(t)
}

func TestGenerateInvoice_ContainsRequiredFields(t *testing.T) {
	repo := new(MockPaymentRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	payment := &domainPayment.Payment{
		ID:            5,
		OrderID:       10,
		Amount:        80.0,
		Tip:           8.0,
		Total:         88.0,
		PaymentMethod: "tarjeta",
		Reference:     "REF-001",
	}
	repo.On("GetByID", ctx, int64(5)).Return(payment, nil)

	invoice, err := service.GenerateInvoice(ctx, 5)
	assert.NoError(t, err)
	assert.NotNil(t, invoice)

	// Verificar campos obligatorios de la factura
	assert.Contains(t, invoice, "factura_id")
	assert.Contains(t, invoice, "pago_id")
	assert.Contains(t, invoice, "orden_id")
	assert.Contains(t, invoice, "monto")
	assert.Contains(t, invoice, "propina")
	assert.Contains(t, invoice, "total")
	assert.Contains(t, invoice, "metodo")
	assert.Equal(t, float64(80), invoice["monto"])
	assert.Equal(t, float64(8), invoice["propina"])
	assert.Equal(t, float64(88), invoice["total"])
	repo.AssertExpectations(t)
}
