package ingredient

import (
	"context"
	"errors"
	"fmt"

	"github.com/AXONcompany/POS/internal/domain/ingredient"
	"github.com/jackc/pgx/v5"
)

type IngredientRepository interface {
	NewIngredient(ctx context.Context, ing ingredient.Ingredient) (*ingredient.Ingredient, error)
	GetByID(ctx context.Context, id int64) (*ingredient.Ingredient, error)
	GetAllIngredients(ctx context.Context, page, pageSize int) ([]ingredient.Ingredient, error)
	UpdateIngredient(ctx context.Context, ing ingredient.Ingredient) (ingredient.Ingredient, error)
	DeleteIngredient(ctx context.Context, id int64) error
}

type IngredientService struct {
	repo IngredientRepository
}

func NewIngredientService(repo IngredientRepository) *IngredientService {
	return &IngredientService{
		repo: repo,
	}
}

func (s *IngredientService) CreateIngredient(ctx context.Context, ing ingredient.Ingredient) (*ingredient.Ingredient, error) {

	if len(ing.Name) == 0 || ing.Name == "" {
		return nil, ingredient.ErrNameEmpty
	}

	if ing.Stock < 0 {
		return nil, ingredient.ErrNegativeStock
	}

	ingr, err := s.repo.NewIngredient(ctx, ing)
	if err != nil {
		return nil, ingredient.ErrAlreadyExists
	}

	return ingr, nil

}

func (s *IngredientService) GetIngredient(ctxt context.Context, id int64) (*ingredient.Ingredient, error) {

	if id <= 0 {
		return nil, ingredient.ErrInvalidID
	}

	ing, err := s.repo.GetByID(ctxt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ingredient.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get ingredient %w", err)
	}

	return ing, nil

}

func (s *IngredientService) GetAllIngredients(ctxt context.Context, page, pageSize int) ([]ingredient.Ingredient, error) {

	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.GetAllIngredients(ctxt, page, pageSize)
}
