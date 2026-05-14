package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/AXONcompany/POS/internal/domain/table"
)

type TableRepository struct {
	db *DB
}

func NewTableRepository(db *DB) *TableRepository {
	return &TableRepository{db: db}
}

const tableColumns = `
	id_table, venue_id, created_at, updated_at, deleted_at,
	table_number, capacity, status, arrival_time,
	table_name, x, y, width, height, shape, rotation, color, floor,
	is_merged, merged_from, guests, assigned_waiter_id`

func (r *TableRepository) scanTable(row pgx.Row) (table.Table, error) {
	var t table.Table
	var deletedAt pgtype.Timestamptz
	var updatedAt pgtype.Timestamptz
	var arrivalTime pgtype.Timestamptz
	var tableName pgtype.Text
	var color pgtype.Text
	var assignedWaiterID pgtype.Int4
	var mergedFromJSON []byte

	err := row.Scan(
		&t.ID, &t.VenueID, &t.CreatedAt, &updatedAt, &deletedAt,
		&t.Number, &t.Capacity, &t.Status, &arrivalTime,
		&tableName, &t.X, &t.Y, &t.Width, &t.Height, &t.Shape, &t.Rotation, &color, &t.Floor,
		&t.IsMerged, &mergedFromJSON, &t.Guests, &assignedWaiterID,
	)
	if err != nil {
		return table.Table{}, err
	}

	if updatedAt.Valid {
		v := updatedAt.Time
		t.UpdatedAt = &v
	}
	if arrivalTime.Valid {
		v := arrivalTime.Time
		t.ArrivalTime = &v
	}
	if tableName.Valid {
		t.Name = tableName.String
	}
	if color.Valid {
		v := color.String
		t.Color = &v
	}
	if assignedWaiterID.Valid {
		v := int(assignedWaiterID.Int32)
		t.AssignedWaiterID = &v
	}
	if mergedFromJSON != nil {
		json.Unmarshal(mergedFromJSON, &t.MergedFrom)
	}
	return t, nil
}

func (r *TableRepository) Create(ctx context.Context, t *table.Table) error {
	if t.Width == 0 {
		t.Width = 110
	}
	if t.Height == 0 {
		t.Height = 110
	}
	if t.Shape == "" {
		t.Shape = "square"
	}
	if t.Floor == 0 {
		t.Floor = 1
	}

	var mergedFromJSON []byte
	if len(t.MergedFrom) > 0 {
		mergedFromJSON, _ = json.Marshal(t.MergedFrom)
	}

	var colorParam interface{} = nil
	if t.Color != nil {
		colorParam = *t.Color
	}
	var waiterParam interface{} = nil
	if t.AssignedWaiterID != nil {
		waiterParam = *t.AssignedWaiterID
	}

	query := `
		INSERT INTO tables (
			venue_id, table_number, capacity, status,
			table_name, x, y, width, height, shape, rotation, color, floor,
			is_merged, merged_from, guests, assigned_waiter_id
		)
		VALUES (
			$1,
			(SELECT COALESCE(MAX(table_number), 0) + 1 FROM tables WHERE venue_id = $1 AND deleted_at IS NULL),
			$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
		RETURNING ` + tableColumns

	row := r.db.Pool.QueryRow(ctx, query,
		t.VenueID, t.Capacity, t.Status,
		t.Name, t.X, t.Y, t.Width, t.Height, t.Shape, t.Rotation, colorParam, t.Floor,
		t.IsMerged, mergedFromJSON, t.Guests, waiterParam,
	)

	created, err := r.scanTable(row)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}
	*t = created
	return nil
}

func (r *TableRepository) FindAll(ctx context.Context, venueID int) ([]table.Table, error) {
	query := `SELECT ` + tableColumns + `
		FROM tables
		WHERE venue_id = $1 AND deleted_at IS NULL
		ORDER BY table_number`

	rows, err := r.db.Pool.Query(ctx, query, venueID)
	if err != nil {
		return nil, fmt.Errorf("list tables: %w", err)
	}
	defer rows.Close()

	var tables []table.Table
	for rows.Next() {
		t, err := r.scanTable(rows)
		if err != nil {
			return nil, fmt.Errorf("scan table: %w", err)
		}
		tables = append(tables, t)
	}
	return tables, nil
}

