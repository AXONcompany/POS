package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/role"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
)

type RoleRepository struct {
	q  *sqlc.Queries
	db *DB
}

func NewRoleRepository(db *DB) *RoleRepository {
	return &RoleRepository{
		q:  sqlc.New(db.Pool),
		db: db,
	}
}

func toDomainRole(p sqlc.Role) *role.Role {
	return &role.Role{
		ID:          int(p.ID),
		Name:        p.Name,
		Description: p.Description.String,
	}
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*role.Role, error) {
	result, err := r.q.GetRoleByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return toDomainRole(result), nil
}

func (r *RoleRepository) GetByID(ctx context.Context, id int) (*role.Role, error) {
	result, err := r.q.GetRoleByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return toDomainRole(result), nil
}
