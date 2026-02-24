package http

import (
	"github.com/AXONcompany/POS/internal/config"
	"github.com/AXONcompany/POS/internal/http/ingredient"
	"github.com/AXONcompany/POS/internal/http/product"
	"github.com/AXONcompany/POS/internal/http/table"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, ingredientHandler *ingredient.IngredientHandler, productHandler *product.Handler, tableHandler *table.Handler) *gin.Engine {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	RegisterRouters(r, ingredientHandler, productHandler, tableHandler)
	return r

}
