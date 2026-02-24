package rest

import (
	"log"
	"net/http"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/ingredient"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/product"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	orderrest "github.com/AXONcompany/POS/internal/infrastructure/rest/order"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine, ingredientHandler *ingredient.IngredientHandler, productHandler *product.Handler, orderHandler *orderrest.Handler, jwtSecret []byte) {

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

	// Menu
	menu := r.Group("/menu")
	{
		menu.GET("", productHandler.GetMenu)
		menu.POST("", productHandler.CreateMenuItem)
		menu.PATCH("/:id", productHandler.UpdateMenuItem)
	}

	// Roles mappings (adjust map to your DB setup)
	const RolePropietario = 1
	const RoleCajero = 2
	const RoleMesero = 3

	// Orders (Protected)
	orders := r.Group("/orders")
	orders.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Both Mesero & Cajero can create orders
		orders.POST("", orderHandler.CreateOrder)

		// Only Cajero can checkout
		checkout := orders.Group("/:id/checkout")
		checkout.Use(middleware.RequireRole(RoleCajero))
		checkout.POST("", orderHandler.CheckoutOrder)
	}

	log.Printf("Registered POST /ingredients, /categories, /products, /menu, /orders")

}
