package table

import (
	"net/http"
	"strconv"

	usecase "github.com/AXONcompany/POS/internal/usecase/table"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *usecase.Usecase
}

func NewHandler(u *usecase.Usecase) *Handler {
	return &Handler{uc: u}
}

// Create maneja POST /tables
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	domainTable := ToDomain(req)

	if err := h.uc.Create(c.Request.Context(), domainTable); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ToResponse(domainTable))
}

// GetAll maneja GET /tables
func (h *Handler) GetAll(c *gin.Context) {
	tables, err := h.uc.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToResponseList(tables))
}

// GetByID maneja GET /tables/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	t, err := h.uc.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mesa no encontrada"})
		return
	}

	c.JSON(http.StatusOK, ToResponse(t))
}

// Update maneja PATCH /tables/:id
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := ToUpdateDomain(req)

	if err := h.uc.Update(c.Request.Context(), id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesa actualizada correctamente"})
}

// Delete maneja DELETE /tables/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesa eliminada"})
}

// AssignWaitress maneja POST /tables/:id/assign
func (h *Handler) AssignWaitress(c *gin.Context) {
	tableID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de mesa inválido"})
		return
	}

	var req AssignWaitressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.AssignWaitress(c.Request.Context(), tableID, req.WaitressID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesero asignado exitosamente"})
}
