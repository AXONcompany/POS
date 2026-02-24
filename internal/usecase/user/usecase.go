package user

import (
	"context"
	"fmt"

	domainUser "github.com/AXONcompany/POS/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, u *domainUser.User) (*domainUser.User, error)
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
