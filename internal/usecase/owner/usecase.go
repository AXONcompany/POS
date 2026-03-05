package owner

import (
	"context"
	"errors"
	"fmt"

	domainOwner "github.com/AXONcompany/POS/internal/domain/owner"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Create(ctx context.Context, o *domainOwner.Owner) (*domainOwner.Owner, error)
	GetByID(ctx context.Context, id int) (*domainOwner.Owner, error)
	GetByEmail(ctx context.Context, email string) (*domainOwner.Owner, error)
	Update(ctx context.Context, o *domainOwner.Owner) (*domainOwner.Owner, error)
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) CreateOwner(ctx context.Context, name, email, rawPassword string) (*domainOwner.Owner, error) {
	if name == "" || email == "" || rawPassword == "" {
		return nil, errors.New("nombre, email y password son requeridos")
	}

	existing, _ := uc.repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email ya registrado")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	o := &domainOwner.Owner{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		IsActive:     true,
	}

	return uc.repo.Create(ctx, o)
}

func (uc *Usecase) GetOwnerByID(ctx context.Context, id int) (*domainOwner.Owner, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *Usecase) UpdateOwner(ctx context.Context, id int, name, email string) (*domainOwner.Owner, error) {
	current, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("propietario no encontrado: %w", err)
	}

	if name != "" {
		current.Name = name
	}
	if email != "" {
		current.Email = email
	}

	return uc.repo.Update(ctx, current)
}
