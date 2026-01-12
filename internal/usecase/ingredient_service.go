package usecase

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/ingredient"
)

type IngredientRepository interface {
	NewIngredient(ctx context.Context, ing ingredient.Ingredient) (*ingredient.Ingredient, error)
	GetByID(ctx context.Context, id int64) (*ingredient.Ingredient, error)
	GetAllIngredients(ctx context.Context) ([]ingredient.Ingredient, error)
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
