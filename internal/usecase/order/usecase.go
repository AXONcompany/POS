package order

import (
	"context"
	"fmt"
	"time"

	domainAudit "github.com/AXONcompany/POS/internal/domain/audit"
	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
)

type Repository interface {
	Create(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error)
	GetByID(ctx context.Context, id int64, venueID int) (*domainOrder.Order, error)
	GetStatusByID(ctx context.Context, id int64, venueID int) (int, error)
	GetOrderItem(ctx context.Context, itemID, orderID int64) (*domainOrder.OrderItem, error)
	UpdateStatus(ctx context.Context, id int64, venueID int, statusID int) error
	ListByTable(ctx context.Context, tableID int64, venueID int) ([]domainOrder.Order, error)
	AddItemsWithInventory(ctx context.Context, orderID int64, venueID int, items []domainOrder.OrderItem, deductions []domainOrder.StockDeduction) error
	CancelItemWithInventoryRestore(ctx context.Context, itemID, orderID int64, venueID int, restorations []domainOrder.StockDeduction) error
	CreateDivisions(ctx context.Context, divisions []domainOrder.OrderDivision) error
	GetDivisionsByOrderID(ctx context.Context, orderID int64, venueID int) ([]domainOrder.OrderDivision, error)
}

type ProductInventoryRepository interface {
	GetProductPrice(ctx context.Context, productID int64, venueID int) (float64, error)
}

type AuditRepository interface {
	SaveAudit(ctx context.Context, entry *domainAudit.AuditEntry) error
}

type Usecase struct {
	repo      Repository
	invRepo   ProductInventoryRepository
	auditRepo AuditRepository
}

func NewUsecase(repo Repository, invRepo ProductInventoryRepository, auditRepo AuditRepository) *Usecase {
	return &Usecase{repo: repo, invRepo: invRepo, auditRepo: auditRepo}
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
	current, err := uc.repo.GetStatusByID(ctx, orderID, venueID)
	if err != nil {
		return err
	}
	// Solo se permiten ordenes en estado PENDING(1) o SENT(2)
	if current != 1 && current != 2 {
		return domainOrder.ErrInvalidStatusTransition
	}

	for i := range items {
		price, err := uc.invRepo.GetProductPrice(ctx, items[i].ProductID, venueID)
		if err != nil {
			return fmt.Errorf("get product price: %w", err)
		}
		items[i].UnitPrice = price
	}

	return uc.repo.AddItemsWithInventory(ctx, orderID, venueID, items, nil)
}

func (uc *Usecase) GetOrderByID(ctx context.Context, venueID int, orderID int64) (*domainOrder.Order, error) {
	return uc.repo.GetByID(ctx, orderID, venueID)
}

func (uc *Usecase) CheckoutOrder(ctx context.Context, venueID int, orderID int64) error {
	current, err := uc.repo.GetStatusByID(ctx, orderID, venueID)
	if err != nil {
		return err
	}
	if !domainOrder.CanTransitionTo(current, 5) {
		return domainOrder.ErrInvalidStatusTransition
	}
	return uc.repo.UpdateStatus(ctx, orderID, venueID, 5)
}

func (uc *Usecase) ListOrdersByTable(ctx context.Context, venueID int, tableID int64) ([]domainOrder.Order, error) {
	return uc.repo.ListByTable(ctx, tableID, venueID)
}

func (uc *Usecase) UpdateOrderStatus(ctx context.Context, venueID int, orderID int64, statusID int) error {
	current, err := uc.repo.GetStatusByID(ctx, orderID, venueID)
	if err != nil {
		return err
	}
	if !domainOrder.CanTransitionTo(current, statusID) {
		return domainOrder.ErrInvalidStatusTransition
	}
	return uc.repo.UpdateStatus(ctx, orderID, venueID, statusID)
}

