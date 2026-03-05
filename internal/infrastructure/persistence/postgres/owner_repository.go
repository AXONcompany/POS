package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/owner"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type OwnerRepository struct {
	q *sqlc.Queries
}

func NewOwnerRepository(db *DB) *OwnerRepository {
	return &OwnerRepository{q: sqlc.New(db.Pool)}
}

func toDomainOwner(o sqlc.Owner) *owner.Owner {
	return &owner.Owner{
		ID:           int(o.ID),
		Name:         o.Name,
		Email:        o.Email,
		PasswordHash: o.PasswordHash,
		IsActive:     o.IsActive.Bool,
		CreatedAt:    o.CreatedAt.Time,
		UpdatedAt:    o.UpdatedAt.Time,
	}
}

func (r *OwnerRepository) Create(ctx context.Context, o *owner.Owner) (*owner.Owner, error) {
	row, err := r.q.CreateOwner(ctx, sqlc.CreateOwnerParams{
		Name:         o.Name,
		Email:        o.Email,
		PasswordHash: o.PasswordHash,
	})
	if err != nil {
		return nil, err
	}
	return toDomainOwner(row), nil
}

func (r *OwnerRepository) GetByID(ctx context.Context, id int) (*owner.Owner, error) {
	row, err := r.q.GetOwnerByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return toDomainOwner(row), nil
}

func (r *OwnerRepository) GetByEmail(ctx context.Context, email string) (*owner.Owner, error) {
	row, err := r.q.GetOwnerByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return toDomainOwner(row), nil
}

func (r *OwnerRepository) Update(ctx context.Context, o *owner.Owner) (*owner.Owner, error) {
	row, err := r.q.UpdateOwner(ctx, sqlc.UpdateOwnerParams{
		ID:    int32(o.ID),
		Name:  o.Name,
		Email: o.Email,
		IsActive: pgtype.Bool{
			Bool:  o.IsActive,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}
	return toDomainOwner(row), nil
}
