package ingredient_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AXONcompany/POS/internal/domain/ingredient"
	uc "github.com/AXONcompany/POS/internal/usecase/ingredient"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIngredientRepo struct {
	mock.Mock
}

func (m *MockIngredientRepo) NewIngredient(ctx context.Context, ing ingredient.Ingredient) (*ingredient.Ingredient, error) {
	args := m.Called(ctx, ing)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ingredient.Ingredient), args.Error(1)
}

func (m *MockIngredientRepo) GetByID(ctx context.Context, id int64, venueID int) (*ingredient.Ingredient, error) {
	args := m.Called(ctx, id, venueID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ingredient.Ingredient), args.Error(1)
}

func (m *MockIngredientRepo) GetAllIngredients(ctx context.Context, venueID int, page, pageSize int) ([]ingredient.Ingredient, error) {
	args := m.Called(ctx, venueID, page, pageSize)
	return args.Get(0).([]ingredient.Ingredient), args.Error(1)
}

func (m *MockIngredientRepo) UpdateIngredient(ctx context.Context, ing ingredient.Ingredient) (ingredient.Ingredient, error) {
	args := m.Called(ctx, ing)
	return args.Get(0).(ingredient.Ingredient), args.Error(1)
}

func (m *MockIngredientRepo) DeleteIngredient(ctx context.Context, id int64, venueID int) error {
	args := m.Called(ctx, id, venueID)
	return args.Error(0)
}

func (m *MockIngredientRepo) GetAllInventory(ctx context.Context, venueID int) ([]ingredient.Ingredient, error) {
	args := m.Called(ctx, venueID)
	return args.Get(0).([]ingredient.Ingredient), args.Error(1)
}

// --- CreateIngredient ---

func TestCreateIngredient_EmptyName(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)

	_, err := service.CreateIngredient(context.Background(), ingredient.Ingredient{Name: "", Stock: 10})
	assert.ErrorIs(t, err, ingredient.ErrNameEmpty)
	repo.AssertNotCalled(t, "NewIngredient")
}

func TestCreateIngredient_NegativeStock(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)

	_, err := service.CreateIngredient(context.Background(), ingredient.Ingredient{Name: "Tomate", Stock: -1})
	assert.ErrorIs(t, err, ingredient.ErrNegativeStock)
	repo.AssertNotCalled(t, "NewIngredient")
}

func TestCreateIngredient_ZeroStock_IsValid(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	ing := ingredient.Ingredient{Name: "Sal", VenueID: 1, Stock: 0}
	expected := &ingredient.Ingredient{ID: 1, Name: "Sal", Stock: 0}
	repo.On("NewIngredient", ctx, ing).Return(expected, nil)

	result, err := service.CreateIngredient(ctx, ing)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestCreateIngredient_RepoError_ReturnsAlreadyExists(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	ing := ingredient.Ingredient{Name: "Tomate", Stock: 5}
	repo.On("NewIngredient", ctx, ing).Return(nil, errors.New("unique constraint violated"))

	_, err := service.CreateIngredient(ctx, ing)
	assert.ErrorIs(t, err, ingredient.ErrAlreadyExists)
	repo.AssertExpectations(t)
}

// --- GetIngredient ---

func TestGetIngredient_InvalidID(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)

	tests := []struct{ id int64 }{
		{0},
		{-1},
		{-100},
	}
	for _, tt := range tests {
		_, err := service.GetIngredient(context.Background(), tt.id, 1)
		assert.ErrorIs(t, err, ingredient.ErrInvalidID)
	}
	repo.AssertNotCalled(t, "GetByID")
}

func TestGetIngredient_NotFound_ViaNoRows(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, int64(99), 1).Return(nil, pgx.ErrNoRows)

	_, err := service.GetIngredient(ctx, 99, 1)
	assert.ErrorIs(t, err, ingredient.ErrNotFound)
	repo.AssertExpectations(t)
}

