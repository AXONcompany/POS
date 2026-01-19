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

func NewTableRepository(db sqlc.DBTX) *TableRepository {
	return &TableRepository{
		queries: sqlc.New(db),
	}
}

func (r *TableRepository) Create(ctx context.Context, tbl *table.Table) error {
	params := sqlc.CreateTableParams{
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

	tbl.ID = generated.ID
	tbl.CreatedAt = generated.CreatedAt.Time
	return nil
}

func (r *TableRepository) FindAll(ctx context.Context) ([]table.Table, error) {
	rows, err := r.queries.ListTables(ctx)
	if err != nil {
		return nil, err
	}

	tables := make([]table.Table, len(rows))
	for i, row := range rows {
		tables[i] = r.mapToDomain(row)
	}
	return tables, nil
}

func (r *TableRepository) FindByID(ctx context.Context, id int64) (*table.Table, error) {
	row, err := r.queries.GetTable(ctx, id)
	if err != nil {
		return nil, err
	}

	t := r.mapToDomain(row)
	return &t, nil
}
func (r *TableRepository) Update(ctx context.Context, id int64, updates *table.TableUpdates) error {
	params := sqlc.UpdateTableParams{
		ID: id,
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

func (r *TableRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteTable(ctx, id)
}

func (r *TableRepository) AssignWaitressToTable(ctx context.Context, tableID int64, waitressID int64) error {
	params := sqlc.AssignWaitressToTableParams{
		TableID:    tableID,
		WaitressID: waitressID,
	}
	_, err := r.queries.AssignWaitressToTable(ctx, params)
	return err
}

func (r *TableRepository) RemoveWaitressFromTable(ctx context.Context, tableID int64, waitressID int64) error {
	return r.queries.RemoveWaitressFromTable(ctx, tableID)
}

func (r *TableRepository) FindWaitressesByTableID(ctx context.Context, tableID int64) ([]table.TableWaitress, error) {

	row, err := r.queries.GetWaitressByTable(ctx, tableID)
	if err != nil {
		return nil, err
	}

	return []table.TableWaitress{
		{
			ID:         row.ID,
			TableID:    row.TableID,
			WaitressID: row.WaitressID,
			CreatedAt:  row.CreatedAt.Time,
		},
	}, nil
}
func (r *TableRepository) mapToDomain(row sqlc.Table) table.Table {
	t := table.Table{
		ID:        row.ID,
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
