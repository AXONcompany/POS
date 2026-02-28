package order

import (
	"net/http"
	"strconv"

	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
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
	TableID *int64                  `json:"table_id"`
	Items   []domainOrder.OrderItem `json:"items" binding:"required"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	restaurantID, _ := c.Get(middleware.RestaurantIDKey)
	userID, _ := c.Get(middleware.UserIDKey)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	order, err := h.uc.CreateOrder(c.Request.Context(), restaurantID.(int), int(userID.(float64)), req.TableID, req.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) CheckoutOrder(c *gin.Context) {
	restaurantID, _ := c.Get(middleware.RestaurantIDKey)
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
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

type UpdateOrderStatusRequest struct {
	StatusID int `json:"status_id" binding:"required"`
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	restaurantID, _ := c.Get(middleware.RestaurantIDKey)
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = h.uc.UpdateOrderStatus(c.Request.Context(), restaurantID.(int), orderID, req.StatusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}

func (h *Handler) ListByTable(c *gin.Context) {
	restaurantID, _ := c.Get(middleware.RestaurantIDKey)
	tableIDStr := c.Query("table_id")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table id"})
		return
	}

	orders, err := h.uc.ListOrdersByTable(c.Request.Context(), restaurantID.(int), tableID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
