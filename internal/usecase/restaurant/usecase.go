package restaurant

import (
	"context"

	domainRest "github.com/AXONcompany/POS/internal/domain/restaurant"
)

type RestaurantRepository interface {
	GetByID(ctx context.Context, id int) (*domainRest.Restaurant, error)
	Update(ctx context.Context, rest *domainRest.Restaurant) (*domainRest.Restaurant, error)
}

type Usecase struct {
	repo RestaurantRepository
}

func NewUsecase(repo RestaurantRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) GetMyRestaurant(ctx context.Context, restaurantID int) (*domainRest.Restaurant, error) {
	return uc.repo.GetByID(ctx, restaurantID)
}

func (uc *Usecase) UpdateRestaurantInfo(ctx context.Context, restaurantID int, name, address, phone string) (*domainRest.Restaurant, error) {
	rest, err := uc.repo.GetByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	rest.Name = name
	rest.Address = address
	rest.Phone = phone

	return uc.repo.Update(ctx, rest)
}
