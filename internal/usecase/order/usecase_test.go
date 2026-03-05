package order_test

import (
	"context"
	"errors"
	"testing"

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

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()
	tableID := int64(1)
	items := []domainOrder.OrderItem{
		{ProductID: 1, Quantity: 2, UnitPrice: 10.0},
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		uc := order.NewUsecase(repo)

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
		uc := order.NewUsecase(repo)

		repoErr := errors.New("db error")
		repo.On("Create", ctx, mock.Anything).Return((*domainOrder.Order)(nil), repoErr)

		createdOrder, err := uc.CreateOrder(ctx, 1, 1, &tableID, items)

		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, createdOrder)
		repo.AssertExpectations(t)
	})

	t.Run("invalid order - no items", func(t *testing.T) {
		repo := new(MockRepository)
		uc := order.NewUsecase(repo)

		createdOrder, err := uc.CreateOrder(ctx, 1, 1, &tableID, []domainOrder.OrderItem{})

		assert.ErrorIs(t, err, domainOrder.ErrInvalidOrderItems)
		assert.Nil(t, createdOrder)
		repo.AssertExpectations(t)
	})
}

func TestAddProductToOrder(t *testing.T) {
	repo := new(MockRepository)
	uc := order.NewUsecase(repo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		err := uc.AddProductToOrder(ctx, 1, 1, []domainOrder.OrderItem{})
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetOrderByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		uc := order.NewUsecase(repo)
		expectedOrder := &domainOrder.Order{ID: 1, StatusID: 1}

		repo.On("GetByID", ctx, int64(1), 1).Return(expectedOrder, nil)

		o, err := uc.GetOrderByID(ctx, 1, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, o)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockRepository)
		uc := order.NewUsecase(repo)

		repo.On("GetByID", ctx, int64(99), 1).Return((*domainOrder.Order)(nil), errors.New("not found"))

		o, err := uc.GetOrderByID(ctx, 1, 99)

		assert.Error(t, err)
		assert.Nil(t, o)
		repo.AssertExpectations(t)
	})
}

func TestUpdateOrderStatus_And_Checkout(t *testing.T) {
	ctx := context.Background()

	t.Run("success update status", func(t *testing.T) {
		repo := new(MockRepository)
		uc := order.NewUsecase(repo)

		repo.On("UpdateStatus", ctx, int64(1), 1, 3).Return(nil)

		err := uc.UpdateOrderStatus(ctx, 1, 1, 3)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("success checkout", func(t *testing.T) {
		repo := new(MockRepository)
		uc := order.NewUsecase(repo)

		repo.On("UpdateStatus", ctx, int64(1), 1, 5).Return(nil)

		err := uc.CheckoutOrder(ctx, 1, 1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestListOrdersByTable(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		uc := order.NewUsecase(repo)

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
}
