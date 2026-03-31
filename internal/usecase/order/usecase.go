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
}

type ProductInventoryRepository interface {
	GetProductPrice(ctx context.Context, productID int64, venueID int) (float64, error)
	GetRecipeLines(ctx context.Context, productID int64) ([]domainOrder.RecipeLine, error)
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

	// Acumular deducciones por ingrediente y establecer precio real en cada item
	deductionMap := make(map[int64]float64)
	for i := range items {
		price, err := uc.invRepo.GetProductPrice(ctx, items[i].ProductID, venueID)
		if err != nil {
			return fmt.Errorf("get product price: %w", err)
		}
		items[i].UnitPrice = price

		lines, err := uc.invRepo.GetRecipeLines(ctx, items[i].ProductID)
		if err != nil {
			return fmt.Errorf("get recipe lines: %w", err)
		}
		for _, line := range lines {
			deductionMap[line.IngredientID] += line.QuantityRequired * float64(items[i].Quantity)
		}
	}

	deductions := make([]domainOrder.StockDeduction, 0, len(deductionMap))
	for ingredientID, qty := range deductionMap {
		deductions = append(deductions, domainOrder.StockDeduction{
			IngredientID: ingredientID,
			VenueID:      venueID,
			Quantity:     qty,
		})
	}

	return uc.repo.AddItemsWithInventory(ctx, orderID, venueID, items, deductions)
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

	// Calcular restauraciones de stock
	lines, err := uc.invRepo.GetRecipeLines(ctx, item.ProductID)
	if err != nil {
		return fmt.Errorf("get recipe lines: %w", err)
	}
	restorations := make([]domainOrder.StockDeduction, 0, len(lines))
	for _, line := range lines {
		restorations = append(restorations, domainOrder.StockDeduction{
			IngredientID: line.IngredientID,
			VenueID:      venueID,
			Quantity:     line.QuantityRequired * float64(item.Quantity),
		})
	}

	// TX atómica: cancelar item + restaurar stock + ajustar total
	if err := uc.repo.CancelItemWithInventoryRestore(ctx, itemID, orderID, venueID, restorations); err != nil {
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
