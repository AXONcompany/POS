package postgres

import (
	"context"
	
	domainSale "github.com/AXONcompany/POS/internal/domain/sale"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type SaleRepository struct {
	q  *sqlc.Queries
}


func NewSaleRepository(db *DB) *SaleRepository {
	return &SaleRepository{q: sqlc.New(db.Pool)}
}

func toDomainSale(s sqlc.CreateSaleRow) *domainSale.Sale {
	total, _ := s.Total.Float64Value()

	return &domainSale.Sale{
		ID:            s.ID,
		Total:         total.Float64,
		PaymentMethod: s.PaymentMethod,
		Date:          s.Date.Time,
		OrderID:       int64(s.OrderID),
		CreatedAt:     s.CreatedAt.Time,
		UpdatedAt:     ptrTime(s.UpdatedAt),
		DeletedAt:     ptrTime(s.DeletedAt),
	}
}

func toDomainSaleFromGet(s sqlc.GetSaleByIDRow) *domainSale.Sale {
	total, _ := s.Total.Float64Value()

	return &domainSale.Sale{
		ID:            s.ID,
		Total:         total.Float64,
		PaymentMethod: s.PaymentMethod,
		Date:          s.Date.Time,
		OrderID:       int64(s.OrderID),
		CreatedAt:     s.CreatedAt.Time,
		UpdatedAt:     ptrTime(s.UpdatedAt),
		DeletedAt:     ptrTime(s.DeletedAt),
	}
}

func (r *SaleRepository) CreateSale(ctx context.Context, s domainSale.Sale) (*domainSale.Sale, error) {
	var total pgtype.Numeric
	total.Scan(s.Total)

	row, err := r.q.CreateSale(ctx, sqlc.CreateSaleParams{
		Total:         total,
		PaymentMethod: s.PaymentMethod,
		Date:          pgtype.Timestamptz{Time: s.Date, Valid: true},
		OrderID:       int32(s.OrderID),
	})
	if err != nil {
		return nil, err
	}

	created := toDomainSale(row)
	return created, nil
}

func (r *SaleRepository) GetByID(ctx context.Context, id int64) (*domainSale.Sale, error) {
	row, err := r.q.GetSaleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sale := toDomainSaleFromGet(row)
	return sale, nil
}