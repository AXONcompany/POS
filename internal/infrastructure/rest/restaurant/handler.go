package restaurant

import (
	"net/http"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uc "github.com/AXONcompany/POS/internal/usecase/restaurant"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *uc.Usecase
}

func NewHandler(usecase *uc.Usecase) *Handler {
	return &Handler{uc: usecase}
}

func (h *Handler) GetMyRestaurant(c *gin.Context) {
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rest, err := h.uc.GetMyRestaurant(c.Request.Context(), restaurantID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch restaurant"})
		return
	}

	c.JSON(http.StatusOK, rest)
}

type UpdateRestaurantRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func (h *Handler) UpdateMyRestaurant(c *gin.Context) {
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	rest, err := h.uc.UpdateRestaurantInfo(c.Request.Context(), restaurantID.(int), req.Name, req.Address, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update restaurant"})
		return
	}

	c.JSON(http.StatusOK, rest)
}
