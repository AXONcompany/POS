package rest

import (
	"log"
	"net/http"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/auth"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/ingredient"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	orderrest "github.com/AXONcompany/POS/internal/infrastructure/rest/order"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/owner"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/payment"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/pos"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/product"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/report"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/table"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/user"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/venue"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine, ingredientHandler *ingredient.IngredientHandler, productHandler *product.Handler, authHandler *auth.Handler, orderHandler *orderrest.Handler, tableHandler *table.Handler, userHandler *user.Handler, paymentHandler *payment.Handler, reportHandler *report.Handler, ownerHandler *owner.Handler, venueHandler *venue.Handler, posHandler *pos.Handler, jwtSecret []byte) {

	log.Printf("RegisterRouters called, ingredientHandler is nil: %v", ingredientHandler == nil)

	// Health checks
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "server say: pong")
	})

	// Roles mappings
	const RolePropietario = 1
	const RoleCajero = 2
	const RoleMesero = 3

	// --- AUTH (publico) ---
	authPublic := r.Group("/auth")
	{
		authPublic.POST("/login", authHandler.Login)
		authPublic.POST("/register-owner", authHandler.RegisterOwner)
		authPublic.POST("/refresh", authHandler.Refresh)
	}

	// --- AUTH (protegido) ---
	authProtected := r.Group("/auth")
	authProtected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		authProtected.POST("/register", middleware.RequireRoles(RolePropietario, RoleCajero), authHandler.Register)
		authProtected.GET("/me", authHandler.Me)
		authProtected.POST("/logout", authHandler.Logout)
		authProtected.POST("/switch-sede", middleware.RequireRoles(RolePropietario), authHandler.SwitchSede)
	}

	// Protected API Group
	api := r.Group("")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// --- PROPIETARIO ---
		propietario := api.Group("/propietario")
		propietario.Use(middleware.RequireRoles(RolePropietario))
		{
			propietario.GET("", ownerHandler.GetMe)
			propietario.PATCH("", ownerHandler.Update)
		}

		// --- SEDES ---
		sedes := api.Group("/sedes")
		{
			sedes.GET("/mi-sede", middleware.RequireRoles(RoleCajero, RoleMesero, RolePropietario), venueHandler.GetMyVenue)
			sedes.POST("", middleware.RequireRoles(RolePropietario), venueHandler.Create)
			sedes.GET("", middleware.RequireRoles(RolePropietario), venueHandler.List)
			sedes.GET("/:id", middleware.RequireRoles(RolePropietario), venueHandler.GetByID)
			sedes.PATCH("/:id", middleware.RequireRoles(RolePropietario), venueHandler.Update)
		}

		// --- TERMINALES POS ---
		terminales := api.Group("/terminales")
		{
			terminales.POST("", middleware.RequireRoles(RolePropietario), posHandler.Create)
			terminales.GET("", middleware.RequireRoles(RolePropietario, RoleCajero), posHandler.List)
			terminales.GET("/:id", middleware.RequireRoles(RolePropietario, RoleCajero), posHandler.GetByID)
			terminales.PATCH("/:id", middleware.RequireRoles(RolePropietario), posHandler.Update)
		}

		// --- MESAS ---
		mesas := api.Group("/mesas")
		{
			mesas.GET("", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), tableHandler.GetAll)
			mesas.GET("/:id", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), tableHandler.GetByID)

			mesas.POST("", middleware.RequireRoles(RoleCajero, RolePropietario), tableHandler.Create)
			mesas.PATCH("/:id/estado", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), tableHandler.UpdateEstado)
			mesas.DELETE("/:id", middleware.RequireRoles(RoleCajero, RolePropietario), tableHandler.Delete)
			mesas.POST("/:id/asignar", middleware.RequireRoles(RoleCajero, RolePropietario), tableHandler.AssignWaiter)
			mesas.GET("/:id/asignaciones", middleware.RequireRoles(RoleCajero, RolePropietario), tableHandler.GetAssignments)
		}

		// --- INGREDIENTES ---
		ingredientes := api.Group("/ingredientes")
		ingredientes.Use(middleware.RequireRoles(RolePropietario))
		{
			ingredientes.GET("", ingredientHandler.GetAll)
			ingredientes.GET("/report", ingredientHandler.GetInventoryReport)
			ingredientes.POST("", ingredientHandler.Create)
			ingredientes.GET("/:id", ingredientHandler.GetByID)
			ingredientes.PUT("/:id", ingredientHandler.Update)
			ingredientes.PATCH("/:id/stock", ingredientHandler.UpdateStock)
			ingredientes.DELETE("/:id", ingredientHandler.Delete)
		}

		// --- CATEGORIAS ---
		categorias := api.Group("/categorias")
		categorias.Use(middleware.RequireRoles(RolePropietario))
		{
			categorias.POST("", productHandler.CreateCategory)
			categorias.GET("", productHandler.GetAllCategories)
		}

		// --- PRODUCTS (interno, mantener compatibilidad) ---
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

		// --- ORDENES ---
		ordenes := api.Group("/ordenes")
		{
			ordenes.POST("", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), orderHandler.CreateOrder)
			ordenes.GET("", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), orderHandler.ListByTable)
			ordenes.GET("/:id", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), orderHandler.GetByID)
			ordenes.POST("/:id/items", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), orderHandler.AddItems)
			ordenes.DELETE("/:id/items/:item_id", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), orderHandler.CancelItem)
			ordenes.POST("/:id/enviar-cocina", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), orderHandler.SendToKitchen)
			ordenes.PATCH("/:id/status", middleware.RequireRoles(RoleMesero, RoleCajero, RolePropietario), orderHandler.UpdateOrderStatus)
			ordenes.POST("/:id/dividir", middleware.RequireRoles(RoleCajero, RolePropietario), orderHandler.DivideOrder)

			// Checkout
			checkout := ordenes.Group("/:id/checkout")
			checkout.Use(middleware.RequireRoles(RoleCajero, RolePropietario))
			checkout.POST("", orderHandler.CheckoutOrder)
		}

		// --- USUARIOS ---
		usuarios := api.Group("/usuarios")
		{
			usuarios.GET("", middleware.RequireRoles(RolePropietario, RoleCajero), userHandler.GetAll)
			usuarios.GET("/:id", middleware.RequireRoles(RolePropietario), userHandler.GetByID)
			usuarios.PATCH("/:id", middleware.RequireRoles(RolePropietario), userHandler.Update)
			usuarios.DELETE("/:id", middleware.RequireRoles(RolePropietario), userHandler.Delete)
			usuarios.POST("/mesero", middleware.RequireRoles(RolePropietario, RoleCajero), authHandler.RegisterWaiter)
		}

		// --- PAGOS ---
		pagos := api.Group("/pagos")
		pagos.Use(middleware.RequireRoles(RoleCajero, RolePropietario))
		{
			pagos.POST("", paymentHandler.ProcessPayment)
			pagos.GET("/:id/factura", paymentHandler.GetInvoice)
		}

		// --- REPORTES ---
		reportes := api.Group("/reportes")
		reportes.Use(middleware.RequireRoles(RolePropietario))
		{
			reportes.GET("/ventas", reportHandler.GetSalesReport)
			reportes.GET("/inventario", reportHandler.GetInventoryReport)
			reportes.GET("/propinas", reportHandler.GetTipsReport)
		}
	}

}
