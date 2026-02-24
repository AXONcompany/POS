package table

import (
	"context"
	"errors"
	"time"

	"github.com/AXONcompany/POS/internal/domain/table"
)

type TableRepository interface {
	Create(ctx context.Context, t *table.Table) error
	FindAll(ctx context.Context) ([]table.Table, error)
	FindByID(ctx context.Context, id int64) (*table.Table, error)
	Update(ctx context.Context, id int64, updates *table.TableUpdates) error
	Delete(ctx context.Context, id int64) error
	AssignWaitressToTable(ctx context.Context, tableID int64, waitressID int64) error
	RemoveWaitressFromTable(ctx context.Context, tableID int64, waitressID int64) error
	FindWaitressesByTableID(ctx context.Context, tableID int64) ([]table.TableWaitress, error)
}

type Usecase struct {
	repo           TableRepository
	contextTimeout time.Duration
}

func NewUsecase(repo TableRepository) *Usecase {
	return &Usecase{
		repo:           repo,
		contextTimeout: time.Duration(2) * time.Second,
	}
}

func (s *Usecase) Create(c context.Context, t *table.Table) error {
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

func (s *Usecase) FindAll(c context.Context) ([]table.Table, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.FindAll(ctx)
}

func (s *Usecase) FindByID(c context.Context, id int64) (*table.Table, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.FindByID(ctx, id)
}

func (s *Usecase) Update(c context.Context, id int64, updates *table.TableUpdates) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	if updates.Capacity != nil && *updates.Capacity <= 0 {
		return errors.New("la capacidad no puede ser negativa")
	}

	return s.repo.Update(ctx, id, updates)
}

func (s *Usecase) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.Delete(ctx, id)
}

func (s *Usecase) AssignWaitress(c context.Context, tableID int64, waitressID int64) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	return s.repo.AssignWaitressToTable(ctx, tableID, waitressID)
}

func (s *Usecase) FindWaitressesByTable(c context.Context, tableID int64) ([]table.TableWaitress, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.FindWaitressesByTableID(ctx, tableID)
}