func TestGetIngredient_Success(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	expected := &ingredient.Ingredient{ID: 1, Name: "Cebolla", Stock: 50}
	repo.On("GetByID", ctx, int64(1), 1).Return(expected, nil)

	result, err := service.GetIngredient(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

// --- GetAllIngredients (paginación) ---

func TestGetAllIngredients_Pagination_Defaults(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	// page=0 debe normalizarse a 1, pageSize=0 a 20
	repo.On("GetAllIngredients", ctx, 1, 1, 20).Return([]ingredient.Ingredient{}, nil)

	_, err := service.GetAllIngredients(ctx, 1, 0, 0)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetAllIngredients_PageSize_ExceedsMax(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	// pageSize=200 supera el máximo de 100, debe normalizarse a 20
	repo.On("GetAllIngredients", ctx, 1, 1, 20).Return([]ingredient.Ingredient{}, nil)

	_, err := service.GetAllIngredients(ctx, 1, 1, 200)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetAllIngredients_ValidPagination(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	expected := []ingredient.Ingredient{{ID: 1, Name: "Arroz"}}
	repo.On("GetAllIngredients", ctx, 1, 2, 10).Return(expected, nil)

	result, err := service.GetAllIngredients(ctx, 1, 2, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

// --- UpdateIngredient ---

func TestUpdateIngredient_InvalidID(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)

	_, err := service.UpdateIngredient(context.Background(), 0, 1, ingredient.PartialIngredient{})
	assert.ErrorIs(t, err, ingredient.ErrInvalidID)
	repo.AssertNotCalled(t, "GetByID")
}

func TestUpdateIngredient_EmptyName(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	current := &ingredient.Ingredient{ID: 1, Name: "Tomate", Stock: 10}
	repo.On("GetByID", ctx, int64(1), 1).Return(current, nil)

	emptyName := ""
	_, err := service.UpdateIngredient(ctx, 1, 1, ingredient.PartialIngredient{Name: &emptyName})
	assert.ErrorIs(t, err, ingredient.ErrNameEmpty)
}

func TestUpdateIngredient_NegativeStock(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	current := &ingredient.Ingredient{ID: 1, Name: "Tomate", Stock: 10}
	repo.On("GetByID", ctx, int64(1), 1).Return(current, nil)

	negStock := int64(-5)
	_, err := service.UpdateIngredient(ctx, 1, 1, ingredient.PartialIngredient{Stock: &negStock})
	assert.ErrorIs(t, err, ingredient.ErrNegativeStock)
}

func TestUpdateIngredient_NotFound(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, int64(99), 1).Return(nil, pgx.ErrNoRows)

	_, err := service.UpdateIngredient(ctx, 99, 1, ingredient.PartialIngredient{})
	assert.ErrorIs(t, err, ingredient.ErrNotFound)
}

func TestUpdateIngredient_Success(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	current := &ingredient.Ingredient{ID: 1, Name: "Tomate", Stock: 10, VenueID: 1}
	repo.On("GetByID", ctx, int64(1), 1).Return(current, nil)

	newName := "Tomate Cherry"
	updated := ingredient.Ingredient{ID: 1, Name: newName, Stock: 10, VenueID: 1}
	repo.On("UpdateIngredient", ctx, updated).Return(updated, nil)

	result, err := service.UpdateIngredient(ctx, 1, 1, ingredient.PartialIngredient{Name: &newName})
	assert.NoError(t, err)
	assert.Equal(t, newName, result.Name)
	repo.AssertExpectations(t)
}

// --- DeleteIngredient ---

func TestDeleteIngredient_InvalidID(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)

	err := service.DeleteIngredient(context.Background(), 0, 1)
	assert.ErrorIs(t, err, ingredient.ErrInvalidID)

	err = service.DeleteIngredient(context.Background(), -5, 1)
	assert.ErrorIs(t, err, ingredient.ErrInvalidID)

	repo.AssertNotCalled(t, "DeleteIngredient")
}

func TestDeleteIngredient_Success(t *testing.T) {
	repo := new(MockIngredientRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	repo.On("DeleteIngredient", ctx, int64(1), 1).Return(nil)

	err := service.DeleteIngredient(ctx, 1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
