package venue

import (
	"net/http"
	"strconv"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uc "github.com/AXONcompany/POS/internal/usecase/venue"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *uc.Usecase
}

func NewHandler(usecase *uc.Usecase) *Handler {
	return &Handler{uc: usecase}
}

type CreateVenueRequest struct {
	Name    string `json:"nombre" binding:"required"`
	Address string `json:"direccion"`
	Phone   string `json:"telefono"`
}

type UpdateVenueRequest struct {
	Name    string `json:"nombre"`
	Address string `json:"direccion"`
	Phone   string `json:"telefono"`
}

// Create crea una nueva sede para el owner autenticado.
func (h *Handler) Create(c *gin.Context) {
	var req CreateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	userID, _ := c.Get(middleware.UserIDKey)
	created, err := h.uc.CreateVenue(c.Request.Context(), userID.(int), req.Name, req.Address, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusCreated, httputil.SuccessResponse(gin.H{
		"id":        created.ID,
		"nombre":    created.Name,
		"direccion": created.Address,
		"telefono":  created.Phone,
		"activo":    created.IsActive,
		"creado_en": created.CreatedAt,
	}))
}

// List lista sedes del owner autenticado.
func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	venues, err := h.uc.ListVenuesByOwner(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	results := make([]gin.H, len(venues))
	for i, v := range venues {
		results[i] = gin.H{
			"id":        v.ID,
			"nombre":    v.Name,
			"direccion": v.Address,
			"telefono":  v.Phone,
			"activo":    v.IsActive,
			"creado_en": v.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(results))
}

// GetByID obtiene una sede por ID.
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	v, err := h.uc.GetVenueByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Sede no encontrada", "NOT_FOUND"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"id":        v.ID,
		"nombre":    v.Name,
		"direccion": v.Address,
		"telefono":  v.Phone,
		"activo":    v.IsActive,
		"creado_en": v.CreatedAt,
	}))
}

// Update actualiza una sede.
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	var req UpdateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	updated, err := h.uc.UpdateVenue(c.Request.Context(), id, req.Name, req.Address, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"id":        updated.ID,
		"nombre":    updated.Name,
		"direccion": updated.Address,
		"telefono":  updated.Phone,
		"activo":    updated.IsActive,
	}))
}
