package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupRouter(jwtSecret []byte) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(middleware.AuthMiddleware(jwtSecret))

	// Dummy endpoint to verify context population
	r.GET("/protected", func(c *gin.Context) {
		roleID, _ := c.Get(middleware.RoleIDKey)
		c.JSON(http.StatusOK, gin.H{"role": roleID})
	})

	r.GET("/admin-only", middleware.RequireRoles(1), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	return r
}

func generateValidToken(secret []byte, roleID int, expired bool) string {
	exp := time.Now().Add(1 * time.Hour).Unix()
	if expired {
		exp = time.Now().Add(-1 * time.Hour).Unix()
	}

	claims := jwt.MapClaims{
		"sub":           float64(123),
		"role_id":       float64(roleID),
		"restaurant_id": float64(1),
		"exp":           exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString(secret)
	return s
}

func TestAuthMiddleware(t *testing.T) {
	secret := []byte("test-secret-key")
	router := setupRouter(secret)

	tests := []struct {
		name         string
		setupHeader  func(req *http.Request)
		expectedCode int
	}{
		{
			name:         "Missing Auth Header",
			setupHeader:  func(req *http.Request) {},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid Format (No Bearer)",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Basic token123")
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Malformed Token (Malicious User)",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid.token.str")
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Expired Token",
			setupHeader: func(req *http.Request) {
				token := generateValidToken(secret, 2, true)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid Signature (Tampered Token)",
			setupHeader: func(req *http.Request) {
				token := generateValidToken([]byte("wrong-secret"), 2, false)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Valid Token Correct Access",
			setupHeader: func(req *http.Request) {
				token := generateValidToken(secret, 2, false)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/protected", nil)
			tc.setupHeader(req)

			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestRequireRoleMiddleware(t *testing.T) {
	secret := []byte("test-secret-key")
	router := setupRouter(secret)

	t.Run("Valid Admin Role", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin-only", nil)

		token := generateValidToken(secret, 1, false) // 1 = Admin Propietario
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Forbidden Non-Admin Role (Cajero Trying to access Admin)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin-only", nil)

		token := generateValidToken(secret, 2, false) // 2 = Cajero
		req.Header.Set("Authorization", "Bearer "+token)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
