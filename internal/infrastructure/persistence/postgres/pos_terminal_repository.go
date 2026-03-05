package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/pos"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type POSTerminalRepository struct {
	q *sqlc.Queries
}

func NewPOSTerminalRepository(db *DB) *POSTerminalRepository {
	return &POSTerminalRepository{q: sqlc.New(db.Pool)}
}

func toDomainTerminal(t sqlc.PosTerminal) *pos.Terminal {
	return &pos.Terminal{
		ID:           int(t.ID),
		VenueID:      int(t.VenueID),
		TerminalName: t.TerminalName,
		IsActive:     t.IsActive.Bool,
		CreatedAt:    t.CreatedAt.Time,
		UpdatedAt:    t.UpdatedAt.Time,
	}
}

func (r *POSTerminalRepository) Create(ctx context.Context, t *pos.Terminal) (*pos.Terminal, error) {
	row, err := r.q.CreateTerminal(ctx, sqlc.CreateTerminalParams{
		VenueID:      int32(t.VenueID),
		TerminalName: t.TerminalName,
	})
	if err != nil {
		return nil, err
	}
	return toDomainTerminal(row), nil
}

func (r *POSTerminalRepository) GetByID(ctx context.Context, id int) (*pos.Terminal, error) {
	row, err := r.q.GetTerminalByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return toDomainTerminal(row), nil
}

func (r *POSTerminalRepository) ListByVenue(ctx context.Context, venueID int) ([]*pos.Terminal, error) {
	rows, err := r.q.ListTerminalsByVenue(ctx, int32(venueID))
	if err != nil {
		return nil, err
	}
	terminals := make([]*pos.Terminal, len(rows))
	for i, row := range rows {
		terminals[i] = toDomainTerminal(row)
	}
	return terminals, nil
}

func (r *POSTerminalRepository) Update(ctx context.Context, t *pos.Terminal) (*pos.Terminal, error) {
	row, err := r.q.UpdateTerminal(ctx, sqlc.UpdateTerminalParams{
		ID:           int32(t.ID),
		TerminalName: t.TerminalName,
		IsActive:     pgtype.Bool{Bool: t.IsActive, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return toDomainTerminal(row), nil
}
