package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domainUser "github.com/AXONcompany/POS/internal/domain/user"
	uc "github.com/AXONcompany/POS/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *domainUser.User) (*domainUser.User, error) {
	args := m.Called(ctx, u)
	if args.Get(0) != nil {
		return args.Get(0).(*domainUser.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*domainUser.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, u *domainUser.User) (*domainUser.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserRepository) ListByRestaurant(ctx context.Context, restaurantID int) ([]*domainUser.User, error) {
	args := m.Called(ctx, restaurantID)
	return args.Get(0).([]*domainUser.User), args.Error(1)
}

func TestCreateUser_PasswordHashing(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.NewUsecase(mockRepo)

	ctx := context.Background()
	rawPassword := "SecurePass123"
	u := &domainUser.User{
		Name:         "John Doe",
		Email:        "john@test.com",
		RoleID:       2,
		RestaurantID: 1,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Run(func(args mock.Arguments) {
		userArg := args.Get(1).(*domainUser.User)
		userArg.ID = 100
		userArg.CreatedAt = time.Now()
	}).Return(u, nil)

	createdUser, err := usecase.CreateUser(ctx, u, rawPassword)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, 100, createdUser.ID)
	assert.NotEqual(t, rawPassword, createdUser.PasswordHash)

	err = bcrypt.CompareHashAndPassword([]byte(createdUser.PasswordHash), []byte(rawPassword))
	assert.NoError(t, err, "The hashed password should match the raw password when compared")

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_DBError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.NewUsecase(mockRepo)

	ctx := context.Background()
	u := &domainUser.User{Name: "Fail User"}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil, errors.New("db connection lost"))

	_, err := usecase.CreateUser(ctx, u, "pass123")

	assert.Error(t, err)
	assert.Equal(t, "db connection lost", err.Error())
	mockRepo.AssertExpectations(t)
}
