package user

import (
	"log"
	"net/http"

	"github.com/AXONcompany/POS/internal/domain/user"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uc "github.com/AXONcompany/POS/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *uc.Usecase
}

func NewHandler(usecase *uc.Usecase) *Handler {
	return &Handler{uc: usecase}
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   int    `json:"role_id" binding:"required"`
}

func (h *Handler) CreateEmployee(c *gin.Context) {
	restaurantID, exists := c.Get(middleware.RestaurantIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	u := &user.User{
		RestaurantID: restaurantID.(int),
		RoleID:       req.RoleID,
		Name:         req.Name,
		Email:        req.Email,
		IsActive:     true,
	}

	createdUser, err := h.uc.CreateUser(c.Request.Context(), u, req.Password)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}
