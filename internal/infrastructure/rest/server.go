package rest

import (
	"github.com/AXONcompany/POS/internal/config"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/auth"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/ingredient"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/order"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/product"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/table"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, ingredientHandler *ingredient.IngredientHandler, productHandler *product.Handler, authHandler *auth.Handler, orderHandler *order.Handler, tableHandler *table.Handler) *gin.Engine {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	jwtSecret := []byte(cfg.JWTSecret) // assuming config has JWTSecret

	RegisterRouters(r, ingredientHandler, productHandler, authHandler, orderHandler, tableHandler, jwtSecret)

	// Auth Routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
	}

	return r

}
