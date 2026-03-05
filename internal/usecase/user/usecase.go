package user

import (
	"context"
	"fmt"

	domainUser "github.com/AXONcompany/POS/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, u *domainUser.User) (*domainUser.User, error)
	GetByID(ctx context.Context, id int) (*domainUser.User, error)
	GetByEmail(ctx context.Context, email string) (*domainUser.User, error)
	Update(ctx context.Context, u *domainUser.User) (*domainUser.User, error)
	ListByVenue(ctx context.Context, venueID int) ([]*domainUser.User, error)
	UpdateLastAccess(ctx context.Context, id int) error
}

type Usecase struct {
	repo UserRepository
}

func NewUsecase(repo UserRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) CreateUser(ctx context.Context, u *domainUser.User, rawPassword string) (*domainUser.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u.PasswordHash = string(hashedPassword)
	return uc.repo.Create(ctx, u)
}

func (uc *Usecase) GetAllUsers(ctx context.Context, venueID int) ([]*domainUser.User, error) {
	return uc.repo.ListByVenue(ctx, venueID)
}

func (uc *Usecase) GetUserByID(ctx context.Context, id int) (*domainUser.User, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *Usecase) UpdateUser(ctx context.Context, u *domainUser.User) (*domainUser.User, error) {
	return uc.repo.Update(ctx, u)
}

func (uc *Usecase) DeleteUser(ctx context.Context, id int) error {
	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	u.IsActive = false
	_, err = uc.repo.Update(ctx, u)
	return err
}
