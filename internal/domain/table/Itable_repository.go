package table

import "context"

type TableRepository interface {
	FindAll(ctx context.Context) ([]Table, error)
	FindByID(ctx context.Context, id int64) (*Table, error)
	Create(ctx context.Context, table *Table) error
	Update(ctx context.Context, id int64, updates *TableUpdates) error
	Delete(ctx context.Context, id int64) error

	AssignWaitressToTable(ctx context.Context, tableID int64, waitressID int64) error
	RemoveWaitressFromTable(ctx context.Context, tableID int64, waitressID int64) error
	FindWaitressesByTableID(ctx context.Context, tableID int64) ([]TableWaitress, error)
}
