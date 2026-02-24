package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey       = "user_id"
	RoleIDKey       = "role_id"
	RestaurantIDKey = "restaurant_id"
	EmailKey        = "email"
)

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		if userID, ok := claims["sub"].(float64); ok {
			c.Set(UserIDKey, int(userID))
		}
		if roleID, ok := claims["role_id"].(float64); ok {
			c.Set(RoleIDKey, int(roleID))
		}
		if restID, ok := claims["restaurant_id"].(float64); ok {
			c.Set(RestaurantIDKey, int(restID))
		}
		if email, ok := claims["email"].(string); ok {
			c.Set(EmailKey, email)
		}

		c.Next()
	}
}

// RequireRole checks if the user has a specific role ID.
func RequireRole(roleID int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleID, exists := c.Get(RoleIDKey)
		if !exists || userRoleID.(int) != roleID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
			return
		}
		c.Next()
	}
}
