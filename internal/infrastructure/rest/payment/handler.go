package payment

import (
	"net/http"
	"strconv"

	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uc "github.com/AXONcompany/POS/internal/usecase/payment"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *uc.Usecase
}

func NewHandler(usecase *uc.Usecase) *Handler {
	return &Handler{uc: usecase}
}

// --- DTOs ---

type ProcessPaymentRequest struct {
	OrdenID      string              `json:"orden_id" binding:"required"`
	DivisionID   string              `json:"division_id"`
	MetodoPago   string              `json:"metodo_pago" binding:"required"`
	Monto        float64             `json:"monto" binding:"required"`
	Propina      float64             `json:"propina"`
	DetallesPago *PaymentDetailsJSON `json:"detalles_pago"`
}

type PaymentDetailsJSON struct {
	Efectivo          float64 `json:"efectivo"`
	Tarjeta           float64 `json:"tarjeta"`
	ReferenciaTarjeta string  `json:"referencia_tarjeta"`
}

// --- Handlers ---

// ProcessPayment maneja POST /pagos
func (h *Handler) ProcessPayment(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)
	userID, _ := c.Get(middleware.UserIDKey)

	var req ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	// Validar metodo de pago
	switch req.MetodoPago {
	case "efectivo", "tarjeta", "multiple":
	default:
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Metodo de pago invalido", "BAD_REQUEST"))
		return
	}

	orderID, err := strconv.ParseInt(req.OrdenID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("orden_id invalido", "BAD_REQUEST"))
		return
	}

	var divisionID *string
	if req.DivisionID != "" {
		divisionID = &req.DivisionID
	}

	reference := ""
	if req.DetallesPago != nil {
		reference = req.DetallesPago.ReferenciaTarjeta
	}

	payment, err := h.uc.ProcessPayment(
		c.Request.Context(),
		orderID,
		divisionID,
		req.MetodoPago,
		req.Monto,
		req.Propina,
		reference,
		venueID.(int),
		userID.(int),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al procesar pago", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"id":          payment.ID,
		"orden_id":    payment.OrderID,
		"metodo_pago": payment.PaymentMethod,
		"monto":       payment.Amount,
		"propina":     payment.Tip,
		"total":       payment.Total,
		"estado":      payment.Status,
		"referencia":  payment.Reference,
		"fecha":       payment.CreatedAt,
	}))
}

// GetInvoice maneja GET /pagos/:id/factura
func (h *Handler) GetInvoice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	invoice, err := h.uc.GenerateInvoice(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Pago no encontrado", "NOT_FOUND"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(invoice))
}
