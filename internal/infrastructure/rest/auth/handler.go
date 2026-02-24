package auth

import (
	"net/http"
	"time"

	"github.com/AXONcompany/POS/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *auth.Usecase
}

func NewHandler(uc *auth.Usecase) *Handler {
	return &Handler{uc: uc}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ipAddress := c.ClientIP()
	deviceInfo := c.GetHeader("User-Agent")

	tokens, err := h.uc.Login(c.Request.Context(), req.Email, req.Password, deviceInfo, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set refresh token as HttpOnly Cookie
	c.SetCookie("refresh_token", tokens.RefreshToken, int(7*24*time.Hour.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
	})
}
