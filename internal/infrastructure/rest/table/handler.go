package table

import (
	"net/http"
	"strconv"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	usecase "github.com/AXONcompany/POS/internal/usecase/table"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *usecase.Usecase
}

func NewHandler(u *usecase.Usecase) *Handler {
	return &Handler{uc: u}
}

func tblVenueID(c *gin.Context) int {
	v, _ := c.Get(middleware.VenueIDKey)
	return v.(int)
}

// Create maneja POST /mesas
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos: "+err.Error(), "BAD_REQUEST"))
		return
	}

	venueID := tblVenueID(c)
	domainTable := ToDomain(req)
	domainTable.VenueID = venueID

	if err := h.uc.Create(c.Request.Context(), domainTable); err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusCreated, httputil.SuccessResponse(ToResponse(domainTable)))
}

// GetAll maneja GET /mesas
func (h *Handler) GetAll(c *gin.Context) {
	venueID := tblVenueID(c)
	tables, err := h.uc.FindAll(c.Request.Context(), venueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(ToResponseList(tables)))
}

// GetByID maneja GET /mesas/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	venueID := tblVenueID(c)
	t, err := h.uc.FindByID(c.Request.Context(), id, venueID)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Mesa no encontrada", "NOT_FOUND"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(ToResponse(t)))
}

// UpdateEstado maneja PATCH /mesas/:id/estado
func (h *Handler) UpdateEstado(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	var req UpdateEstadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse(err.Error(), "BAD_REQUEST"))
		return
	}

	venueID := tblVenueID(c)
	updates := ToUpdateDomain(UpdateRequest{Status: &req.Estado})

	if err := h.uc.Update(c.Request.Context(), id, venueID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	t, err := h.uc.FindByID(c.Request.Context(), id, venueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(ToResponse(t)))
}

// Delete maneja DELETE /mesas/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	venueID := tblVenueID(c)
	if err := h.uc.Delete(c.Request.Context(), id, venueID); err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessMessageResponse("Mesa eliminada"))
}
