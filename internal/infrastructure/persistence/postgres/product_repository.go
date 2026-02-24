package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/product"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductRepository struct {
	q  *sqlc.Queries
	db *DB
}

func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{
		q:  sqlc.New(db.Pool),
		db: db,
	}
}

func toDomainProduct(p sqlc.Product) product.Product {
	val, _ := p.SalesPrice.Float64Value()
	return product.Product{
		ID:         p.ID,
		Name:       p.ProductName,
		SalesPrice: val.Float64,
		IsActive:   p.IsActive,
		CreatedAt:  p.CreatedAt.Time,
		UpdatedAt:  ptrTime(p.UpdatedAt),
		DeletedAt:  ptrTime(p.DeletedAt),
	}
}

func floatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(f)
	return n
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p product.Product) (*product.Product, error) {
	row, err := r.q.CreateProduct(ctx, sqlc.CreateProductParams{
		ProductName: p.Name,
		SalesPrice:  floatToNumeric(p.SalesPrice),
		IsActive:    p.IsActive,
	})
	if err != nil {
		return nil, err
	}
	dp := toDomainProduct(row)
	return &dp, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*product.Product, error) {
	row, err := r.q.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	dp := toDomainProduct(row)
	return &dp, nil
}

func (r *ProductRepository) GetAllProducts(ctx context.Context, page, pageSize int) ([]product.Product, error) {
	offset := (page - 1) * pageSize
	rows, err := r.q.ListProducts(ctx, sqlc.ListProductsParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	items := make([]product.Product, len(rows))
	for i, row := range rows {
		items[i] = toDomainProduct(row)
	}
	return items, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, p product.Product) (*product.Product, error) {
	row, err := r.q.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:          p.ID,
		ProductName: p.Name,
		SalesPrice:  floatToNumeric(p.SalesPrice),
		IsActive:    p.IsActive,
	})
	if err != nil {
		return nil, err
	}
	dp := toDomainProduct(row)
	return &dp, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	return r.q.DeleteProduct(ctx, id)
}

func (r *ProductRepository) CreateProductWithRecipe(ctx context.Context, p product.Product, items []product.RecipeItem) (*product.Product, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	// 1. Create Product
	prodRow, err := qtx.CreateProduct(ctx, sqlc.CreateProductParams{
		ProductName: p.Name,
		SalesPrice:  floatToNumeric(p.SalesPrice),
		IsActive:    p.IsActive,
	})
	if err != nil {
		return nil, err
	}

	// 2. Create Recipe Items
	for _, item := range items {
		_, err := qtx.AddRecipeItem(ctx, sqlc.AddRecipeItemParams{
			ProductID:        prodRow.ID,
			IngredientID:     item.IngredientID,
			QuantityRequired: floatToNumeric(item.QuantityRequired),
		})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	res := toDomainProduct(prodRow)
	return &res, nil
}
