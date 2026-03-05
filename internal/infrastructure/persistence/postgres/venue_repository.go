package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/venue"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type VenueRepository struct {
	q *sqlc.Queries
}

func NewVenueRepository(db *DB) *VenueRepository {
	return &VenueRepository{q: sqlc.New(db.Pool)}
}

func toDomainVenue(v sqlc.Venue) *venue.Venue {
	return &venue.Venue{
		ID:        int(v.ID),
		OwnerID:   int(v.OwnerID),
		Name:      v.Name,
		Address:   v.Address.String,
		Phone:     v.Phone.String,
		IsActive:  v.IsActive.Bool,
		CreatedAt: v.CreatedAt.Time,
		UpdatedAt: v.UpdatedAt.Time,
	}
}

func (r *VenueRepository) Create(ctx context.Context, v *venue.Venue) (*venue.Venue, error) {
	row, err := r.q.CreateVenue(ctx, sqlc.CreateVenueParams{
		OwnerID: int32(v.OwnerID),
		Name:    v.Name,
		Address: pgtype.Text{String: v.Address, Valid: v.Address != ""},
		Phone:   pgtype.Text{String: v.Phone, Valid: v.Phone != ""},
	})
	if err != nil {
		return nil, err
	}
	return toDomainVenue(row), nil
}

func (r *VenueRepository) GetByID(ctx context.Context, id int) (*venue.Venue, error) {
	row, err := r.q.GetVenueByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return toDomainVenue(row), nil
}

func (r *VenueRepository) ListByOwner(ctx context.Context, ownerID int) ([]*venue.Venue, error) {
	rows, err := r.q.ListVenuesByOwner(ctx, int32(ownerID))
	if err != nil {
		return nil, err
	}
	venues := make([]*venue.Venue, len(rows))
	for i, row := range rows {
		venues[i] = toDomainVenue(row)
	}
	return venues, nil
}

func (r *VenueRepository) Update(ctx context.Context, v *venue.Venue) (*venue.Venue, error) {
	row, err := r.q.UpdateVenue(ctx, sqlc.UpdateVenueParams{
		ID:       int32(v.ID),
		Name:     v.Name,
		Address:  pgtype.Text{String: v.Address, Valid: v.Address != ""},
		Phone:    pgtype.Text{String: v.Phone, Valid: v.Phone != ""},
		IsActive: pgtype.Bool{Bool: v.IsActive, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return toDomainVenue(row), nil
}
