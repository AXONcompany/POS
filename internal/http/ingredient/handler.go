package ingredient

import (
	"net/http"
	"strconv"

	ding "github.com/AXONcompany/POS/internal/domain/ingredient"   //domain ingredient
	uing "github.com/AXONcompany/POS/internal/usecase/ingredients" //use case ingredient
	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	service *uing.IngredientService
}

func NewIngredientHandler(service *uing.IngredientService) *IngredientHandler {
	return &IngredientHandler{service: service}
}

// POST /ingredients
func (h *IngredientHandler) Create(c *gin.Context) {
	var req CreateIngredientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	ing := ding.Ingredient{
		Name:           req.Name,
		UnitOfMeasure:  req.UnitOfMeasure,
		IngredientType: req.Type,
		Stock:          req.Stock,
	}

	created, err := h.service.CreateIngredient(c.Request.Context(), ing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, toIngredientResponse(created))
}

//GET /ingredients/:id

func (h *IngredientHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ing, err := h.service.GetIngredient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toIngredientResponse(ing))
}

// PUT /ingredients/:id
func (h *IngredientHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := FormatValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Convertir DTO a IngredientUpdates del dominio
	updates := ding.IngredientUpdates{
		Name:           req.Name,
		UnitOfMeasure:  req.UnitOfMeasure,
		IngredientType: req.Type,
		Stock:          req.Stock,
	}

	updated, err := h.service.UpdateIngredient(c.Request.Context(), id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toIngredientResponse(updated))
}

// GET /ingredients?page=1&page_size=20
func (h *IngredientHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	ings, err := h.service.GetAllIngredients(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]IngredientResponse, len(ings))
	for i, ing := range ings {
		response[i] = toIngredientResponse(&ing)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      response,
		"page":      page,
		"page_size": pageSize,
	})
}
