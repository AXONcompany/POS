package table_test

import (
	"context"
	"errors"
	"testing"

	domainTable "github.com/AXONcompany/POS/internal/domain/table"
	uc "github.com/AXONcompany/POS/internal/usecase/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTableRepo struct {
	mock.Mock
}

func (m *MockTableRepo) Create(ctx context.Context, t *domainTable.Table) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockTableRepo) FindAll(ctx context.Context, venueID int) ([]domainTable.Table, error) {
	args := m.Called(ctx, venueID)
	return args.Get(0).([]domainTable.Table), args.Error(1)
}

func (m *MockTableRepo) FindByID(ctx context.Context, id int64, venueID int) (*domainTable.Table, error) {
	args := m.Called(ctx, id, venueID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainTable.Table), args.Error(1)
}

func (m *MockTableRepo) Update(ctx context.Context, id int64, venueID int, updates *domainTable.TableUpdates) error {
	args := m.Called(ctx, id, venueID, updates)
	return args.Error(0)
}

func (m *MockTableRepo) Delete(ctx context.Context, id int64, venueID int) error {
	args := m.Called(ctx, id, venueID)
	return args.Error(0)
}

// --- Create ---

func TestCreate_InvalidNumber(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)

	cases := []int{0, -1, -100}
	for _, n := range cases {
		err := service.Create(context.Background(), &domainTable.Table{Number: n, Capacity: 4})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "positivo")
	}
	repo.AssertNotCalled(t, "Create")
}

func TestCreate_InvalidCapacity(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)

	cases := []int{0, -1, -10}
	for _, c := range cases {
		err := service.Create(context.Background(), &domainTable.Table{Number: 1, Capacity: c})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "capacidad")
	}
	repo.AssertNotCalled(t, "Create")
}

func TestCreate_Success(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	table := &domainTable.Table{Number: 5, Capacity: 4, VenueID: 1}
	repo.On("Create", mock.Anything, table).Return(nil)

	err := service.Create(ctx, table)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCreate_RepoError(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	table := &domainTable.Table{Number: 1, Capacity: 2, VenueID: 1}
	repo.On("Create", mock.Anything, table).Return(errors.New("duplicate table number"))

	err := service.Create(ctx, table)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// --- Update ---

func TestUpdate_NegativeCapacity(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)

	neg := -1
	err := service.Update(context.Background(), 1, 1, &domainTable.TableUpdates{Capacity: &neg})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "negativa")
	repo.AssertNotCalled(t, "Update")
}

func TestUpdate_ZeroCapacity(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)

	zero := 0
	err := service.Update(context.Background(), 1, 1, &domainTable.TableUpdates{Capacity: &zero})
	assert.Error(t, err)
	repo.AssertNotCalled(t, "Update")
}

func TestUpdate_NilCapacity_Passthrough(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	// Capacity=nil significa "no cambiar", debe pasar al repo sin error
	updates := &domainTable.TableUpdates{Capacity: nil}
	repo.On("Update", mock.Anything, int64(1), 1, updates).Return(nil)

	err := service.Update(ctx, 1, 1, updates)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUpdate_ValidCapacity(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	cap := 8
	updates := &domainTable.TableUpdates{Capacity: &cap}
	repo.On("Update", mock.Anything, int64(1), 1, updates).Return(nil)

	err := service.Update(ctx, 1, 1, updates)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// --- FindAll ---

func TestFindAll_Success(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	tables := []domainTable.Table{{ID: 1, Number: 1, Capacity: 4}}
	repo.On("FindAll", mock.Anything, 1).Return(tables, nil)

	result, err := service.FindAll(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestFindAll_EmptyVenue(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	repo.On("FindAll", mock.Anything, 99).Return([]domainTable.Table{}, nil)

	result, err := service.FindAll(ctx, 99)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

// --- Delete ---

func TestDelete_Success(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	repo.On("Delete", mock.Anything, int64(1), 1).Return(nil)

	err := service.Delete(ctx, 1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDelete_NotFound(t *testing.T) {
	repo := new(MockTableRepo)
	service := uc.NewUsecase(repo)
	ctx := context.Background()

	repo.On("Delete", mock.Anything, int64(99), 1).Return(errors.New("not found"))

	err := service.Delete(ctx, 99, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}