func (r *TableRepository) FindByID(ctx context.Context, id int64, venueID int) (*table.Table, error) {
	query := `SELECT ` + tableColumns + `
		FROM tables WHERE id_table = $1 AND venue_id = $2 AND deleted_at IS NULL`

	row := r.db.Pool.QueryRow(ctx, query, id, venueID)
	t, err := r.scanTable(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("mesa no encontrada")
		}
		return nil, fmt.Errorf("get table: %w", err)
	}
	return &t, nil
}

// Update makes a partial update (status + arrival_time only, used by PATCH /estado).
func (r *TableRepository) Update(ctx context.Context, id int64, venueID int, updates *table.TableUpdates) error {
	query := `
		UPDATE tables
		SET
			table_number = COALESCE($1, table_number),
			capacity     = COALESCE($2, capacity),
			status       = COALESCE($3, status),
			arrival_time = COALESCE($4, arrival_time),
			updated_at   = now()
		WHERE id_table = $5 AND venue_id = $6 AND deleted_at IS NULL`

	var number, capacity pgtype.Int4
	var status pgtype.Text
	var arrivalTime pgtype.Timestamptz

	if updates.Number != nil {
		number = pgtype.Int4{Int32: int32(*updates.Number), Valid: true}
	}
	if updates.Capacity != nil {
		capacity = pgtype.Int4{Int32: int32(*updates.Capacity), Valid: true}
	}
	if updates.Status != nil {
		status = pgtype.Text{String: *updates.Status, Valid: true}
	}
	if updates.ArrivalTime != nil {
		arrivalTime = pgtype.Timestamptz{Time: *updates.ArrivalTime, Valid: true}
	}

	_, err := r.db.Pool.Exec(ctx, query, number, capacity, status, arrivalTime, id, venueID)
	if err != nil {
		return fmt.Errorf("update table: %w", err)
	}
	return nil
}

// FullUpdate replaces all mutable fields of a table (used by PUT /mesas/:id).
func (r *TableRepository) FullUpdate(ctx context.Context, id int64, venueID int, t table.Table) error {
	var mergedFromJSON []byte
	if len(t.MergedFrom) > 0 {
		mergedFromJSON, _ = json.Marshal(t.MergedFrom)
	}

	var colorParam interface{} = nil
	if t.Color != nil {
		colorParam = *t.Color
	}
	var waiterParam interface{} = nil
	if t.AssignedWaiterID != nil {
		waiterParam = *t.AssignedWaiterID
	}
	var arrivalParam interface{} = nil
	if t.ArrivalTime != nil {
		arrivalParam = *t.ArrivalTime
	}

	query := `
		UPDATE tables SET
			table_name        = $1,
			capacity          = $2,
			status            = $3,
			arrival_time      = $4,
			x                 = $5,
			y                 = $6,
			width             = $7,
			height            = $8,
			shape             = $9,
			rotation          = $10,
			color             = $11,
			floor             = $12,
			is_merged         = $13,
			merged_from       = $14,
			guests            = $15,
			assigned_waiter_id= $16,
			updated_at        = now()
		WHERE id_table = $17 AND venue_id = $18 AND deleted_at IS NULL`

	_, err := r.db.Pool.Exec(ctx, query,
		t.Name, t.Capacity, t.Status, arrivalParam,
		t.X, t.Y, t.Width, t.Height, t.Shape, t.Rotation, colorParam, t.Floor,
		t.IsMerged, mergedFromJSON, t.Guests, waiterParam,
		id, venueID,
	)
	if err != nil {
		return fmt.Errorf("full update table: %w", err)
	}
	return nil
}

func (r *TableRepository) Delete(ctx context.Context, id int64, venueID int) error {
	_, err := r.db.Pool.Exec(ctx,
		`DELETE FROM tables WHERE id_table = $1 AND venue_id = $2`,
		id, venueID,
	)
	return err
}
