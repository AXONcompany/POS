package http

import (
	"log"
	"net/http"

	"github.com/AXONcompany/POS/internal/http/ingredient"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine, ingredientHandler *ingredient.IngredientHandler) {

	log.Printf("RegisterRouters called, ingredientHandler is nil: %v", ingredientHandler == nil)

	//ver si está vivo
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	//r.GET("/ping", func(c *gin.Context) {
	//	c.String(http.StatusOK, "server say: pong")
	//})

	//ingredientes
	ingredients := r.Group("/ingredients")
	{
		ingredients.POST("", ingredientHandler.Create)
	}
	log.Printf("Registered POST /ingredients")

}
