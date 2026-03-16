package table

import (
	"context"
	"errors"
	"time"

	"github.com/AXONcompany/POS/internal/domain/table"
)

type TableRepository interface {
	Create(ctx context.Context, t *table.Table) error
	FindAll(ctx context.Context, venueID int) ([]table.Table, error)
	FindByID(ctx context.Context, id int64, venueID int) (*table.Table, error)
	Update(ctx context.Context, id int64, venueID int, updates *table.TableUpdates) error
	Delete(ctx context.Context, id int64, venueID int) error
}

type TableAssignmentRepository interface {
	Assign(ctx context.Context, tableID int64, userID, venueID int) (*table.Assignment, error)
	Unassign(ctx context.Context, tableID int64, venueID int) error
	GetActive(ctx context.Context, tableID int64, venueID int) (*table.Assignment, error)
	ListByTable(ctx context.Context, tableID int64, venueID int) ([]table.AssignmentDetail, error)
}

type Usecase struct {
	repo           TableRepository
	assignRepo     TableAssignmentRepository
	contextTimeout time.Duration
}

func NewUsecase(repo TableRepository, assignRepo TableAssignmentRepository) *Usecase {
	return &Usecase{
		repo:           repo,
		assignRepo:     assignRepo,
		contextTimeout: time.Duration(2) * time.Second,
	}
}

func (s *Usecase) Create(c context.Context, t *table.Table) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	if t.Number <= 0 {
		return errors.New("el numero de mesa debe ser positivo")
	}
	if t.Capacity <= 0 {
		return errors.New("la capacidad de la mesa debe ser mayor a 0")
	}

	return s.repo.Create(ctx, t)
}

func (s *Usecase) FindAll(c context.Context, venueID int) ([]table.Table, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.FindAll(ctx, venueID)
}

func (s *Usecase) FindByID(c context.Context, id int64, venueID int) (*table.Table, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.FindByID(ctx, id, venueID)
}

func (s *Usecase) Update(c context.Context, id int64, venueID int, updates *table.TableUpdates) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	if updates.Capacity != nil && *updates.Capacity <= 0 {
		return errors.New("la capacidad no puede ser negativa")
	}

	return s.repo.Update(ctx, id, venueID, updates)
}

func (s *Usecase) Delete(c context.Context, id int64, venueID int) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.repo.Delete(ctx, id, venueID)
}

// AssignWaiter asigna un mesero a una mesa. Si ya hay uno asignado, lo desasigna primero.
func (s *Usecase) AssignWaiter(c context.Context, tableID int64, userID, venueID int) (*table.Assignment, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.assignRepo.Assign(ctx, tableID, userID, venueID)
}

// UnassignWaiter desasigna el mesero activo de una mesa.
func (s *Usecase) UnassignWaiter(c context.Context, tableID int64, venueID int) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.assignRepo.Unassign(ctx, tableID, venueID)
}

// GetAssignments devuelve el historial de asignaciones de una mesa.
func (s *Usecase) GetAssignments(c context.Context, tableID int64, venueID int) ([]table.AssignmentDetail, error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()
	return s.assignRepo.ListByTable(ctx, tableID, venueID)
}
