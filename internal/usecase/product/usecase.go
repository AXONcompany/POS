package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/AXONcompany/POS/internal/domain/product"
	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, p product.Product) (*product.Product, error)
	GetByID(ctx context.Context, id int64) (*product.Product, error)
	GetAllProducts(ctx context.Context, page, pageSize int) ([]product.Product, error)
	UpdateProduct(ctx context.Context, p product.Product) (*product.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	CreateProductWithRecipe(ctx context.Context, p product.Product, items []product.RecipeItem) (*product.Product, error)
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, c product.Category) (*product.Category, error)
	GetByID(ctx context.Context, id int64) (*product.Category, error)
	GetAllCategories(ctx context.Context, page, pageSize int) ([]product.Category, error)
	UpdateCategory(ctx context.Context, c product.Category) (*product.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
}

type RecipeRepository interface {
	AddRecipeItem(ctx context.Context, item product.RecipeItem) (*product.RecipeItem, error)
	GetByProductID(ctx context.Context, productID int64) ([]product.RecipeItem, error)
}

type Usecase struct {
	productRepo  ProductRepository
	categoryRepo CategoryRepository
	recipeRepo   RecipeRepository
}

func NewUsecase(productRepo ProductRepository, categoryRepo CategoryRepository, recipeRepo RecipeRepository) *Usecase {
	return &Usecase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		recipeRepo:   recipeRepo,
	}
}

// Category Methods

func (s *Usecase) CreateCategory(ctx context.Context, c product.Category) (*product.Category, error) {
	if c.Name == "" {
		return nil, product.ErrNameEmpty
	}

	created, err := s.categoryRepo.CreateCategory(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return created, nil
}

func (s *Usecase) GetCategory(ctx context.Context, id int64) (*product.Category, error) {
	if id <= 0 {
		return nil, product.ErrInvalidID
	}
	c, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrCategoryNotFound
		}
		return nil, err
	}
	return c, nil
}

func (s *Usecase) GetAllCategories(ctx context.Context, page, pageSize int) ([]product.Category, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.categoryRepo.GetAllCategories(ctx, page, pageSize)
}

func (s *Usecase) UpdateCategory(ctx context.Context, id int64, name string) (*product.Category, error) {
	if id <= 0 {
		return nil, product.ErrInvalidID
	}
	if name == "" {
		return nil, product.ErrNameEmpty
	}

	current, err := s.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	current.Name = name
	return s.categoryRepo.UpdateCategory(ctx, *current)
}

func (s *Usecase) DeleteCategory(ctx context.Context, id int64) error {
	if id <= 0 {
		return product.ErrInvalidID
	}
	return s.categoryRepo.DeleteCategory(ctx, id)
}

// Product Methods

func (s *Usecase) CreateProduct(ctx context.Context, p product.Product) (*product.Product, error) {
	if p.Name == "" {
		return nil, product.ErrNameEmpty
	}
	if p.SalesPrice < 0 {
		return nil, product.ErrPriceNegative
	}

	created, err := s.productRepo.CreateProduct(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return created, nil
}

func (s *Usecase) GetProduct(ctx context.Context, id int64) (*product.Product, error) {
	if id <= 0 {
		return nil, product.ErrInvalidID
	}
	p, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *Usecase) GetAllProducts(ctx context.Context, page, pageSize int) ([]product.Product, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.productRepo.GetAllProducts(ctx, page, pageSize)
}

func (s *Usecase) UpdateProduct(ctx context.Context, id int64, p product.Product) (*product.Product, error) {
	if id <= 0 {
		return nil, product.ErrInvalidID
	}

	current, err := s.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates (simplified logic here, usually pass DTO or Partial)
	if p.Name != "" {
		current.Name = p.Name
	}
	if p.SalesPrice >= 0 {
		current.SalesPrice = p.SalesPrice
	}
	current.IsActive = p.IsActive

	return s.productRepo.UpdateProduct(ctx, *current)
}

func (s *Usecase) DeleteProduct(ctx context.Context, id int64) error {
	if id <= 0 {
		return product.ErrInvalidID
	}
	return s.productRepo.DeleteProduct(ctx, id)
}

// Recipe Methods

func (s *Usecase) AddIngredient(ctx context.Context, productID, ingredientID int64, quantity float64) (*product.RecipeItem, error) {
	if productID <= 0 || ingredientID <= 0 {
		return nil, product.ErrInvalidID
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}

	item := product.RecipeItem{
		ProductID:        productID,
		IngredientID:     ingredientID,
		QuantityRequired: quantity,
	}

	return s.recipeRepo.AddRecipeItem(ctx, item)
}

func (s *Usecase) GetProductIngredients(ctx context.Context, productID int64) ([]product.RecipeItem, error) {
	if productID <= 0 {
		return nil, product.ErrInvalidID
	}
	return s.recipeRepo.GetByProductID(ctx, productID)
}

// Menu Methods

func (s *Usecase) CreateMenuItem(ctx context.Context, name string, price float64, ingredients []product.RecipeItem) (*product.Product, error) {
	if name == "" {
		return nil, product.ErrNameEmpty
	}
	if price < 0 {
		return nil, product.ErrPriceNegative
	}
	if len(ingredients) == 0 {
		return nil, errors.New("menu item must have at least one ingredient")
	}

	prod := product.Product{
		Name:       name,
		SalesPrice: price,
		IsActive:   true,
	}

	return s.productRepo.CreateProductWithRecipe(ctx, prod, ingredients)
}
