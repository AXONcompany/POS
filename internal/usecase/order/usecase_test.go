package order_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domainAudit "github.com/AXONcompany/POS/internal/domain/audit"
	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
	"github.com/AXONcompany/POS/internal/usecase/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a testify mock for order.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, o *domainOrder.Order) (*domainOrder.Order, error) {
	args := m.Called(ctx, o)
	if args.Get(0) != nil {
		return args.Get(0).(*domainOrder.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64, venueID int) (*domainOrder.Order, error) {
	args := m.Called(ctx, id, venueID)
	if args.Get(0) != nil {
		return args.Get(0).(*domainOrder.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) GetStatusByID(ctx context.Context, id int64, venueID int) (int, error) {
	args := m.Called(ctx, id, venueID)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) UpdateStatus(ctx context.Context, id int64, venueID int, statusID int) error {
	args := m.Called(ctx, id, venueID, statusID)
	return args.Error(0)
}

func (m *MockRepository) ListByTable(ctx context.Context, tableID int64, venueID int) ([]domainOrder.Order, error) {
	args := m.Called(ctx, tableID, venueID)
	if args.Get(0) != nil {
		return args.Get(0).([]domainOrder.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) AddItemsWithInventory(ctx context.Context, orderID int64, venueID int, items []domainOrder.OrderItem, deductions []domainOrder.StockDeduction) error {
	args := m.Called(ctx, orderID, venueID, items, deductions)
	return args.Error(0)
}

func (m *MockRepository) GetOrderItem(ctx context.Context, itemID, orderID int64) (*domainOrder.OrderItem, error) {
	args := m.Called(ctx, itemID, orderID)
	if args.Get(0) != nil {
		return args.Get(0).(*domainOrder.OrderItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) CancelItemWithInventoryRestore(ctx context.Context, itemID, orderID int64, venueID int, restorations []domainOrder.StockDeduction) error {
	args := m.Called(ctx, itemID, orderID, venueID, restorations)
	return args.Error(0)
}

// MockAuditRepository is a testify mock for order.AuditRepository
type MockAuditRepository struct {
	mock.Mock
}

func (m *MockAuditRepository) SaveAudit(ctx context.Context, entry *domainAudit.AuditEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

// MockProductInventoryRepository is a testify mock for order.ProductInventoryRepository
type MockProductInventoryRepository struct {
	mock.Mock
}

func (m *MockProductInventoryRepository) GetProductPrice(ctx context.Context, productID int64, venueID int) (float64, error) {
	args := m.Called(ctx, productID, venueID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockProductInventoryRepository) GetRecipeLines(ctx context.Context, productID int64) ([]domainOrder.RecipeLine, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) != nil {
		return args.Get(0).([]domainOrder.RecipeLine), args.Error(1)
	}
	return nil, args.Error(1)
}

func newUsecase(repo *MockRepository, invRepo *MockProductInventoryRepository) *order.Usecase {
	return order.NewUsecase(repo, invRepo, new(MockAuditRepository))
}

func newUsecaseWithAudit(repo *MockRepository, invRepo *MockProductInventoryRepository, auditRepo *MockAuditRepository) *order.Usecase {
	return order.NewUsecase(repo, invRepo, auditRepo)
}

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()
	tableID := int64(1)
	items := []domainOrder.OrderItem{
		{ProductID: 1, Quantity: 2, UnitPrice: 10.0},
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		expectedOrder := &domainOrder.Order{
			VenueID:     1,
			UserID:      1,
			TableID:     &tableID,
			StatusID:    1,
			TotalAmount: 20.0,
			Items:       items,
		}

		repo.On("Create", ctx, expectedOrder).Return(expectedOrder, nil)

		createdOrder, err := uc.CreateOrder(ctx, 1, 1, &tableID, items)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, createdOrder)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		repoErr := errors.New("db error")
		repo.On("Create", ctx, mock.Anything).Return((*domainOrder.Order)(nil), repoErr)

		createdOrder, err := uc.CreateOrder(ctx, 1, 1, &tableID, items)

		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, createdOrder)
		repo.AssertExpectations(t)
	})

	t.Run("invalid order - no items", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		createdOrder, err := uc.CreateOrder(ctx, 1, 1, &tableID, []domainOrder.OrderItem{})

		assert.ErrorIs(t, err, domainOrder.ErrInvalidOrderItems)
		assert.Nil(t, createdOrder)
		repo.AssertExpectations(t)
	})
}

func TestAddProductToOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("success with stock deduction", func(t *testing.T) {
		repo := new(MockRepository)
		invRepo := new(MockProductInventoryRepository)
		uc := newUsecase(repo, invRepo)

		items := []domainOrder.OrderItem{
			{ProductID: 10, Quantity: 2},
		}
		recipe := []domainOrder.RecipeLine{
			{IngredientID: 5, QuantityRequired: 100},
		}
		expectedDeductions := []domainOrder.StockDeduction{
			{IngredientID: 5, VenueID: 1, Quantity: 200},
		}
		expectedItems := []domainOrder.OrderItem{
			{ProductID: 10, Quantity: 2, UnitPrice: 15.0},
		}

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(1, nil)
		invRepo.On("GetProductPrice", ctx, int64(10), 1).Return(15.0, nil)
		invRepo.On("GetRecipeLines", ctx, int64(10)).Return(recipe, nil)
		repo.On("AddItemsWithInventory", ctx, int64(1), 1, expectedItems, expectedDeductions).Return(nil)

		err := uc.AddProductToOrder(ctx, 1, 1, items)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
		invRepo.AssertExpectations(t)
	})

	t.Run("reject when order not in editable state", func(t *testing.T) {
		repo := new(MockRepository)
		invRepo := new(MockProductInventoryRepository)
		uc := newUsecase(repo, invRepo)

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(3, nil) // PREPARING

		err := uc.AddProductToOrder(ctx, 1, 1, []domainOrder.OrderItem{{ProductID: 10, Quantity: 1}})

		assert.ErrorIs(t, err, domainOrder.ErrInvalidStatusTransition)
		repo.AssertExpectations(t)
		invRepo.AssertExpectations(t)
	})

	t.Run("propagate ErrInsufficientStock from repo", func(t *testing.T) {
		repo := new(MockRepository)
		invRepo := new(MockProductInventoryRepository)
		uc := newUsecase(repo, invRepo)

		items := []domainOrder.OrderItem{{ProductID: 10, Quantity: 1}}

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(2, nil) // SENT
		invRepo.On("GetProductPrice", ctx, int64(10), 1).Return(10.0, nil)
		invRepo.On("GetRecipeLines", ctx, int64(10)).Return([]domainOrder.RecipeLine{
			{IngredientID: 5, QuantityRequired: 50},
		}, nil)
		repo.On("AddItemsWithInventory", ctx, int64(1), 1, mock.Anything, mock.Anything).Return(domainOrder.ErrInsufficientStock)

		err := uc.AddProductToOrder(ctx, 1, 1, items)

		assert.ErrorIs(t, err, domainOrder.ErrInsufficientStock)
		repo.AssertExpectations(t)
		invRepo.AssertExpectations(t)
	})

	t.Run("accumulate deductions for shared ingredient", func(t *testing.T) {
		repo := new(MockRepository)
		invRepo := new(MockProductInventoryRepository)
		uc := newUsecase(repo, invRepo)

		items := []domainOrder.OrderItem{
			{ProductID: 10, Quantity: 1},
			{ProductID: 20, Quantity: 1},
		}
		// Ambos productos usan el ingrediente 5
		recipe10 := []domainOrder.RecipeLine{{IngredientID: 5, QuantityRequired: 100}}
		recipe20 := []domainOrder.RecipeLine{{IngredientID: 5, QuantityRequired: 50}}

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(1, nil)
		invRepo.On("GetProductPrice", ctx, int64(10), 1).Return(10.0, nil)
		invRepo.On("GetProductPrice", ctx, int64(20), 1).Return(20.0, nil)
		invRepo.On("GetRecipeLines", ctx, int64(10)).Return(recipe10, nil)
		invRepo.On("GetRecipeLines", ctx, int64(20)).Return(recipe20, nil)
		// La deduccion acumulada debe ser 150 (100 + 50)
		repo.On("AddItemsWithInventory", ctx, int64(1), 1, mock.Anything,
			mock.MatchedBy(func(deductions []domainOrder.StockDeduction) bool {
				if len(deductions) != 1 {
					return false
				}
				return deductions[0].IngredientID == 5 && deductions[0].Quantity == 150
			}),
		).Return(nil)

		err := uc.AddProductToOrder(ctx, 1, 1, items)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
		invRepo.AssertExpectations(t)
	})
}

func TestGetOrderByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))
		expectedOrder := &domainOrder.Order{ID: 1, StatusID: 1}

		repo.On("GetByID", ctx, int64(1), 1).Return(expectedOrder, nil)

		o, err := uc.GetOrderByID(ctx, 1, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, o)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		repo.On("GetByID", ctx, int64(99), 1).Return((*domainOrder.Order)(nil), errors.New("not found"))

		o, err := uc.GetOrderByID(ctx, 1, 99)

		assert.Error(t, err)
		assert.Nil(t, o)
		repo.AssertExpectations(t)
	})
}

