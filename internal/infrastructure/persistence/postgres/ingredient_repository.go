package postgres

import (
	"context"
	"time"

	"github.com/AXONcompany/POS/internal/domain/ingredient"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type IngredientRepository struct {
	q *sqlc.Queries
}

func ptrTime(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func toDomainIngredient(ing sqlc.Ingredient) ingredient.Ingredient {
	return ingredient.Ingredient{
		ID:             ing.ID,
		Name:           ing.IngredientName,
		UnitOfMeasure:  ing.UnitOfMeasure,
		IngredientType: ing.IngredientType,
		CreatedAt:      ing.CreatedAt.Time,
		Stock:          ing.Stock,
		UpdatedAt:      ptrTime(ing.UpdatedAt),
		DeletedAt:      ptrTime(ing.DeletedAt),
	}
}

func NewIngredientRepository(db *DB) *IngredientRepository {
	return &IngredientRepository{q: sqlc.New(db.Pool)}
}

func (r *IngredientRepository) GetByID(ctx context.Context, id int64) (*ingredient.Ingredient, error) {

	row, err := r.q.GetIngredientByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ing := toDomainIngredient(row)
	return &ing, nil

}

func (r *IngredientRepository) NewIngredient(ctx context.Context, ing ingredient.Ingredient) (*ingredient.Ingredient, error) {
	row, err := r.q.CreateIngredient(ctx, sqlc.CreateIngredientParams{
		IngredientName: ing.Name,
		UnitOfMeasure:  ing.UnitOfMeasure,
		IngredientType: ing.IngredientType,
	})
	if err != nil {
		return nil, err
	}

	created := toDomainIngredient(row)
	return &created, nil
}

func (r *IngredientRepository) DeleteIngredient(ctx context.Context, id int64) error {
	return r.q.DeleteIngredient(ctx, id)
}

func (r *IngredientRepository) GetAllIngredients(ctx context.Context) ([]ingredient.Ingredient, error) {

	rows, err := r.q.ListIngredients(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]ingredient.Ingredient, len(rows))

	for i := range rows {
		items[i] = toDomainIngredient(rows[i])
	}

	return items, nil
}

func (r *IngredientRepository) UpdateIngredient(ctx context.Context, ing ingredient.Ingredient) (ingredient.Ingredient, error) {

	row, err := r.q.UpdateIngredient(
		ctx,
		sqlc.UpdateIngredientParams{
			ID:             ing.ID,
			IngredientName: ing.Name,
			UnitOfMeasure:  ing.UnitOfMeasure,
			IngredientType: ing.IngredientType,
		},
	)

	if err != nil {
		return ingredient.Ingredient{}, err
	}

	return toDomainIngredient(row), nil

}
