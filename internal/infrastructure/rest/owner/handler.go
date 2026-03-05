package owner

import (
	"net/http"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uc "github.com/AXONcompany/POS/internal/usecase/owner"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *uc.Usecase
}

func NewHandler(usecase *uc.Usecase) *Handler {
	return &Handler{uc: usecase}
}

type CreateOwnerRequest struct {
	Name     string `json:"nombre" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateOwnerRequest struct {
	Name  string `json:"nombre"`
	Email string `json:"email"`
}

// GetMe devuelve la info del owner autenticado.
func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse("No autorizado", "UNAUTHORIZED"))
		return
	}

	o, err := h.uc.GetOwnerByID(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Propietario no encontrado", "NOT_FOUND"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"id":        o.ID,
		"nombre":    o.Name,
		"email":     o.Email,
		"activo":    o.IsActive,
		"creado_en": o.CreatedAt,
	}))
}

// Update actualiza datos del owner autenticado.
func (h *Handler) Update(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse("No autorizado", "UNAUTHORIZED"))
		return
	}

	var req UpdateOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	updated, err := h.uc.UpdateOwner(c.Request.Context(), userID.(int), req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"id":        updated.ID,
		"nombre":    updated.Name,
		"email":     updated.Email,
		"activo":    updated.IsActive,
		"creado_en": updated.CreatedAt,
	}))
}
