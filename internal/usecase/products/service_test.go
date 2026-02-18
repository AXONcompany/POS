package products

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AXONcompany/POS/internal/domain/product"
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

func (m *MockProductRepo) GetByID(ctx context.Context, id int64) (*product.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductRepo) GetAllProducts(ctx context.Context, page, pageSize int) ([]product.Product, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]product.Product), args.Error(1)
}

func (m *MockProductRepo) UpdateProduct(ctx context.Context, p product.Product) (*product.Product, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockProductRepo) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepo) CreateProductWithRecipe(ctx context.Context, p product.Product, items []product.RecipeItem) (*product.Product, error) {
	args := m.Called(ctx, p, items)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
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

func (m *MockCategoryRepo) GetByID(ctx context.Context, id int64) (*product.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Category), args.Error(1)
}

func (m *MockCategoryRepo) GetAllCategories(ctx context.Context, page, pageSize int) ([]product.Category, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]product.Category), args.Error(1)
}

func (m *MockCategoryRepo) UpdateCategory(ctx context.Context, c product.Category) (*product.Category, error) {
	args := m.Called(ctx, c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Category), args.Error(1)
}

func (m *MockCategoryRepo) DeleteCategory(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockRecipeRepo struct {
	mock.Mock
}

func (m *MockRecipeRepo) AddRecipeItem(ctx context.Context, item product.RecipeItem) (*product.RecipeItem, error) {
	args := m.Called(ctx, item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.RecipeItem), args.Error(1)
}

func (m *MockRecipeRepo) GetByProductID(ctx context.Context, productID int64) ([]product.RecipeItem, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]product.RecipeItem), args.Error(1)
}

// Tests

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := NewService(mockRepo, nil, nil)

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
	service := NewService(mockRepo, nil, nil)

	ctx := context.Background()

	// Empty Name
	_, err := service.CreateProduct(ctx, product.Product{Name: "", SalesPrice: 10})
	assert.ErrorIs(t, err, product.ErrNameEmpty)

	// Negative Price
	_, err = service.CreateProduct(ctx, product.Product{Name: "Burger", SalesPrice: -1})
	assert.ErrorIs(t, err, product.ErrPriceNegative)

	mockRepo.AssertNotCalled(t, "CreateProduct")
}

func TestCreateCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepo)
	service := NewService(nil, mockRepo, nil)

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
	service := NewService(nil, mockRepo, nil)

	ctx := context.Background()

	_, err := service.CreateCategory(ctx, product.Category{Name: ""})
	assert.ErrorIs(t, err, product.ErrNameEmpty)

	mockRepo.AssertNotCalled(t, "CreateCategory")
}

func TestAddIngredient_Success(t *testing.T) {
	mockProdRepo := new(MockProductRepo)
	mockRecipeRepo := new(MockRecipeRepo)
	service := NewService(mockProdRepo, nil, mockRecipeRepo)

	ctx := context.Background()
	productID := int64(1)
	ingredientID := int64(10)
	quantity := 2.5

	existingProduct := &product.Product{ID: productID, Name: "Burger"}
	mockProdRepo.On("GetByID", ctx, productID).Return(existingProduct, nil)

	expectedItem := &product.RecipeItem{
		ID:               100,
		ProductID:        productID,
		IngredientID:     ingredientID,
		QuantityRequired: quantity,
	}
	mockRecipeRepo.On("AddRecipeItem", ctx, product.RecipeItem{
		ProductID:        productID,
		IngredientID:     ingredientID,
		QuantityRequired: quantity,
	}).Return(expectedItem, nil)

	result, err := service.AddIngredient(ctx, productID, ingredientID, quantity)

	assert.NoError(t, err)
	assert.Equal(t, expectedItem, result)
}

func TestAddIngredient_Cases(t *testing.T) {
	mockProdRepo := new(MockProductRepo)
	mockRecipeRepo := new(MockRecipeRepo)
	service := NewService(mockProdRepo, nil, mockRecipeRepo)

	ctx := context.Background()

	// Invalid IDs
	_, err := service.AddIngredient(ctx, 0, 10, 1.0)
	assert.ErrorIs(t, err, product.ErrInvalidID)

	// Negative Quantity
	_, err = service.AddIngredient(ctx, 1, 10, -5.0)
	assert.ErrorContains(t, err, "quantity must be positive")

	// Product Not Found
	mockProdRepo.On("GetByID", ctx, int64(99)).Return(nil, errors.New("not found"))
	// Note: In real implementation we check specifically for pgx.ErrNoRows or allow repository to return a clean error.
	// We'll simulate a generic error here as the service doesn't wrap "not found" unless it matches pgx.ErrNoRows.
	// For this test, let's assume we want to ensure it fails if repo fails.

	_, err = service.AddIngredient(ctx, 99, 10, 1.0)
	assert.Error(t, err)
}

func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := NewService(mockRepo, nil, nil)
	ctx := context.Background()

	id := int64(1)
	current := &product.Product{ID: id, Name: "Old Name", SalesPrice: 10, IsActive: true}

	// Mock GetByID
	mockRepo.On("GetByID", ctx, id).Return(current, nil)

	// Mock UpdateProduct
	updatedName := "New Name"
	updatedProduct := &product.Product{ID: id, Name: updatedName, SalesPrice: 10, IsActive: true}
	mockRepo.On("UpdateProduct", ctx, mock.MatchedBy(func(p product.Product) bool {
		return p.Name == updatedName && p.SalesPrice == 10
	})).Return(updatedProduct, nil)

	// Call
	result, err := service.UpdateProduct(ctx, id, product.Product{
		Name:       updatedName,
		SalesPrice: 10,
		IsActive:   true,
	})

	assert.NoError(t, err)
	assert.Equal(t, updatedName, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCreateMenuItem_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := NewService(mockRepo, nil, nil)
	ctx := context.Background()

	name := "Special Burger"
	price := 15.0
	ingredients := []product.RecipeItem{
		{IngredientID: 1, QuantityRequired: 2},
		{IngredientID: 2, QuantityRequired: 1},
	}

	expected := &product.Product{ID: 1, Name: name, SalesPrice: price, IsActive: true}

	mockRepo.On("CreateProductWithRecipe", ctx, product.Product{
		Name:       name,
		SalesPrice: price,
		IsActive:   true,
	}, ingredients).Return(expected, nil)

	result, err := service.CreateMenuItem(ctx, name, price, ingredients)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateMenuItem_Validation(t *testing.T) {
	mockRepo := new(MockProductRepo)
	service := NewService(mockRepo, nil, nil)
	ctx := context.Background()

	// Empty Name
	_, err := service.CreateMenuItem(ctx, "", 10, []product.RecipeItem{{IngredientID: 1}})
	assert.ErrorIs(t, err, product.ErrNameEmpty)

	// Negative Price
	_, err = service.CreateMenuItem(ctx, "Burger", -5, []product.RecipeItem{{IngredientID: 1}})
	assert.ErrorIs(t, err, product.ErrPriceNegative)

	// No Ingredients
	_, err = service.CreateMenuItem(ctx, "Burger", 10, nil)
	assert.ErrorContains(t, err, "menu item must have at least one ingredient")

	mockRepo.AssertNotCalled(t, "CreateProductWithRecipe")
}
