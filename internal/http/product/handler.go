package product

import (
	"net/http"
	"strconv"

	"github.com/AXONcompany/POS/internal/domain/product"
	usecase "github.com/AXONcompany/POS/internal/usecase/products"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *usecase.Service
}

func NewHandler(service *usecase.Service) *Handler {
	return &Handler{service: service}
}

func toCategoryResponse(c *product.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toProductResponse(p *product.Product) ProductResponse {
	return ProductResponse{
		ID:         p.ID,
		Name:       p.Name,
		SalesPrice: p.SalesPrice,
		IsActive:   p.IsActive,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

// Categories

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationErrors(err)})
		return
	}

	cat := product.Category{
		Name: req.Name,
	}

	created, err := h.service.CreateCategory(c.Request.Context(), cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toCategoryResponse(created))
}

func (h *Handler) GetAllCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	cats, err := h.service.GetAllCategories(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]CategoryResponse, len(cats))
	for i, cat := range cats {
		response[i] = toCategoryResponse(&cat)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      response,
		"page":      page,
		"page_size": pageSize,
	})
}

// Products

func (h *Handler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationErrors(err)})
		return
	}

	prod := product.Product{
		Name:       req.Name,
		SalesPrice: req.SalesPrice,
		IsActive:   req.IsActive,
	}

	created, err := h.service.CreateProduct(c.Request.Context(), prod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(created))
}

func (h *Handler) GetAllProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	prods, err := h.service.GetAllProducts(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]ProductResponse, len(prods))
	for i, prod := range prods {
		response[i] = toProductResponse(&prod)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      response,
		"page":      page,
		"page_size": pageSize,
	})
}

// Recipes

func (h *Handler) AddIngredient(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req AddIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationErrors(err)})
		return
	}

	item, err := h.service.AddIngredient(c.Request.Context(), productID, req.IngredientID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, RecipeItemResponse{
		ID:           item.ID,
		ProductID:    item.ProductID,
		IngredientID: item.IngredientID,
		Quantity:     item.QuantityRequired,
	})
}

func (h *Handler) GetIngredients(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	items, err := h.service.GetProductIngredients(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]RecipeItemResponse, len(items))
	for i, item := range items {
		response[i] = RecipeItemResponse{
			ID:             item.ID,
			ProductID:      item.ProductID,
			IngredientID:   item.IngredientID,
			IngredientName: item.IngredientName,
			UnitOfMeasure:  item.UnitOfMeasure,
			Quantity:       item.QuantityRequired,
		}
	}

	c.JSON(http.StatusOK, response)
}
