package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/product"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
)

type RecipeRepository struct {
	q *sqlc.Queries
}

func NewRecipeRepository(db *DB) *RecipeRepository {
	return &RecipeRepository{q: sqlc.New(db.Pool)}
}

func (r *RecipeRepository) AddRecipeItem(ctx context.Context, item product.RecipeItem) (*product.RecipeItem, error) {
	row, err := r.q.AddRecipeItem(ctx, sqlc.AddRecipeItemParams{
		ProductID:        item.ProductID,
		IngredientID:     item.IngredientID,
		QuantityRequired: floatToNumeric(item.QuantityRequired),
	})
	if err != nil {
		return nil, err
	}

	val, _ := row.QuantityRequired.Float64Value()

	return &product.RecipeItem{
		ID:               row.ID,
		ProductID:        row.ProductID,
		IngredientID:     row.IngredientID,
		QuantityRequired: val.Float64,
	}, nil
}

func (r *RecipeRepository) GetByProductID(ctx context.Context, productID int64) ([]product.RecipeItem, error) {
	rows, err := r.q.GetRecipeByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	items := make([]product.RecipeItem, len(rows))
	for i, row := range rows {
		val, _ := row.QuantityRequired.Float64Value()
		items[i] = product.RecipeItem{
			ID:               row.ID,
			ProductID:        row.ProductID,
			IngredientID:     row.IngredientID,
			IngredientName:   row.IngredientName,
			UnitOfMeasure:    row.UnitOfMeasure,
			QuantityRequired: val.Float64,
		}
	}
	return items, nil
}