// CancelOrderItem cancela un item de una orden, restaura el stock y registra auditoría.
func (uc *Usecase) CancelOrderItem(ctx context.Context, venueID, userID int, orderID, itemID int64) error {
	// Validar estado de la orden — no se puede cancelar items de ordenes PAID o CANCELLED
	current, err := uc.repo.GetStatusByID(ctx, orderID, venueID)
	if err != nil {
		return err
	}
	if current == 5 || current == 6 {
		return domainOrder.ErrInvalidStatusTransition
	}

	// Snapshot "before"
	item, err := uc.repo.GetOrderItem(ctx, itemID, orderID)
	if err != nil {
		return err
	}
	if item.CancelledAt != nil {
		return domainOrder.ErrItemAlreadyCancelled
	}
	before := *item

	// TX atómica: cancelar item + ajustar total
	if err := uc.repo.CancelItemWithInventoryRestore(ctx, itemID, orderID, venueID, nil); err != nil {
		return err
	}

	// Snapshot "after"
	now := time.Now()
	after := before
	after.CancelledAt = &now

	// Registro de auditoría
	return uc.auditRepo.SaveAudit(ctx, &domainAudit.AuditEntry{
		EntityType: "order_item",
		EntityID:   itemID,
		Action:     "cancel",
		OldValue:   before,
		NewValue:   after,
		UserID:     userID,
		VenueID:    venueID,
	})
}

func (uc *Usecase) DivideOrder(ctx context.Context, venueID int, orderID int64, divisionType string, numParts int, customAmounts []float64) ([]domainOrder.OrderDivision, error) {
	order, err := uc.repo.GetByID(ctx, orderID, venueID)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	subtotal := order.TotalAmount
	impuestos := subtotal * 0.19
	total := subtotal + impuestos

	var divisions []domainOrder.OrderDivision

	switch divisionType {
	case "partes_iguales":
		parts := numParts
		if parts <= 0 {
			parts = 2
		}
		partSubtotal := subtotal / float64(parts)
		partImpuestos := impuestos / float64(parts)
		partTotal := total / float64(parts)

		divisions = make([]domainOrder.OrderDivision, parts)
		for i := 0; i < parts; i++ {
			divisions[i] = domainOrder.OrderDivision{
				ID:           fmt.Sprintf("div_%d_%d", orderID, i+1),
				OrderID:      orderID,
				VenueID:      venueID,
				DivisionType: divisionType,
				Amount:       partSubtotal,
				Tax:          partImpuestos,
				Total:        partTotal,
				IsPaid:       false,
			}
		}

	case "por_monto":
		remaining := total
		divisions = make([]domainOrder.OrderDivision, 0, len(customAmounts))
		for i, amount := range customAmounts {
			divTotal := amount
			if divTotal > remaining {
				divTotal = remaining
			}
			divSubtotal := divTotal / 1.19
			divImpuestos := divTotal - divSubtotal

			divisions = append(divisions, domainOrder.OrderDivision{
				ID:           fmt.Sprintf("div_%d_%d", orderID, i+1),
				OrderID:      orderID,
				VenueID:      venueID,
				DivisionType: divisionType,
				Amount:       divSubtotal,
				Tax:          divImpuestos,
				Total:        divTotal,
				IsPaid:       false,
			})
			remaining -= divTotal
		}

	case "por_item":
		parts := len(customAmounts)
		if parts == 0 {
			parts = 1
		}
		divTotal := total / float64(parts)
		divSubtotal := divTotal / 1.19
		divImpuestos := divTotal - divSubtotal

		divisions = make([]domainOrder.OrderDivision, 0, parts)
		for i := 0; i < parts; i++ {
			divisions = append(divisions, domainOrder.OrderDivision{
				ID:           fmt.Sprintf("div_%d_%d", orderID, i+1),
				OrderID:      orderID,
				VenueID:      venueID,
				DivisionType: divisionType,
				Amount:       divSubtotal,
				Tax:          divImpuestos,
				Total:        divTotal,
				IsPaid:       false,
			})
		}

	default:
		return nil, fmt.Errorf("invalid division type: %s", divisionType)
	}

	err = uc.repo.CreateDivisions(ctx, divisions)
	if err != nil {
		return nil, fmt.Errorf("create divisions: %w", err)
	}

	return divisions, nil
}

func (uc *Usecase) GetDivisionsByOrder(ctx context.Context, venueID int, orderID int64) ([]domainOrder.OrderDivision, error) {
	return uc.repo.GetDivisionsByOrderID(ctx, orderID, venueID)
}
