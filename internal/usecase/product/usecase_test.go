package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/AXONcompany/POS/internal/domain/product"
	uc "github.com/AXONcompany/POS/internal/usecase/product"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) CreateProduct(ctx context.Context, p product.Product) (*product.Product, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductRepo) GetByID(ctx context.Context, id int64, venueID int) (*product.Product, error) {
	args := m.Called(ctx, id, venueID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductRepo) GetAllProducts(ctx context.Context, venueID int, page, pageSize int) ([]product.Product, error) {
	args := m.Called(ctx, venueID, page, pageSize)
	return args.Get(0).([]product.Product), args.Error(1)
}

func (m *MockProductRepo) UpdateProduct(ctx context.Context, p product.Product) (*product.Product, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductRepo) DeleteProduct(ctx context.Context, id int64, venueID int) error {
	args := m.Called(ctx, id, venueID)
	return args.Error(0)
}

type MockCategoryRepo struct {
	mock.Mock
}

func (m *MockCategoryRepo) CreateCategory(ctx context.Context, c product.Category) (*product.Category, error) {
	args := m.Called(ctx, c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Category), args.Error(1)
}

func (m *MockCategoryRepo) GetByID(ctx context.Context, id int64, venueID int) (*product.Category, error) {
	args := m.Called(ctx, id, venueID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Category), args.Error(1)
}

func (m *MockCategoryRepo) GetAllCategories(ctx context.Context, venueID int, page, pageSize int) ([]product.Category, error) {
	args := m.Called(ctx, venueID, page, pageSize)
	return args.Get(0).([]product.Category), args.Error(1)
}

func (m *MockCategoryRepo) UpdateCategory(ctx context.Context, c product.Category) (*product.Category, error) {
	args := m.Called(ctx, c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Category), args.Error(1)
}

func (m *MockCategoryRepo) DeleteCategory(ctx context.Context, id int64, venueID int) error {
	args := m.Called(ctx, id, venueID)
	return args.Error(0)
}

// Tests

const testVenueID = 1

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := uc.NewUsecase(mockRepo, nil)

	ctx := context.Background()
	input := product.Product{Name: "Burger", SalesPrice: 10.5, IsActive: true}
	expected := &product.Product{ID: 1, Name: "Burger", SalesPrice: 10.5, IsActive: true}

	mockRepo.On("CreateProduct", ctx, input).Return(expected, nil)

	result, err := service.CreateProduct(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_ValidationError(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := uc.NewUsecase(mockRepo, nil)

	ctx := context.Background()

	_, err := service.CreateProduct(ctx, product.Product{Name: "", SalesPrice: 10})
	assert.ErrorIs(t, err, product.ErrNameEmpty)

	_, err = service.CreateProduct(ctx, product.Product{Name: "Burger", SalesPrice: -1})
	assert.ErrorIs(t, err, product.ErrPriceNegative)

	mockRepo.AssertNotCalled(t, "CreateProduct")
}

func TestCreateCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	service := uc.NewUsecase(nil, mockRepo)

	ctx := context.Background()
	input := product.Category{Name: "Drinks"}
	expected := &product.Category{ID: 1, Name: "Drinks", CreatedAt: time.Now()}

	mockRepo.On("CreateCategory", ctx, input).Return(expected, nil)

	result, err := service.CreateCategory(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateCategory_ValidationError(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	service := uc.NewUsecase(nil, mockRepo)

	ctx := context.Background()

	_, err := service.CreateCategory(ctx, product.Category{Name: ""})
	assert.ErrorIs(t, err, product.ErrNameEmpty)

	mockRepo.AssertNotCalled(t, "CreateCategory")
}

func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := uc.NewUsecase(mockRepo, nil)
	ctx := context.Background()

	id := int64(1)
	current := &product.Product{ID: id, Name: "Old Name", SalesPrice: 10, IsActive: true}

	mockRepo.On("GetByID", ctx, id, testVenueID).Return(current, nil)

	updatedName := "New Name"
	updatedProduct := &product.Product{ID: id, Name: updatedName, SalesPrice: 10, IsActive: true}
	mockRepo.On("UpdateProduct", ctx, mock.MatchedBy(func(p product.Product) bool {
		return p.Name == updatedName && p.SalesPrice == 10
	})).Return(updatedProduct, nil)

	result, err := service.UpdateProduct(ctx, id, testVenueID, product.Product{
		Name:       updatedName,
		SalesPrice: 10,
		IsActive:   true,
	})

	assert.NoError(t, err)
	assert.Equal(t, updatedName, result.Name)
	mockRepo.AssertExpectations(t)
}
