package ingredient

import (
	"net/http"

	"github.com/AXONcompany/POS/internal/domain/ingredient"
	"github.com/AXONcompany/POS/internal/usecase"
	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	service *usecase.IngredientService
}

func NewIngredientHandler(service *usecase.IngredientService) *IngredientHandler {
	return &IngredientHandler{service: service}
}

// POST /ingredients
func (h *IngredientHandler) Create(c *gin.Context) {
	var req CreateIngredientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	ing := ingredient.Ingredient{
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
	panic("to do")
}