func TestUpdateOrderStatus_And_Checkout(t *testing.T) {
	ctx := context.Background()

	t.Run("success update status SENT->PREPARING", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(2, nil) // current: SENT
		repo.On("UpdateStatus", ctx, int64(1), 1, 3).Return(nil)

		err := uc.UpdateOrderStatus(ctx, 1, 1, 3)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("invalid transition PAID->PENDING", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(5, nil) // current: PAID

		err := uc.UpdateOrderStatus(ctx, 1, 1, 1)

		assert.ErrorIs(t, err, domainOrder.ErrInvalidStatusTransition)
		repo.AssertExpectations(t)
	})

	t.Run("success checkout READY->PAID", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(4, nil) // current: READY
		repo.On("UpdateStatus", ctx, int64(1), 1, 5).Return(nil)

		err := uc.CheckoutOrder(ctx, 1, 1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("checkout invalid when already PAID", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(5, nil) // current: PAID (terminal)

		err := uc.CheckoutOrder(ctx, 1, 1)

		assert.ErrorIs(t, err, domainOrder.ErrInvalidStatusTransition)
		repo.AssertExpectations(t)
	})
}

func TestListOrdersByTable(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		tableID := int64(1)
		expectedOrders := []domainOrder.Order{
			{ID: 1, TableID: &tableID},
		}

		repo.On("ListByTable", ctx, int64(1), 1).Return(expectedOrders, nil)

		list, err := uc.ListOrdersByTable(ctx, 1, 1)

		assert.NoError(t, err)
		assert.Len(t, list, 1)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		repo.On("ListByTable", ctx, int64(1), 1).Return(nil, errors.New("db error"))

		list, err := uc.ListOrdersByTable(ctx, 1, 1)

		assert.Error(t, err)
		assert.Nil(t, list)
		repo.AssertExpectations(t)
	})
}

