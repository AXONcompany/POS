package postgres

import (
	"context"
	"time"

	"github.com/AXONcompany/POS/internal/domain/table"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
)

type TableAssignmentRepository struct {
	queries *sqlc.Queries
}

func NewTableAssignmentRepository(db *DB) *TableAssignmentRepository {
	return &TableAssignmentRepository{
		queries: sqlc.New(db.Pool),
	}
}

func (r *TableAssignmentRepository) Assign(ctx context.Context, tableID int64, userID, venueID int) (*table.Assignment, error) {
	// Primero cerrar asignacion activa si existe
	_ = r.queries.UnassignWaiterFromTable(ctx, sqlc.UnassignWaiterFromTableParams{
		TableID: tableID,
		VenueID: int32(venueID),
	})

	row, err := r.queries.AssignWaiterToTable(ctx, sqlc.AssignWaiterToTableParams{
		TableID: tableID,
		UserID:  int32(userID),
		VenueID: int32(venueID),
	})
	if err != nil {
		return nil, err
	}

	return r.mapToDomain(row), nil
}

func (r *TableAssignmentRepository) Unassign(ctx context.Context, tableID int64, venueID int) error {
	return r.queries.UnassignWaiterFromTable(ctx, sqlc.UnassignWaiterFromTableParams{
		TableID: tableID,
		VenueID: int32(venueID),
	})
}

func (r *TableAssignmentRepository) GetActive(ctx context.Context, tableID int64, venueID int) (*table.Assignment, error) {
	row, err := r.queries.GetActiveAssignment(ctx, sqlc.GetActiveAssignmentParams{
		TableID: tableID,
		VenueID: int32(venueID),
	})
	if err != nil {
		return nil, err
	}
	return r.mapToDomain(row), nil
}

func (r *TableAssignmentRepository) ListByTable(ctx context.Context, tableID int64, venueID int) ([]table.AssignmentDetail, error) {
	rows, err := r.queries.ListAssignmentsByTable(ctx, sqlc.ListAssignmentsByTableParams{
		TableID: tableID,
		VenueID: int32(venueID),
	})
	if err != nil {
		return nil, err
	}

	result := make([]table.AssignmentDetail, len(rows))
	for i, row := range rows {
		result[i] = table.AssignmentDetail{
			ID:         row.ID,
			TableID:    row.TableID,
			UserID:     int(row.UserID),
			VenueID:    int(row.VenueID),
			WaiterName: row.WaiterName,
			AssignedAt: row.AssignedAt.Time,
		}
		if row.UnassignedAt.Valid {
			t := row.UnassignedAt.Time
			result[i].UnassignedAt = &t
		}
	}
	return result, nil
}

func (r *TableAssignmentRepository) mapToDomain(row sqlc.TableAssignment) *table.Assignment {
	a := &table.Assignment{
		ID:         row.ID,
		TableID:    row.TableID,
		UserID:     int(row.UserID),
		VenueID:    int(row.VenueID),
		AssignedAt: row.AssignedAt.Time,
	}
	if row.UnassignedAt.Valid {
		t := row.UnassignedAt.Time
		a.UnassignedAt = &t
	}
	return a
}

// Compile-time check
var _ interface {
	Assign(ctx context.Context, tableID int64, userID, venueID int) (*table.Assignment, error)
	Unassign(ctx context.Context, tableID int64, venueID int) error
	GetActive(ctx context.Context, tableID int64, venueID int) (*table.Assignment, error)
	ListByTable(ctx context.Context, tableID int64, venueID int) ([]table.AssignmentDetail, error)
} = (*TableAssignmentRepository)(nil)

// Ensure time import is used
var _ = time.Now
