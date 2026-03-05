package payment

import (
	"context"
	"fmt"
	"time"

	domainPayment "github.com/AXONcompany/POS/internal/domain/payment"
)

type Repository interface {
	Create(ctx context.Context, p *domainPayment.Payment) (*domainPayment.Payment, error)
	GetByID(ctx context.Context, id int64) (*domainPayment.Payment, error)
	GetByOrderID(ctx context.Context, orderID int64) ([]*domainPayment.Payment, error)
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

// ProcessPayment procesa un pago para una orden.
func (uc *Usecase) ProcessPayment(ctx context.Context, orderID int64, divisionID *string, method string, amount, tip float64, reference string, venueID, userID int) (*domainPayment.Payment, error) {
	total := amount + tip

	p := &domainPayment.Payment{
		OrderID:       orderID,
		DivisionID:    divisionID,
		PaymentMethod: method,
		Amount:        amount,
		Tip:           tip,
		Total:         total,
		Status:        "aprobado",
		Reference:     reference,
		VenueID:       venueID,
		UserID:        userID,
	}

	return uc.repo.Create(ctx, p)
}

// GetPayment obtiene un pago por ID.
func (uc *Usecase) GetPayment(ctx context.Context, id int64) (*domainPayment.Payment, error) {
	return uc.repo.GetByID(ctx, id)
}

// GenerateInvoice genera los datos de factura para un pago.
func (uc *Usecase) GenerateInvoice(ctx context.Context, paymentID int64) (map[string]interface{}, error) {
	p, err := uc.repo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	invoice := map[string]interface{}{
		"factura_id": fmt.Sprintf("FACT-%d-%06d", time.Now().Year(), p.ID),
		"fecha":      p.CreatedAt,
		"pago_id":    p.ID,
		"orden_id":   p.OrderID,
		"monto":      p.Amount,
		"propina":    p.Tip,
		"total":      p.Total,
		"metodo":     p.PaymentMethod,
		"referencia": p.Reference,
		"codigo":     fmt.Sprintf("%x", p.CreatedAt.UnixNano()),
		"url_pdf":    fmt.Sprintf("/facturas/FACT-%d-%06d.pdf", time.Now().Year(), p.ID),
	}

	return invoice, nil
}