func TestCreateOrderWithoutItems(t *testing.T) {
	ctx := context.Background()
	tableID := int64(3)

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		expected := &domainOrder.Order{
			VenueID:     1,
			UserID:      2,
			TableID:     &tableID,
			StatusID:    1,
			TotalAmount: 0,
		}
		repo.On("Create", ctx, expected).Return(expected, nil)

		result, err := uc.CreateOrderWithoutItems(ctx, 1, 2, &tableID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, float64(0), result.TotalAmount)
		assert.Equal(t, 1, result.StatusID)
		repo.AssertExpectations(t)
	})

	t.Run("nil tableID (take-away)", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		expected := &domainOrder.Order{
			VenueID:     1,
			UserID:      1,
			TableID:     nil,
			StatusID:    1,
			TotalAmount: 0,
		}
		repo.On("Create", ctx, expected).Return(expected, nil)

		result, err := uc.CreateOrderWithoutItems(ctx, 1, 1, nil)

		assert.NoError(t, err)
		assert.Nil(t, result.TableID)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		repo.On("Create", ctx, mock.Anything).Return((*domainOrder.Order)(nil), errors.New("connection timeout"))

		_, err := uc.CreateOrderWithoutItems(ctx, 1, 1, &tableID)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestUpdateOrderStatus_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	uc := newUsecase(repo, new(MockProductInventoryRepository))

	// current: PENDING(1) → next: SENT(2) es válido, pero UpdateStatus falla
	repo.On("GetStatusByID", ctx, int64(1), 1).Return(1, nil)
	repo.On("UpdateStatus", ctx, int64(1), 1, 2).Return(errors.New("order not found"))

	err := uc.UpdateOrderStatus(ctx, 1, 1, 2)

	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCheckoutOrder_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := new(MockRepository)
	uc := newUsecase(repo, new(MockProductInventoryRepository))

	// current: READY(4) → PAID(5) es válido, pero UpdateStatus falla
	repo.On("GetStatusByID", ctx, int64(1), 1).Return(4, nil)
	repo.On("UpdateStatus", ctx, int64(1), 1, 5).Return(errors.New("order already closed"))

	err := uc.CheckoutOrder(ctx, 1, 1)

	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCancelOrderItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success: cancela item y restaura stock", func(t *testing.T) {
		repo := new(MockRepository)
		invRepo := new(MockProductInventoryRepository)
		auditRepo := new(MockAuditRepository)
		uc := newUsecaseWithAudit(repo, invRepo, auditRepo)

		item := &domainOrder.OrderItem{
			ID:        10,
			OrderID:   1,
			ProductID: 5,
			Quantity:  2,
			UnitPrice: 15.0,
		}
		recipe := []domainOrder.RecipeLine{
			{IngredientID: 3, QuantityRequired: 100},
		}
		expectedRestorations := []domainOrder.StockDeduction{
			{IngredientID: 3, VenueID: 1, Quantity: 200},
		}

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(2, nil) // SENT
		repo.On("GetOrderItem", ctx, int64(10), int64(1)).Return(item, nil)
		invRepo.On("GetRecipeLines", ctx, int64(5)).Return(recipe, nil)
		repo.On("CancelItemWithInventoryRestore", ctx, int64(10), int64(1), 1, expectedRestorations).Return(nil)
		auditRepo.On("SaveAudit", ctx, mock.MatchedBy(func(e *domainAudit.AuditEntry) bool {
			return e.EntityType == "order_item" && e.EntityID == 10 && e.Action == "cancel" &&
				e.UserID == 42 && e.VenueID == 1
		})).Return(nil)

		err := uc.CancelOrderItem(ctx, 1, 42, 1, 10)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
		invRepo.AssertExpectations(t)
		auditRepo.AssertExpectations(t)
	})

	t.Run("rechazo: orden en estado PAID", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(5, nil) // PAID

		err := uc.CancelOrderItem(ctx, 1, 1, 1, 10)

		assert.ErrorIs(t, err, domainOrder.ErrInvalidStatusTransition)
		repo.AssertExpectations(t)
	})

	t.Run("rechazo: item ya cancelado", func(t *testing.T) {
		repo := new(MockRepository)
		uc := newUsecase(repo, new(MockProductInventoryRepository))

		cancelledAt := func() *time.Time { ts := time.Now(); return &ts }()
		alreadyCancelledItem := &domainOrder.OrderItem{
			ID:          10,
			OrderID:     1,
			ProductID:   5,
			Quantity:    1,
			CancelledAt: cancelledAt,
		}

		repo.On("GetStatusByID", ctx, int64(1), 1).Return(1, nil) // PENDING
		repo.On("GetOrderItem", ctx, int64(10), int64(1)).Return(alreadyCancelledItem, nil)

		err := uc.CancelOrderItem(ctx, 1, 1, 1, 10)

		assert.ErrorIs(t, err, domainOrder.ErrItemAlreadyCancelled)
		repo.AssertExpectations(t)
	})
}
