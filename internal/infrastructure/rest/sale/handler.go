package sale

import (
	"fmt"
	"net/http"
	"strconv"

	usale "github.com/AXONcompany/POS/internal/usecase/sales"
	"github.com/gin-gonic/gin"
)

type SaleHandler struct {
	uc *usale.Usecase
}

func NewSaleHandler(uc *usale.Usecase) *SaleHandler {
	return &SaleHandler{uc: uc}
}

// POST /payments
func (h *SaleHandler) ProcessPayment(c *gin.Context) {
	var req ProcessPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sale, err := h.uc.ProcessPayment(c.Request.Context(), req.OrderID, req.RestaurantID, req.Total, req.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toSaleResponse(sale))
}

// GET /payments/:id/invoice
func (h *SaleHandler) GetInvoice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sale, err := h.uc.GetInvoice(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toSaleResponse(sale))
}

// POST /orders/:id/split
func (h *SaleHandler) SplitOrder(c *gin.Context) {
	var req SplitOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.uc.SplitOrder(c.Request.Context(), req.Total, req.People)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, result)
}