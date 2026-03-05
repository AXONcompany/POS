package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/product"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
)

type CategoryRepository struct {
	q *sqlc.Queries
}

func NewCategoryRepository(db *DB) *CategoryRepository {
	return &CategoryRepository{q: sqlc.New(db.Pool)}
}

func toDomainCategory(c sqlc.Category) product.Category {
	return product.Category{
		ID:        c.ID,
		VenueID:   int(c.VenueID),
		Name:      c.CategoryName,
		CreatedAt: c.CreatedAt.Time,
		UpdatedAt: ptrTime(c.UpdatedAt),
		DeletedAt: ptrTime(c.DeletedAt),
	}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, c product.Category) (*product.Category, error) {
	row, err := r.q.CreateCategory(ctx, sqlc.CreateCategoryParams{
		VenueID:      int32(c.VenueID),
		CategoryName: c.Name,
	})
	if err != nil {
		return nil, err
	}
	pc := toDomainCategory(row)
	return &pc, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64, venueID int) (*product.Category, error) {
	row, err := r.q.GetCategory(ctx, sqlc.GetCategoryParams{
		ID:      id,
		VenueID: int32(venueID),
	})
	if err != nil {
		return nil, err
	}
	pc := toDomainCategory(row)
	return &pc, nil
}

func (r *CategoryRepository) GetAllCategories(ctx context.Context, venueID int, page, pageSize int) ([]product.Category, error) {
	offset := (page - 1) * pageSize
	rows, err := r.q.ListCategories(ctx, sqlc.ListCategoriesParams{
		VenueID: int32(venueID),
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	items := make([]product.Category, len(rows))
	for i, row := range rows {
		items[i] = toDomainCategory(row)
	}
	return items, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, c product.Category) (*product.Category, error) {
	row, err := r.q.UpdateCategory(ctx, sqlc.UpdateCategoryParams{
		ID:           c.ID,
		VenueID:      int32(c.VenueID),
		CategoryName: c.Name,
	})
	if err != nil {
		return nil, err
	}
	pc := toDomainCategory(row)
	return &pc, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int64, venueID int) error {
	return r.q.DeleteCategory(ctx, sqlc.DeleteCategoryParams{
		ID:      id,
		VenueID: int32(venueID),
	})
}
