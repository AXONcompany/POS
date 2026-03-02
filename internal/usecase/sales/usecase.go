package sale

import (
	"context"
	"fmt"
	"time"

	domainSale "github.com/AXONcompany/POS/internal/domain/sale"
)

type SaleRepository interface {
	CreateSale(ctx context.Context, s domainSale.Sale) (*domainSale.Sale, error)
	GetByID(ctx context.Context, id int64) (*domainSale.Sale, error)
}

type OrderUsecase interface {
	CheckoutOrder(ctx context.Context, restaurantID int, orderID int64) error
}

type Usecase struct {
	saleRepo     SaleRepository
	orderUsecase OrderUsecase
}

func NewUsecase(saleRepo SaleRepository, orderUsecase OrderUsecase) *Usecase {
	return &Usecase{
		saleRepo:     saleRepo,
		orderUsecase: orderUsecase,
	}
}

// POST /payments
func (u *Usecase) ProcessPayment(ctx context.Context, orderID int64, restaurantID int, total float64, paymentMethod string) (*domainSale.Sale, error) {
	if orderID <= 0 {
		return nil, domainSale.ErrInvalidOrderID
	}
	if paymentMethod == "" {
		return nil, domainSale.ErrPaymentMethodEmpty
	}
	if total <= 0 {
		return nil, domainSale.ErrInvalidTotal
	}

	if err := u.orderUsecase.CheckoutOrder(ctx, restaurantID, orderID); err != nil {
		return nil, err
	}

	s := domainSale.Sale{
		Total:         total,
		PaymentMethod: paymentMethod,
		Date:          time.Now(),
		OrderID:       int64(orderID),
	}

	created, err := u.saleRepo.CreateSale(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("Faild to create sale: %w", err)
	}

	return created, nil
}

// GET /payments/:id/invoice
func (u *Usecase) GetInvoice(ctx context.Context, saleID int64) (*domainSale.Sale, error) {
	if saleID <= 0 {
		return nil, domainSale.ErrInvalidID
	}

	sale, err := u.saleRepo.GetByID(ctx, saleID)
	if err != nil {
		return nil, domainSale.ErrSaleNotFound
	}

	return sale, nil
}

// POST /orders/:id/split
func (u *Usecase) SplitOrder(ctx context.Context, total float64, people int) (*SplitResult, error) {
	if people <= 0 {
		return nil, fmt.Errorf("people must be greater than 0")
	}
	if total <= 0 {
		return nil, fmt.Errorf("total must be greater than 0")
	}

	return &SplitResult{
		Total:           total,
		People:          people,
		AmountPerPerson: total / float64(people),
	}, nil
}

type SplitResult struct {
	Total           float64 `json:"total"`
	People          int     `json:"personas"`
	AmountPerPerson float64 `json:"monto_por_persona"`
}
