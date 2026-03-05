package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/AXONcompany/POS/internal/domain/owner"
	"github.com/AXONcompany/POS/internal/domain/session"
	domainSession "github.com/AXONcompany/POS/internal/domain/session"
	domainUser "github.com/AXONcompany/POS/internal/domain/user"
	"github.com/AXONcompany/POS/internal/domain/venue"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domainUser.User, error)
	GetByID(ctx context.Context, id int) (*domainUser.User, error)
	Create(ctx context.Context, u *domainUser.User) (*domainUser.User, error)
	UpdateLastAccess(ctx context.Context, id int) error
}

type SessionRepository interface {
	Create(ctx context.Context, s *domainSession.Session) (*domainSession.Session, error)
	GetByToken(ctx context.Context, refreshToken string) (*domainSession.Session, error)
	Revoke(ctx context.Context, refreshToken string) error
}

type OwnerRepository interface {
	Create(ctx context.Context, o *owner.Owner) (*owner.Owner, error)
	GetByID(ctx context.Context, id int) (*owner.Owner, error)
	GetByEmail(ctx context.Context, email string) (*owner.Owner, error)
}

type VenueRepository interface {
	Create(ctx context.Context, v *venue.Venue) (*venue.Venue, error)
}

type Usecase struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
	ownerRepo   OwnerRepository
	venueRepo   VenueRepository
	jwtSecret   []byte
}

func NewUsecase(userRepo UserRepository, sessionRepo SessionRepository, jwtSecret string, ownerRepo OwnerRepository, venueRepo VenueRepository) *Usecase {
	return &Usecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		ownerRepo:   ownerRepo,
		venueRepo:   venueRepo,
		jwtSecret:   []byte(jwtSecret),
	}
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	User         *domainUser.User
}

func (uc *Usecase) Login(ctx context.Context, email, password, deviceInfo, ipAddress string) (*TokenResponse, error) {
	u, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !u.IsActive {
		return nil, errors.New("user is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate Access Token (15m)
	accessToken, err := uc.generateToken(u, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token (7d)
	refreshToken, err := uc.generateToken(u, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Save session
	s := &session.Session{
		UserID:       u.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		DeviceInfo:   deviceInfo,
		IPAddress:    ipAddress,
	}
	_, err = uc.sessionRepo.Create(ctx, s)
	if err != nil {
		return nil, err
	}

	// Update last access
	_ = uc.userRepo.UpdateLastAccess(ctx, u.ID)
	// Re-fetch to get updated last_access
	u, _ = uc.userRepo.GetByEmail(ctx, email)

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         u,
	}, nil
}

func (uc *Usecase) generateToken(u *domainUser.User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":      u.ID,
		"email":    u.Email,
		"role_id":  u.RoleID,
		"venue_id": u.VenueID,
		"exp":      time.Now().Add(duration).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.jwtSecret)
}

func (uc *Usecase) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	s, err := uc.sessionRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if s.IsRevoked || s.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("expired or revoked token")
	}

	// Revoke old session
	_ = uc.sessionRepo.Revoke(ctx, refreshToken)

	u, err := uc.userRepo.GetByID(ctx, s.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens directly without re-validating password
	accessToken, err := uc.generateToken(u, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	newRefresh, err := uc.generateToken(u, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	// Create new session
	newSession := &session.Session{
		UserID:       u.ID,
		RefreshToken: newRefresh,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		DeviceInfo:   s.DeviceInfo,
		IPAddress:    s.IPAddress,
	}
	if _, err = uc.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefresh,
		User:         u,
	}, nil
}

// RegisterUser crea un nuevo usuario con password hasheado.
func (uc *Usecase) RegisterUser(ctx context.Context, name, email, rawPassword string, roleID, venueID int, phone string) (*domainUser.User, error) {
	// Verificar que el email no existe
	existing, _ := uc.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u := &domainUser.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       roleID,
		VenueID:      venueID,
		IsActive:     true,
	}

	if phone != "" {
		u.Phone = &phone
	}

	return uc.userRepo.Create(ctx, u)
}

// RegisterOwnerWithVenue registra un propietario nuevo con su primera sede.
// Crea: owner -> venue -> user(role=PROPIETARIO) y retorna tokens para login inmediato.
func (uc *Usecase) RegisterOwnerWithVenue(ctx context.Context, ownerName, email, rawPassword, venueName, address, phone, deviceInfo, ipAddress string) (*TokenResponse, error) {
	// Verificar que el email no existe en users ni en owners
	existingUser, _ := uc.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}
	existingOwner, _ := uc.ownerRepo.GetByEmail(ctx, email)
	if existingOwner != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	// 1. Crear owner
	o := &owner.Owner{
		Name:         ownerName,
		Email:        email,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}
	createdOwner, err := uc.ownerRepo.Create(ctx, o)
	if err != nil {
		return nil, fmt.Errorf("creating owner: %w", err)
	}

	// 2. Crear venue
	v := &venue.Venue{
		OwnerID:  createdOwner.ID,
		Name:     venueName,
		Address:  address,
		Phone:    phone,
		IsActive: true,
	}
	createdVenue, err := uc.venueRepo.Create(ctx, v)
	if err != nil {
		return nil, fmt.Errorf("creating venue: %w", err)
	}

	// 3. Crear user con role PROPIETARIO vinculado a la venue
	u := &domainUser.User{
		Name:         ownerName,
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       1, // PROPIETARIO
		VenueID:      createdVenue.ID,
		IsActive:     true,
	}
	createdUser, err := uc.userRepo.Create(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	// 4. Generar tokens y sesion
	accessToken, err := uc.generateToken(createdUser, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := uc.generateToken(createdUser, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	s := &session.Session{
		UserID:       createdUser.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		DeviceInfo:   deviceInfo,
		IPAddress:    ipAddress,
	}
	_, err = uc.sessionRepo.Create(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         createdUser,
	}, nil
}

// RegisterWaiter crea un mesero con credenciales generadas automaticamente.
// Retorna el usuario creado y el password en texto plano (se muestra una sola vez).
func (uc *Usecase) RegisterWaiter(ctx context.Context, name, email string, venueID int) (*domainUser.User, string, error) {
	existing, _ := uc.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, "", errors.New("email already registered")
	}

	// Generar password aleatorio de 8 caracteres hex (16 chars)
	rawPassword, err := generateRandomPassword(8)
	if err != nil {
		return nil, "", fmt.Errorf("generating password: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("hashing password: %w", err)
	}

	u := &domainUser.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       3, // MESERO
		VenueID:      venueID,
		IsActive:     true,
	}

	created, err := uc.userRepo.Create(ctx, u)
	if err != nil {
		return nil, "", fmt.Errorf("creating waiter: %w", err)
	}

	return created, rawPassword, nil
}

// generateRandomPassword genera un password aleatorio de n bytes (retorna 2n chars hex).
func generateRandomPassword(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GetUserByID obtiene un usuario por su ID.
func (uc *Usecase) GetUserByID(ctx context.Context, id int) (*domainUser.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

// RevokeSession revoca un refresh token.
func (uc *Usecase) RevokeSession(ctx context.Context, refreshToken string) error {
	return uc.sessionRepo.Revoke(ctx, refreshToken)
}
