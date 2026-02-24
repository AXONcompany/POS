package product

import (
	"net/http"
	"strconv"

	"github.com/AXONcompany/POS/internal/domain/product"
	"github.com/gin-gonic/gin"
)

// Menu Handlers

func (h *Handler) GetMenu(c *gin.Context) {
	// For now, return all products. In future we might enrich this.
	// Reusing GetAllProducts logic but specialized for menu view if needed.
	// User req: "GET /menu - obtiene el menu"
	h.GetAllProducts(c)
}

func (h *Handler) CreateMenuItem(c *gin.Context) {
	var req CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationErrors(err)})
		return
	}

	// Map DTO to Domain
	ingredients := make([]product.RecipeItem, len(req.Ingredients))
	for i, ing := range req.Ingredients {
		ingredients[i] = product.RecipeItem{
			IngredientID:     ing.IngredientID,
			QuantityRequired: ing.Quantity,
		}
	}

	created, err := h.uc.CreateMenuItem(c.Request.Context(), req.Name, req.SalesPrice, ingredients)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(created))
}

func (h *Handler) UpdateMenuItem(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateMenuItemRequest
	// Relaxed validation for PATCH
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationErrors(err)})
		return
	}

	// Implementation Note: Since we don't have atomic PATCH in service yet, we'll do:
	// 1. GetProduct
	// 2. Update fields
	// 3. UpdateProduct

	current, err := h.uc.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	if req.Name != "" {
		current.Name = req.Name
	}
	if req.SalesPrice > 0 { // Assuming price update if provided
		current.SalesPrice = req.SalesPrice
	}

	updated, err := h.uc.UpdateProduct(c.Request.Context(), id, *current)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toProductResponse(updated))
}
