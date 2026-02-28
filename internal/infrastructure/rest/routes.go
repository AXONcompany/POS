package rest

import (
	"log"
	"net/http"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/auth"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/ingredient"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	orderrest "github.com/AXONcompany/POS/internal/infrastructure/rest/order"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/product"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/table"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine, ingredientHandler *ingredient.IngredientHandler, productHandler *product.Handler, authHandler *auth.Handler, orderHandler *orderrest.Handler, tableHandler *table.Handler, jwtSecret []byte) {

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

	// Roles mappings (adjust map to your DB setup)
	const RolePropietario = 1
	const RoleCajero = 2
	const RoleMesero = 3

	// Protected API Group
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// --- TABLES ---
		tables := api.Group("/tables")
		{
			tables.GET("", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), tableHandler.GetAll)
			tables.GET("/:id", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), tableHandler.GetByID)

			tables.POST("", middleware.RequireRoles(RoleCajero, RolePropietario), tableHandler.Create)
			tables.PATCH("/:id", middleware.RequireRoles(RoleCajero, RolePropietario), tableHandler.Update)
			tables.DELETE("/:id", middleware.RequireRoles(RoleCajero, RolePropietario), tableHandler.Delete)
			tables.POST("/:id/assign", middleware.RequireRoles(RoleCajero, RolePropietario), tableHandler.AssignWaitress)
		}

		// --- INGREDIENTS ---
		ingredients := api.Group("/ingredients")
		ingredients.Use(middleware.RequireRoles(RolePropietario))
		{
			ingredients.GET("", ingredientHandler.GetAll)
			ingredients.GET("/report", ingredientHandler.GetInventoryReport)
			ingredients.POST("", ingredientHandler.Create)
			ingredients.GET("/:id", ingredientHandler.GetByID)
			ingredients.PUT("/:id", ingredientHandler.Update)
			ingredients.DELETE("/:id", ingredientHandler.Delete)
		}

		// --- CATEGORIES ---
		categories := api.Group("/categories")
		categories.Use(middleware.RequireRoles(RolePropietario))
		{
			categories.POST("", productHandler.CreateCategory)
			categories.GET("", productHandler.GetAllCategories)
		}

		// --- PRODUCTS ---
		products := api.Group("/products")
		products.Use(middleware.RequireRoles(RolePropietario))
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.GetAllProducts)
			products.POST("/:id/ingredients", productHandler.AddIngredient)
			products.GET("/:id/ingredients", productHandler.GetIngredients)
		}

		// --- MENU ---
		menu := api.Group("/menu")
		{
			menu.GET("", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), productHandler.GetMenu)
			menu.POST("", middleware.RequireRoles(RolePropietario), productHandler.CreateMenuItem)
			menu.PATCH("/:id", middleware.RequireRoles(RolePropietario), productHandler.UpdateMenuItem)
		}

		// --- ORDERS ---
		orders := api.Group("/orders")
		{
			// Both Mesero & Cajero can create orders
			orders.POST("", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), orderHandler.CreateOrder)
			orders.GET("", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), orderHandler.ListByTable)
			orders.PATCH("/:id/status", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), orderHandler.UpdateOrderStatus)

			// Only Cajero & Propietario can checkout
			checkout := orders.Group("/:id/checkout")
			checkout.Use(middleware.RequireRoles(RoleCajero, RolePropietario))
			checkout.POST("", orderHandler.CheckoutOrder)
		}
	}

}
