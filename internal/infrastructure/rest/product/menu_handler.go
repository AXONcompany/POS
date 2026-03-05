package product

import (
	"net/http"
	"strconv"

	"github.com/AXONcompany/POS/internal/domain/product"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/gin-gonic/gin"
)

// Menu Handlers

func (h *Handler) GetMenu(c *gin.Context) {
	venueID := prodVenueID(c)
	products, err := h.uc.GetAllProducts(c.Request.Context(), venueID, 1, 1000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	result := make([]gin.H, len(products))
	for i, p := range products {
		result[i] = toMenuItemResponse(&p)
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(result))
}

func toMenuItemResponse(p *product.Product) gin.H {
	return gin.H{
		"id":         p.ID,
		"nombre":     p.Name,
		"precio":     p.SalesPrice,
		"disponible": p.IsActive,
	}
}

func (h *Handler) CreateMenuItem(c *gin.Context) {
	var req CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationErrors(err)})
		return
	}

	venueID := prodVenueID(c)
	ingredients := make([]product.RecipeItem, len(req.Ingredients))
	for i, ing := range req.Ingredients {
		ingredients[i] = product.RecipeItem{
			IngredientID:     ing.IngredientID,
			QuantityRequired: ing.Quantity,
		}
	}

	created, err := h.uc.CreateMenuItem(c.Request.Context(), venueID, req.Name, req.SalesPrice, ingredients)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusCreated, httputil.SuccessResponse(toMenuItemResponse(created)))
}

func (h *Handler) UpdateMenuItem(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID invalido", "BAD_REQUEST"))
		return
	}

	var req UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationErrors(err)})
		return
	}

	venueID := prodVenueID(c)
	current, err := h.uc.GetProduct(c.Request.Context(), id, venueID)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Item no encontrado", "NOT_FOUND"))
		return
	}

	if req.Name != "" {
		current.Name = req.Name
	}
	if req.SalesPrice > 0 {
		current.SalesPrice = req.SalesPrice
	}

	updated, err := h.uc.UpdateProduct(c.Request.Context(), id, venueID, *current)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse(err.Error(), "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(toMenuItemResponse(updated)))
}
