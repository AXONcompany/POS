package http

import (
	"log"
	"net/http"

	"github.com/AXONcompany/POS/internal/http/ingredient"
	tableHttp "github.com/AXONcompany/POS/internal/http/table"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine, ingredientHandler *ingredient.IngredientHandler, tableHandler *tableHttp.Handler) {

	log.Printf("RegisterRouters called, ingredientHandler is nil: %v", ingredientHandler == nil)

	//ver si está vivo
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "server say: pong")
	})

	tables := r.Group("/tables")
	{
		tables.POST("", tableHandler.Create)
		tables.GET("", tableHandler.GetAll)
		tables.GET("/:id", tableHandler.GetByID)
		tables.PATCH("/:id", tableHandler.Update)
		tables.DELETE("/:id", tableHandler.Delete)

		// Ruta especial
		tables.POST("/:id/assign", tableHandler.AssignWaitress)
	}
	log.Printf("Registered /tables routes")

	//ingredientes
	ingredients := r.Group("/ingredients")
	{
		ingredients.GET("", ingredientHandler.GetAll)
		ingredients.POST("", ingredientHandler.Create)
		ingredients.GET("/:id", ingredientHandler.GetByID)
		ingredients.PUT("/:id", ingredientHandler.Update)
		ingredients.DELETE("/:id", ingredientHandler.Delete)
	}
	log.Printf("Registered POST /ingredients")

}
