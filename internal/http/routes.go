package http

import (
	"log"
	"net/http"

	"github.com/AXONcompany/POS/internal/http/ingredient"
	"github.com/AXONcompany/POS/internal/http/product"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine, ingredientHandler *ingredient.IngredientHandler, productHandler *product.Handler) {

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

	//ingredientes
	ingredients := r.Group("/ingredients")
	{
		ingredients.GET("", ingredientHandler.GetAll)
		ingredients.GET("/report", ingredientHandler.GetInventoryReport)
		ingredients.POST("", ingredientHandler.Create)
		ingredients.GET("/:id", ingredientHandler.GetByID)
		ingredients.PUT("/:id", ingredientHandler.Update)
		ingredients.DELETE("/:id", ingredientHandler.Delete)
	}

	// Categories
	categories := r.Group("/categories")
	{
		categories.POST("", productHandler.CreateCategory)
		categories.GET("", productHandler.GetAllCategories)
	}

	// Products
	products := r.Group("/products")
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("", productHandler.GetAllProducts)
		products.POST("/:id/ingredients", productHandler.AddIngredient)
		products.GET("/:id/ingredients", productHandler.GetIngredients)
	}

	log.Printf("Registered POST /ingredients, /categories, /products")

}
