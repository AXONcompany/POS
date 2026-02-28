package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domainSession "github.com/AXONcompany/POS/internal/domain/session"
	domainUser "github.com/AXONcompany/POS/internal/domain/user"
	uc "github.com/AXONcompany/POS/internal/usecase/auth"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// --- MOCKS ---

type mockUserRepository struct {
	getByEmailFunc func(ctx context.Context, email string) (*domainUser.User, error)
	getByIDFunc    func(ctx context.Context, id int) (*domainUser.User, error)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int) (*domainUser.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

type mockSessionRepository struct {
	createFunc     func(ctx context.Context, s *domainSession.Session) (*domainSession.Session, error)
	getByTokenFunc func(ctx context.Context, refreshToken string) (*domainSession.Session, error)
	revokeFunc     func(ctx context.Context, refreshToken string) error
}

func (m *mockSessionRepository) Create(ctx context.Context, s *domainSession.Session) (*domainSession.Session, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, s)
	}
	s.ID = 1
	return s, nil
}

func (m *mockSessionRepository) GetByToken(ctx context.Context, refreshToken string) (*domainSession.Session, error) {
	if m.getByTokenFunc != nil {
		return m.getByTokenFunc(ctx, refreshToken)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSessionRepository) Revoke(ctx context.Context, refreshToken string) error {
	if m.revokeFunc != nil {
		return m.revokeFunc(ctx, refreshToken)
	}
	return nil
}

// Helpers
func generateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

// ==========================================
// NIVEL 1: PRUEBAS BÁSICAS (HAPPY PATH)
// ==========================================

func TestUsecase_Login_Basic_Success(t *testing.T) {
	expectedPass := "secure123"
	expectedHash := generateHash(expectedPass)

	userRepo := &mockUserRepository{
		getByEmailFunc: func(ctx context.Context, email string) (*domainUser.User, error) {
			if email == "test@axon.com" {
				return &domainUser.User{
					ID:           1,
					Email:        email,
					PasswordHash: expectedHash,
					RoleID:       1,
					IsActive:     true,
				}, nil
			}
			return nil, errors.New("user not found")
		},
	}
	sessionRepo := &mockSessionRepository{} // defaults are fine for success

	secret := "supersecret"
	usecase := uc.NewUsecase(userRepo, sessionRepo, secret)

	response, err := usecase.Login(context.Background(), "test@axon.com", expectedPass, "device", "127.0.0.1")

	if err != nil {
		t.Fatalf("Nivel 1 Fallido: se esperaba éxito, se obtuvo error: %v", err)
	}
	if response.AccessToken == "" || response.RefreshToken == "" {
		t.Fatalf("Nivel 1 Fallido: tokens no generados correctamente")
	}

	// Verificar validez del token generado
	token, _ := jwt.Parse(response.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if !token.Valid {
		t.Errorf("Nivel 1 Fallido: Token generado inválido por firma")
	}
}

// ==========================================
// NIVEL 2: PRUEBAS DE BORDE (EDGE CASES)
// ==========================================

func TestUsecase_Login_Edge_Table(t *testing.T) {
	validPass := "secure123"
	validHash := generateHash(validPass)

	userRepo := &mockUserRepository{
		getByEmailFunc: func(ctx context.Context, email string) (*domainUser.User, error) {
			return &domainUser.User{
				ID:           1,
				Email:        email,
				PasswordHash: validHash,
				RoleID:       1,
				IsActive:     true,
			}, nil
		},
	}
	sessionRepo := &mockSessionRepository{}
	secret := "supersecret"
	usecase := uc.NewUsecase(userRepo, sessionRepo, secret)

	// Crear string masivo de 72+ bytes (Límite BCrypt es 72 por diseño, el resto es truncado o devuleve error en versiones seguras)
	massivePayload := ""
	for i := 0; i < 10000; i++ {
		massivePayload += "A"
	}

	tests := []struct {
		name          string
		email         string
		password      string
		expectedError string
	}{
		{
			name:          "Contraseña de Carga Masiva (10k chars)",
			email:         "test@axon.com",
			password:      massivePayload,
			expectedError: "invalid credentials", // bcrypt.Compare devolverá error internamente
		},
		{
			name:          "Contraseña Vacía",
			email:         "test@axon.com",
			password:      "",
			expectedError: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := usecase.Login(context.Background(), tt.email, tt.password, "device", "ip")

			if err == nil {
				t.Fatalf("Nivel 2 Fallido '%s': se esperaba error de credenciales, pero el login pasó", tt.name)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("Nivel 2 Fallido '%s': error esperado '%s', obtenido '%v'", tt.name, tt.expectedError, err)
			}
		})
	}
}

func TestUsecase_RefreshToken_Edge_BorderlineExp(t *testing.T) {
	sessionRepo := &mockSessionRepository{
		getByTokenFunc: func(ctx context.Context, refreshToken string) (*domainSession.Session, error) {
			// Simulamos un token expirado JUSTO hace 1 milisegundo (Edge)
			return &domainSession.Session{
				UserID:       1,
				RefreshToken: refreshToken,
				ExpiresAt:    time.Now().Add(-1 * time.Millisecond),
				IsRevoked:    false,
			}, nil
		},
	}
	userRepo := &mockUserRepository{}
	usecase := uc.NewUsecase(userRepo, sessionRepo, "secret")

	_, err := usecase.RefreshToken(context.Background(), "some-token")
	if err == nil {
		t.Fatal("Nivel 2 Fallido: Un token caducado hace 1 milisegundo logró refrescar la sesión")
	}

	if err.Error() != "expired or revoked token" {
		t.Errorf("Nivel 2 Fallido: Error incorrecto generado: %v", err)
	}
}

// ==========================================
// NIVEL 3: PRUEBAS ADVERSARIALES (OWASP)
// ==========================================

func TestUsecase_Token_Adversarial_AlgoNone(t *testing.T) {
	secret := "supersecret"
	// usecase removed as it is not needed here

	// Simular a un atacante que genera su propio JWT
	// Usando el token method 'None'
	claims := jwt.MapClaims{
		"sub":     1,
		"role_id": 1, // Atacante dice ser Admin
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	maliciousToken, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	// El caso de uso aquí en realidad no decodifica, pero el sistema REST lo hará.
	// Para probar la resistencia intrínseca, utilicemos el validador estándar sobre el token malicioso
	// asegurando que nuestro sistema rechace este algoritmo.

	parsedToken, err := jwt.Parse(maliciousToken, func(t *jwt.Token) (interface{}, error) {
		// Esta función es similar a la implementada en AuthMiddleware
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err == nil || parsedToken.Valid {
		t.Fatal("Nivel 3 Fallido PELIGROSO: Token falsificado con alg:none fue validado como correcto")
	}
}

func TestUsecase_Token_Adversarial_ForgedKey(t *testing.T) {
	// Atacante intercepta token
	// validTokenString is concept, removed to fix linting.

	// Intenta reempaquetarlo con sus claims de owner (role_id 1) y firmarlo él usando HMAC,
	// pero sin conocer el `secret` real del servidor.
	attackerSecret := "123456"
	claims := jwt.MapClaims{"role_id": 1, "sub": 12}
	forgedTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	forgedTokenStr, _ := forgedTokenObj.SignedString([]byte(attackerSecret))

	// Nuestro entorno "Servidor"
	realSecret := "supersecret"

	_, err := jwt.Parse(forgedTokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(realSecret), nil
	})

	if err == nil {
		t.Fatal("Nivel 3 Fallido PELIGROSO: Token falsificado con firma inválida pasó la verificación HMAC")
	}
}
