package table

import (
	"context"
	"errors"
	"time"

	"github.com/AXONcompany/POS/internal/domain/table"
)

type Service struct {
	repo           table.TableRepository
	contextTimeout time.Duration
}

func NewService(repo table.TableRepository) *Service {
	return &Service{
		repo:           repo,
		contextTimeout: time.Duration(2) * time.Second,
	}
}

func (s *Service) Create(c context.Context, t *table.Table) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	if t.Number <= 0 {
		return errors.New("el número de mesa debe ser positivo")
	}
	if t.Capacity <= 0 {
		return errors.New("la capacidad de la mesa debe ser mayor a 0")
	}

	return s.repo.Create(ctx, t)
}

func (s *Service) FindAll(c context.Context) ([]table.Table, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.FindAll(ctx)
}

func (s *Service) FindByID(c context.Context, id int64) (*table.Table, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Update(c context.Context, id int64, updates *table.TableUpdates) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	if updates.Capacity != nil && *updates.Capacity <= 0 {
		return errors.New("la capacidad no puede ser negativa")
	}

	return s.repo.Update(ctx, id, updates)
}

func (s *Service) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.Delete(ctx, id)
}

func (s *Service) AssignWaitress(c context.Context, tableID int64, waitressID int64) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	return s.repo.AssignWaitressToTable(ctx, tableID, waitressID)
}

func (s *Service) FindWaitressesByTable(c context.Context, tableID int64) ([]table.TableWaitress, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.FindWaitressesByTableID(ctx, tableID)
}
