package report

import (
	"net/http"
	"time"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uc "github.com/AXONcompany/POS/internal/usecase/report"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *uc.Usecase
}

func NewHandler(usecase *uc.Usecase) *Handler {
	return &Handler{uc: usecase}
}

// GetSalesReport maneja GET /reportes/ventas
func (h *Handler) GetSalesReport(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)

	startStr := c.Query("fecha_inicio")
	endStr := c.Query("fecha_fin")
	reportType := c.DefaultQuery("tipo", "por_dia")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("fecha_inicio y fecha_fin son requeridos", "BAD_REQUEST"))
		return
	}

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Formato de fecha_inicio invalido (YYYY-MM-DD)", "BAD_REQUEST"))
		return
	}

	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Formato de fecha_fin invalido (YYYY-MM-DD)", "BAD_REQUEST"))
		return
	}

	report, err := h.uc.GetSalesReport(c.Request.Context(), venueID.(int), startDate, endDate, reportType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al generar reporte", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(report))
}

// GetInventoryReport maneja GET /reportes/inventario
func (h *Handler) GetInventoryReport(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)
	report, err := h.uc.GetInventoryReport(c.Request.Context(), venueID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al generar reporte", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(report))
}

// GetTipsReport maneja GET /reportes/propinas
func (h *Handler) GetTipsReport(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)

	startStr := c.Query("fecha_inicio")
	endStr := c.Query("fecha_fin")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("fecha_inicio y fecha_fin son requeridos", "BAD_REQUEST"))
		return
	}

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Formato de fecha invalido (YYYY-MM-DD)", "BAD_REQUEST"))
		return
	}

	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Formato de fecha invalido (YYYY-MM-DD)", "BAD_REQUEST"))
		return
	}

	tips, err := h.uc.GetTipsReport(c.Request.Context(), venueID.(int), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al generar reporte", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(tips))
}
