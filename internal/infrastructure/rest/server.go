package rest

import (
	"github.com/AXONcompany/POS/internal/config"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/auth"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/ingredient"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/order"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/owner"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/payment"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/pos"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/product"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/report"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/table"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/user"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/venue"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, ingredientHandler *ingredient.IngredientHandler, productHandler *product.Handler, authHandler *auth.Handler, orderHandler *order.Handler, tableHandler *table.Handler, userHandler *user.Handler, paymentHandler *payment.Handler, reportHandler *report.Handler, ownerHandler *owner.Handler, venueHandler *venue.Handler, posHandler *pos.Handler) *gin.Engine {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS Setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // TODO: Para prod real, especificar el dominio
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	jwtSecret := []byte(cfg.JWTSecret)

	RegisterRouters(r, ingredientHandler, productHandler, authHandler, orderHandler, tableHandler, userHandler, paymentHandler, reportHandler, ownerHandler, venueHandler, posHandler, jwtSecret)

	return r

}
