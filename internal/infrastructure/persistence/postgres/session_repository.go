package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/session"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type SessionRepository struct {
	q  *sqlc.Queries
	db *DB
}

func NewSessionRepository(db *DB) *SessionRepository {
	return &SessionRepository{
		q:  sqlc.New(db.Pool),
		db: db,
	}
}

func toDomainSession(p sqlc.Session) *session.Session {
	return &session.Session{
		ID:           int(p.ID),
		UserID:       int(p.UserID),
		RefreshToken: p.RefreshToken,
		ExpiresAt:    p.ExpiresAt.Time,
		DeviceInfo:   p.DeviceInfo.String,
		IPAddress:    p.IpAddress.String,
		IsRevoked:    p.IsRevoked.Bool,
		CreatedAt:    p.CreatedAt.Time,
	}
}

func (r *SessionRepository) Create(ctx context.Context, s *session.Session) (*session.Session, error) {
	params := sqlc.CreateSessionParams{
		UserID:       int32(s.UserID),
		RefreshToken: s.RefreshToken,
		ExpiresAt:    pgtype.Timestamptz{Time: s.ExpiresAt, Valid: !s.ExpiresAt.IsZero()},
		DeviceInfo: pgtype.Text{
			String: s.DeviceInfo,
			Valid:  s.DeviceInfo != "",
		},
		IpAddress: pgtype.Text{
			String: s.IPAddress,
			Valid:  s.IPAddress != "",
		},
	}

	result, err := r.q.CreateSession(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainSession(result), nil
}

func (r *SessionRepository) GetByToken(ctx context.Context, token string) (*session.Session, error) {
	result, err := r.q.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return toDomainSession(result), nil
}

func (r *SessionRepository) Revoke(ctx context.Context, token string) error {
	return r.q.RevokeSession(ctx, token)
}
