package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/user"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	q  *sqlc.Queries
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{
		q:  sqlc.New(db.Pool),
		db: db,
	}
}

func toDomainUser(p sqlc.User) *user.User {
	u := &user.User{
		ID:           int(p.ID),
		VenueID:      int(p.VenueID),
		RoleID:       int(p.RoleID),
		Name:         p.Name,
		Email:        p.Email,
		PasswordHash: p.PasswordHash,
		IsActive:     p.IsActive.Bool,
		CreatedAt:    p.CreatedAt.Time,
		UpdatedAt:    p.UpdatedAt.Time,
	}
	if p.Phone.Valid {
		u.Phone = &p.Phone.String
	}
	if p.LastAccess.Valid {
		t := p.LastAccess.Time
		u.LastAccess = &t
	}
	return u
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	params := sqlc.CreateUserParams{
		VenueID:      int32(u.VenueID),
		RoleID:       int32(u.RoleID),
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}

	result, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainUser(result), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*user.User, error) {
	result, err := r.q.GetUserByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return toDomainUser(result), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	result, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return toDomainUser(result), nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) (*user.User, error) {
	params := sqlc.UpdateUserParams{
		ID:    int32(u.ID),
		Name:  u.Name,
		Email: u.Email,
		IsActive: pgtype.Bool{
			Bool:  u.IsActive,
			Valid: true,
		},
		RoleID: int32(u.RoleID),
	}

	result, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainUser(result), nil
}

func (r *UserRepository) ListByVenue(ctx context.Context, venueID int) ([]*user.User, error) {
	results, err := r.q.ListUsersByVenue(ctx, int32(venueID))
	if err != nil {
		return nil, err
	}

	users := make([]*user.User, len(results))
	for i, result := range results {
		users[i] = toDomainUser(result)
	}
	return users, nil
}

func (r *UserRepository) UpdateLastAccess(ctx context.Context, id int) error {
	return r.q.UpdateUserLastAccess(ctx, int32(id))
}
