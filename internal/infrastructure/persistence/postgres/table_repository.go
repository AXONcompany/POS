package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/table"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type TableRepository struct {
	queries *sqlc.Queries
}

func NewTableRepository(db *DB) *TableRepository {
	return &TableRepository{
		queries: sqlc.New(db.Pool),
	}
}

func (r *TableRepository) Create(ctx context.Context, tbl *table.Table) error {
	params := sqlc.CreateTableParams{
		VenueID:     int32(tbl.VenueID),
		TableNumber: int32(tbl.Number),
		Capacity:    int32(tbl.Capacity),
		Status:      tbl.Status,
	}

	if tbl.ArrivalTime != nil {
		params.ArrivalTime = pgtype.Timestamptz{
			Time:  *tbl.ArrivalTime,
			Valid: true,
		}
	} else {
		params.ArrivalTime = pgtype.Timestamptz{
			Valid: false,
		}
	}

	generated, err := r.queries.CreateTable(ctx, params)
	if err != nil {
		return err
	}

	tbl.ID = generated.IDTable
	tbl.CreatedAt = generated.CreatedAt.Time
	return nil
}

func (r *TableRepository) FindAll(ctx context.Context, venueID int) ([]table.Table, error) {
	rows, err := r.queries.ListTables(ctx, int32(venueID))
	if err != nil {
		return nil, err
	}

	tables := make([]table.Table, len(rows))
	for i, row := range rows {
		tables[i] = r.mapToDomain(row)
	}
	return tables, nil
}

func (r *TableRepository) FindByID(ctx context.Context, id int64, venueID int) (*table.Table, error) {
	row, err := r.queries.GetTable(ctx, sqlc.GetTableParams{
		IDTable: id,
		VenueID: int32(venueID),
	})
	if err != nil {
		return nil, err
	}

	t := r.mapToDomain(row)
	return &t, nil
}

func (r *TableRepository) Update(ctx context.Context, id int64, venueID int, updates *table.TableUpdates) error {
	params := sqlc.UpdateTableParams{
		IDTable: id,
		VenueID: int32(venueID),
	}

	if updates.Number != nil {
		params.TableNumber = pgtype.Int4{
			Int32: int32(*updates.Number),
			Valid: true,
		}
	}

	if updates.Capacity != nil {
		params.Capacity = pgtype.Int4{
			Int32: int32(*updates.Capacity),
			Valid: true,
		}
	}

	if updates.Status != nil {
		params.Status = pgtype.Text{
			String: *updates.Status,
			Valid:  true,
		}
	}

	if updates.ArrivalTime != nil {
		params.ArrivalTime = pgtype.Timestamptz{
			Time:  *updates.ArrivalTime,
			Valid: true,
		}
	}
	return r.queries.UpdateTable(ctx, params)
}

func (r *TableRepository) Delete(ctx context.Context, id int64, venueID int) error {
	return r.queries.DeleteTable(ctx, sqlc.DeleteTableParams{
		IDTable: id,
		VenueID: int32(venueID),
	})
}

func (r *TableRepository) mapToDomain(row sqlc.Table) table.Table {
	t := table.Table{
		ID:        row.IDTable,
		VenueID:   int(row.VenueID),
		Number:    int(row.TableNumber),
		Capacity:  int(row.Capacity),
		Status:    row.Status,
		CreatedAt: row.CreatedAt.Time,
	}
	if row.ArrivalTime.Valid {
		t.ArrivalTime = &row.ArrivalTime.Time
	}
	return t
}
