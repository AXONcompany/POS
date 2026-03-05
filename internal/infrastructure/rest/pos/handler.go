package pos

import (
	"net/http"
	"strconv"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uc "github.com/AXONcompany/POS/internal/usecase/pos"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *uc.Usecase
}

func NewHandler(usecase *uc.Usecase) *Handler {
	return &Handler{uc: usecase}
}

type CreateTerminalRequest struct {
	Name string `json:"nombre" binding:"required"`
}

type UpdateTerminalRequest struct {
	Name     string `json:"nombre"`
	IsActive *bool  `json:"activo"`
}

// Create crea un nuevo terminal POS en la sede del usuario autenticado.
func (h *Handler) Create(c *gin.Context) {
	var req CreateTerminalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	venueID, _ := c.Get(middleware.VenueIDKey)
	created, err := h.uc.CreateTerminal(c.Request.Context(), venueID.(int), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusCreated, httputil.SuccessResponse(gin.H{
		"id":       created.ID,
		"nombre":   created.TerminalName,
		"venue_id": created.VenueID,
		"activo":   created.IsActive,
	}))
}

// List lista terminales de la sede del usuario autenticado.
func (h *Handler) List(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)
	terminals, err := h.uc.ListTerminalsByVenue(c.Request.Context(), venueID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	results := make([]gin.H, len(terminals))
	for i, t := range terminals {
		results[i] = gin.H{
			"id":       t.ID,
			"nombre":   t.TerminalName,
			"venue_id": t.VenueID,
			"activo":   t.IsActive,
		}
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(results))
}

// GetByID obtiene un terminal por ID.
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	t, err := h.uc.GetTerminalByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Terminal no encontrado", "NOT_FOUND"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"id":       t.ID,
		"nombre":   t.TerminalName,
		"venue_id": t.VenueID,
		"activo":   t.IsActive,
	}))
}

// Update actualiza un terminal POS.
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	var req UpdateTerminalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	updated, err := h.uc.UpdateTerminal(c.Request.Context(), id, req.Name, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"id":       updated.ID,
		"nombre":   updated.TerminalName,
		"venue_id": updated.VenueID,
		"activo":   updated.IsActive,
	}))
}
