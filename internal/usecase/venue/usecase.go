package venue

import (
	"context"
	"errors"
	"fmt"

	domainVenue "github.com/AXONcompany/POS/internal/domain/venue"
)

type Repository interface {
	Create(ctx context.Context, v *domainVenue.Venue) (*domainVenue.Venue, error)
	GetByID(ctx context.Context, id int) (*domainVenue.Venue, error)
	ListByOwner(ctx context.Context, ownerID int) ([]*domainVenue.Venue, error)
	Update(ctx context.Context, v *domainVenue.Venue) (*domainVenue.Venue, error)
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) CreateVenue(ctx context.Context, ownerID int, name, address, phone string) (*domainVenue.Venue, error) {
	if name == "" {
		return nil, errors.New("nombre de sede es requerido")
	}

	v := &domainVenue.Venue{
		OwnerID:  ownerID,
		Name:     name,
		Address:  address,
		Phone:    phone,
		IsActive: true,
	}

	return uc.repo.Create(ctx, v)
}

func (uc *Usecase) GetVenueByID(ctx context.Context, id int) (*domainVenue.Venue, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *Usecase) ListVenuesByOwner(ctx context.Context, ownerID int) ([]*domainVenue.Venue, error) {
	return uc.repo.ListByOwner(ctx, ownerID)
}

func (uc *Usecase) UpdateVenue(ctx context.Context, id int, name, address, phone string) (*domainVenue.Venue, error) {
	current, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("sede no encontrada: %w", err)
	}

	if name != "" {
		current.Name = name
	}
	if address != "" {
		current.Address = address
	}
	if phone != "" {
		current.Phone = phone
	}

	return uc.repo.Update(ctx, current)
}
