package pos

import (
	"context"
	"errors"
	"fmt"

	domainPos "github.com/AXONcompany/POS/internal/domain/pos"
)

type Repository interface {
	Create(ctx context.Context, t *domainPos.Terminal) (*domainPos.Terminal, error)
	GetByID(ctx context.Context, id int) (*domainPos.Terminal, error)
	ListByVenue(ctx context.Context, venueID int) ([]*domainPos.Terminal, error)
	Update(ctx context.Context, t *domainPos.Terminal) (*domainPos.Terminal, error)
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) CreateTerminal(ctx context.Context, venueID int, terminalName string) (*domainPos.Terminal, error) {
	if terminalName == "" {
		return nil, errors.New("nombre del terminal es requerido")
	}

	t := &domainPos.Terminal{
		VenueID:      venueID,
		TerminalName: terminalName,
		IsActive:     true,
	}

	return uc.repo.Create(ctx, t)
}

func (uc *Usecase) GetTerminalByID(ctx context.Context, id int) (*domainPos.Terminal, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *Usecase) ListTerminalsByVenue(ctx context.Context, venueID int) ([]*domainPos.Terminal, error) {
	return uc.repo.ListByVenue(ctx, venueID)
}

func (uc *Usecase) UpdateTerminal(ctx context.Context, id int, terminalName string, isActive *bool) (*domainPos.Terminal, error) {
	current, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("terminal no encontrado: %w", err)
	}

	if terminalName != "" {
		current.TerminalName = terminalName
	}
	if isActive != nil {
		current.IsActive = *isActive
	}

	return uc.repo.Update(ctx, current)
}
