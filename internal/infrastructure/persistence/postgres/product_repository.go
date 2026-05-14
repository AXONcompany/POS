package postgres

import (
	"context"
	"fmt"

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
	dp := product.Product{
		ID:         p.ID,
		VenueID:    int(p.VenueID),
		Name:       p.ProductName,
		SalesPrice: val.Float64,
		IsActive:   p.IsActive,
		CreatedAt:  p.CreatedAt.Time,
		UpdatedAt:  ptrTime(p.UpdatedAt),
		DeletedAt:  ptrTime(p.DeletedAt),
	}
	if p.Description.Valid {
		dp.Description = p.Description.String
	}
	if p.ImageUrl.Valid {
		dp.ImageURL = p.ImageUrl.String
	}
	if p.CategoryID.Valid {
		v := p.CategoryID.Int64
		dp.CategoryID = &v
	}
	return dp
}

func floatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(fmt.Sprintf("%f", f))
	return n
}

func toOptionalText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func toOptionalInt8(p *int64) pgtype.Int8 {
	if p == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *p, Valid: true}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p product.Product) (*product.Product, error) {
	row, err := r.q.CreateProduct(ctx, sqlc.CreateProductParams{
		VenueID:     int32(p.VenueID),
		ProductName: p.Name,
		SalesPrice:  floatToNumeric(p.SalesPrice),
		IsActive:    p.IsActive,
		Description: toOptionalText(p.Description),
		ImageUrl:    toOptionalText(p.ImageURL),
		CategoryID:  toOptionalInt8(p.CategoryID),
	})
	if err != nil {
		return nil, err
	}
	dp := toDomainProduct(row)
	return &dp, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64, venueID int) (*product.Product, error) {
	row, err := r.q.GetProduct(ctx, sqlc.GetProductParams{
		ID:      id,
		VenueID: int32(venueID),
	})
	if err != nil {
		return nil, err
	}
	dp := toDomainProduct(row)
	return &dp, nil
}

func (r *ProductRepository) GetAllProducts(ctx context.Context, venueID int, page, pageSize int) ([]product.Product, error) {
	offset := (page - 1) * pageSize
	rows, err := r.q.ListProducts(ctx, sqlc.ListProductsParams{
		VenueID: int32(venueID),
		Limit:   int32(pageSize),
		Offset:  int32(offset),
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
		VenueID:     int32(p.VenueID),
		ProductName: p.Name,
		SalesPrice:  floatToNumeric(p.SalesPrice),
		IsActive:    p.IsActive,
		Description: toOptionalText(p.Description),
		ImageUrl:    toOptionalText(p.ImageURL),
		CategoryID:  toOptionalInt8(p.CategoryID),
	})
	if err != nil {
		return nil, err
	}
	dp := toDomainProduct(row)
	return &dp, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int64, venueID int) error {
	return r.q.DeleteProduct(ctx, sqlc.DeleteProductParams{
		ID:      id,
		VenueID: int32(venueID),
	})
}

func (r *ProductRepository) GetProductPrice(ctx context.Context, productID int64, venueID int) (float64, error) {
	row, err := r.q.GetProduct(ctx, sqlc.GetProductParams{
		ID:      productID,
		VenueID: int32(venueID),
	})
	if err != nil {
		return 0, fmt.Errorf("get product price: %w", err)
	}
	val, _ := row.SalesPrice.Float64Value()
	return val.Float64, nil
}
