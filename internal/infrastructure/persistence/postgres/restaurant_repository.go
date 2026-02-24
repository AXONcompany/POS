package postgres

import (
	"context"

	"github.com/AXONcompany/POS/internal/domain/restaurant"
	"github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type RestaurantRepository struct {
	q  *sqlc.Queries
	db *DB
}

func NewRestaurantRepository(db *DB) *RestaurantRepository {
	return &RestaurantRepository{
		q:  sqlc.New(db.Pool),
		db: db,
	}
}

func toDomainRestaurant(p sqlc.Restaurant) *restaurant.Restaurant {
	return &restaurant.Restaurant{
		ID:        int(p.ID),
		Name:      p.Name,
		Address:   p.Address.String,
		Phone:     p.Phone.String,
		IsActive:  p.IsActive.Bool,
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}
}

func (r *RestaurantRepository) Create(ctx context.Context, rest *restaurant.Restaurant) (*restaurant.Restaurant, error) {
	params := sqlc.CreateRestaurantParams{
		Name: rest.Name,
		Address: pgtype.Text{
			String: rest.Address,
			Valid:  rest.Address != "",
		},
		Phone: pgtype.Text{
			String: rest.Phone,
			Valid:  rest.Phone != "",
		},
	}

	result, err := r.q.CreateRestaurant(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainRestaurant(result), nil
}

func (r *RestaurantRepository) GetByID(ctx context.Context, id int) (*restaurant.Restaurant, error) {
	result, err := r.q.GetRestaurantByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return toDomainRestaurant(result), nil
}

func (r *RestaurantRepository) Update(ctx context.Context, rest *restaurant.Restaurant) (*restaurant.Restaurant, error) {
	params := sqlc.UpdateRestaurantParams{
		ID:   int32(rest.ID),
		Name: rest.Name,
		Address: pgtype.Text{
			String: rest.Address,
			Valid:  rest.Address != "",
		},
		Phone: pgtype.Text{
			String: rest.Phone,
			Valid:  rest.Phone != "",
		},
		IsActive: pgtype.Bool{
			Bool:  rest.IsActive,
			Valid: true,
		},
	}

	result, err := r.q.UpdateRestaurant(ctx, params)
	if err != nil {
		return nil, err
	}

	return toDomainRestaurant(result), nil
}
