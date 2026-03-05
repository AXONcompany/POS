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
	GetByID(ctx context.Context, id int64, venueID int) (*ingredient.Ingredient, error)
	GetAllIngredients(ctx context.Context, venueID int, page, pageSize int) ([]ingredient.Ingredient, error)
	UpdateIngredient(ctx context.Context, ing ingredient.Ingredient) (ingredient.Ingredient, error)
	DeleteIngredient(ctx context.Context, id int64, venueID int) error
	GetAllInventory(ctx context.Context, venueID int) ([]ingredient.Ingredient, error)
}

type Usecase struct {
	repo IngredientRepository
}

func NewUsecase(repo IngredientRepository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (s *Usecase) CreateIngredient(ctx context.Context, ing ingredient.Ingredient) (*ingredient.Ingredient, error) {
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

func (s *Usecase) GetIngredient(ctx context.Context, id int64, venueID int) (*ingredient.Ingredient, error) {
	if id <= 0 {
		return nil, ingredient.ErrInvalidID
	}

	ing, err := s.repo.GetByID(ctx, id, venueID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ingredient.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get ingredient %w", err)
	}

	return ing, nil
}

func (s *Usecase) GetAllIngredients(ctx context.Context, venueID int, page, pageSize int) ([]ingredient.Ingredient, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.GetAllIngredients(ctx, venueID, page, pageSize)
}

func (s *Usecase) UpdateIngredient(ctx context.Context, id int64, venueID int, updates ingredient.PartialIngredient) (*ingredient.Ingredient, error) {
	if id <= 0 {
		return nil, ingredient.ErrInvalidID
	}

	current, err := s.GetIngredient(ctx, id, venueID)
	if err != nil {
		return nil, err
	}

	if updates.Name != nil {
		if *updates.Name == "" {
			return nil, ingredient.ErrNameEmpty
		}
		current.Name = *updates.Name
	}
	if updates.UnitOfMeasure != nil {
		current.UnitOfMeasure = *updates.UnitOfMeasure
	}
	if updates.IngredientType != nil {
		current.IngredientType = *updates.IngredientType
	}
	if updates.Stock != nil {
		if *updates.Stock < 0 {
			return nil, ingredient.ErrNegativeStock
		}
		current.Stock = *updates.Stock
	}

	current.ID = id

	updated, err := s.repo.UpdateIngredient(ctx, *current)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ingredient.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update ingredient: %w", err)
	}

	return &updated, nil
}

func (s *Usecase) DeleteIngredient(ctx context.Context, id int64, venueID int) error {
	if id <= 0 {
		return ingredient.ErrInvalidID
	}

	return s.repo.DeleteIngredient(ctx, id, venueID)
}

func (s *Usecase) GetInventoryReport(ctx context.Context, venueID int) ([]ingredient.Ingredient, error) {
	return s.repo.GetAllInventory(ctx, venueID)
}
