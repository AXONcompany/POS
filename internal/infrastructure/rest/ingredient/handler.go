package ingredient

import (
	"net/http"
	"strconv"

	ding "github.com/AXONcompany/POS/internal/domain/ingredient"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uing "github.com/AXONcompany/POS/internal/usecase/ingredient"
	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	uc *uing.Usecase
}

func NewIngredientHandler(uc *uing.Usecase) *IngredientHandler {
	return &IngredientHandler{uc: uc}
}

func getVenueID(c *gin.Context) int {
	v, _ := c.Get(middleware.VenueIDKey)
	return v.(int)
}

// POST /ingredientes
func (h *IngredientHandler) Create(c *gin.Context) {
	var req CreateIngredientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	venueID := getVenueID(c)

	ing := ding.Ingredient{
		VenueID:        venueID,
		Name:           req.Name,
		UnitOfMeasure:  req.UnitOfMeasure,
		IngredientType: req.Type,
		Stock:          req.Stock,
	}

	created, err := h.uc.CreateIngredient(c.Request.Context(), ing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusCreated, httputil.SuccessResponse(toIngredientResponse(created)))
}

// GET /ingredientes/:id
func (h *IngredientHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	venueID := getVenueID(c)
	ing, err := h.uc.GetIngredient(c.Request.Context(), id, venueID)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse(err.Error(), "NOT_FOUND"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(toIngredientResponse(ing)))
}

// PUT /ingredientes/:id
func (h *IngredientHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	var req UpdateIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	venueID := getVenueID(c)
	updates := ding.PartialIngredient{
		Name:           req.Name,
		UnitOfMeasure:  req.UnitOfMeasure,
		IngredientType: req.Type,
		Stock:          req.Stock,
	}

	updated, err := h.uc.UpdateIngredient(c.Request.Context(), id, venueID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(toIngredientResponse(updated)))
}

// GET /ingredientes
func (h *IngredientHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	venueID := getVenueID(c)
	ings, err := h.uc.GetAllIngredients(c.Request.Context(), venueID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	response := make([]IngredientResponse, len(ings))
	for i, ing := range ings {
		response[i] = toIngredientResponse(&ing)
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(response))
}

func (h *IngredientHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	venueID := getVenueID(c)
	err = h.uc.DeleteIngredient(c.Request.Context(), id, venueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *IngredientHandler) GetInventoryReport(c *gin.Context) {
	venueID := getVenueID(c)
	ings, err := h.uc.GetInventoryReport(c.Request.Context(), venueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	response := make([]IngredientResponse, len(ings))
	for i, ing := range ings {
		response[i] = toIngredientResponse(&ing)
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(gin.H{
		"data":  response,
		"count": len(response),
	}))
}

// UpdateStock maneja PATCH /ingredientes/:id/stock
func (h *IngredientHandler) UpdateStock(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	var req UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	venueID := getVenueID(c)
	// Obtener ingrediente actual
	ing, err := h.uc.GetIngredient(c.Request.Context(), id, venueID)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Ingrediente no encontrado", "NOT_FOUND"))
		return
	}

	// Aplicar movimiento de stock
	newStock := ing.Stock
	switch req.TipoMovimiento {
	case "entrada":
		newStock += req.Cantidad
	case "salida":
		newStock -= req.Cantidad
		if newStock < 0 {
			c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Stock insuficiente", "BAD_REQUEST"))
			return
		}
	default:
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("tipo_movimiento debe ser 'entrada' o 'salida'", "BAD_REQUEST"))
		return
	}

	updates := ding.PartialIngredient{
		Stock: &newStock,
	}

	updated, err := h.uc.UpdateIngredient(c.Request.Context(), id, venueID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(toIngredientResponse(updated)))
}
