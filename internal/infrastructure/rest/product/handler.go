package product

import (
	"net/http"
	"strconv"

	"github.com/AXONcompany/POS/internal/domain/product"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	usecase "github.com/AXONcompany/POS/internal/usecase/product"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *usecase.Usecase
}

func NewHandler(uc *usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}

func prodVenueID(c *gin.Context) int {
	v, _ := c.Get(middleware.VenueIDKey)
	return v.(int)
}

func toCategoryResponse(c *product.Category) CategoryResponse {
	return CategoryResponse{
		ID:         c.ID,
		Name:       c.Name,
		ColorClass: c.ColorClass,
		Icon:       c.Icon,
	}
}

func toProductResponse(p *product.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Price:       p.SalesPrice,
		Description: p.Description,
		CategoryId:  p.CategoryID,
		Image:       p.ImageURL,
		IsAvailable: p.IsActive,
	}
}

// Categories

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": FormatValidationErrors(err)})
		return
	}

	venueID := prodVenueID(c)
	cat := product.Category{
		VenueID:    venueID,
		Name:       req.Name,
		ColorClass: req.ColorClass,
		Icon:       req.Icon,
	}

	created, err := h.uc.CreateCategory(c.Request.Context(), cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toCategoryResponse(created))
}

func (h *Handler) GetAllCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	venueID := prodVenueID(c)
	cats, err := h.uc.GetAllCategories(c.Request.Context(), venueID, page, pageSize)
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

	venueID := prodVenueID(c)
	prod := product.Product{
		VenueID:     venueID,
		Name:        req.Name,
		SalesPrice:  req.Price,
		Description: req.Description,
		CategoryID:  req.CategoryId,
		ImageURL:    req.Image,
		IsActive:    req.IsAvailable,
	}

	created, err := h.uc.CreateProduct(c.Request.Context(), prod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(created))
}

func (h *Handler) GetAllProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	venueID := prodVenueID(c)
	prods, err := h.uc.GetAllProducts(c.Request.Context(), venueID, page, pageSize)
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
