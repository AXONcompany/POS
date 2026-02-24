package order

import (
	"net/http"
	"strconv"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uc "github.com/AXONcompany/POS/internal/usecase/order"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *uc.Usecase
}

func NewHandler(usecase *uc.Usecase) *Handler {
	return &Handler{uc: usecase}
}

type CreateOrderRequest struct {
	TableID *int `json:"table_id"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	restaurantID, _ := c.Get(middleware.RestaurantIDKey)
	userID, _ := c.Get(middleware.UserIDKey)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	order, err := h.uc.CreateOrder(c.Request.Context(), restaurantID.(int), userID.(int), req.TableID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) CheckoutOrder(c *gin.Context) {
	restaurantID, _ := c.Get(middleware.RestaurantIDKey)
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	err = h.uc.CheckoutOrder(c.Request.Context(), restaurantID.(int), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to checkout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "PAID"})
}
