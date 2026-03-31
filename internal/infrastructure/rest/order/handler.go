package order

import (
	"errors"
	"net/http"
	"strconv"

	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/httputil"
	"github.com/AXONcompany/POS/internal/infrastructure/rest/middleware"
	uc "github.com/AXONcompany/POS/internal/usecase/order"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *uc.Usecase
}

func NewHandler(usecase *uc.Usecase) *Handler {
	return &Handler{uc: usecase}
}

// --- DTOs ---

type CreateOrderRequest struct {
	MesaID   string `json:"mesa_id"`
	MeseroID string `json:"mesero_id"`
}

type AddItemsRequest struct {
	Items []AddItemRequest `json:"items" binding:"required"`
}

type AddItemRequest struct {
	MenuItemID string `json:"menu_item_id" binding:"required"`
	Cantidad   int    `json:"cantidad" binding:"required,gt=0"`
	Notas      string `json:"notas"`
}

// --- Mapper ---

func toOrdenResponse(o *domainOrder.Order) gin.H {
	// Mapear estado por status ID
	estados := map[int]string{
		1: "abierta",
		2: "enviada",
		3: "en_preparacion",
		4: "lista",
		5: "pagada",
		6: "cancelada",
	}

	estado := estados[o.StatusID]
	if estado == "" {
		estado = o.Status
	}

	items := make([]gin.H, len(o.Items))
	for i, item := range o.Items {
		items[i] = gin.H{
			"id":              item.ID,
			"menu_item_id":    item.ProductID,
			"cantidad":        item.Quantity,
			"precio_unitario": item.UnitPrice,
			"notas":           item.Notes,
		}
	}

	subtotal := o.TotalAmount
	impuestos := subtotal * 0.19
	total := subtotal + impuestos

	return gin.H{
		"id":             o.ID,
		"mesa_id":        o.TableID,
		"mesero_id":      o.UserID,
		"estado":         estado,
		"items":          items,
		"subtotal":       subtotal,
		"impuestos":      impuestos,
		"total":          total,
		"fecha_creacion": o.CreatedAt,
	}
}

// --- Handlers ---

func (h *Handler) CreateOrder(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)
	userID, _ := c.Get(middleware.UserIDKey)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	// Parsear mesa_id como int64 si es posible
	var tableID *int64
	if req.MesaID != "" {
		tid, err := strconv.ParseInt(req.MesaID, 10, 64)
		if err == nil {
			tableID = &tid
		}
	}

	// Crear orden sin items (segun swagger, items no son obligatorios al crear)
	order, err := h.uc.CreateOrderWithoutItems(c.Request.Context(), venueID.(int), userID.(int), tableID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al crear orden", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusCreated, httputil.SuccessResponse(toOrdenResponse(order)))
}

func (h *Handler) GetByID(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID de orden invalido", "BAD_REQUEST"))
		return
	}

	order, err := h.uc.GetOrderByID(c.Request.Context(), venueID.(int), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Orden no encontrada", "NOT_FOUND"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(toOrdenResponse(order)))
}

func (h *Handler) AddItems(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID de orden invalido", "BAD_REQUEST"))
		return
	}

	var req AddItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	// Convertir request a domain items
	items := make([]domainOrder.OrderItem, len(req.Items))
	for i, item := range req.Items {
		productID, _ := strconv.ParseInt(item.MenuItemID, 10, 64)
		items[i] = domainOrder.OrderItem{
			ProductID: productID,
			Quantity:  item.Cantidad,
			Notes:     item.Notas,
		}
	}

	err = h.uc.AddProductToOrder(c.Request.Context(), venueID.(int), orderID, items)
	if err != nil {
		if errors.Is(err, domainOrder.ErrInvalidStatusTransition) {
			c.JSON(http.StatusUnprocessableEntity, httputil.ErrorResponse("La orden no acepta mas items en su estado actual", "INVALID_TRANSITION"))
			return
		}
		if errors.Is(err, domainOrder.ErrInsufficientStock) {
			c.JSON(http.StatusConflict, httputil.ErrorResponse("Stock insuficiente para uno o mas ingredientes", "INSUFFICIENT_STOCK"))
			return
		}
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al agregar items", "INTERNAL_ERROR"))
		return
	}

	// Re-fetch order para devolver estado actualizado
	order, err := h.uc.GetOrderByID(c.Request.Context(), venueID.(int), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al obtener orden", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(toOrdenResponse(order)))
}

func (h *Handler) CancelItem(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID de orden invalido", "BAD_REQUEST"))
		return
	}

	itemID, err := strconv.ParseInt(c.Param("item_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID de item invalido", "BAD_REQUEST"))
		return
	}

	venueID, _ := c.Get(middleware.VenueIDKey)
	userID, _ := c.Get(middleware.UserIDKey)

	err = h.uc.CancelOrderItem(c.Request.Context(), venueID.(int), userID.(int), orderID, itemID)
	if err != nil {
		if errors.Is(err, domainOrder.ErrItemAlreadyCancelled) {
			c.JSON(http.StatusConflict, httputil.ErrorResponse("El item ya fue cancelado", "ITEM_ALREADY_CANCELLED"))
			return
		}
		if errors.Is(err, domainOrder.ErrInvalidStatusTransition) {
			c.JSON(http.StatusUnprocessableEntity, httputil.ErrorResponse("No se pueden cancelar items de una orden en este estado", "INVALID_TRANSITION"))
			return
		}
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al cancelar item", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessMessageResponse("Item cancelado exitosamente"))
}

func (h *Handler) SendToKitchen(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID de orden invalido", "BAD_REQUEST"))
		return
	}

	venueID, _ := c.Get(middleware.VenueIDKey)

	// Cambiar estado a "enviada" (status_id = 2)
	err = h.uc.UpdateOrderStatus(c.Request.Context(), venueID.(int), orderID, 2)
	if err != nil {
		if errors.Is(err, domainOrder.ErrInvalidStatusTransition) {
			c.JSON(http.StatusUnprocessableEntity, httputil.ErrorResponse("Transicion de estado invalida", "INVALID_TRANSITION"))
			return
		}
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al enviar a cocina", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessMessageResponse("Orden enviada a cocina y bar"))
}

func (h *Handler) CheckoutOrder(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID de orden invalido", "BAD_REQUEST"))
		return
	}

	err = h.uc.CheckoutOrder(c.Request.Context(), venueID.(int), orderID)
	if err != nil {
		if errors.Is(err, domainOrder.ErrInvalidStatusTransition) {
			c.JSON(http.StatusUnprocessableEntity, httputil.ErrorResponse("La orden no esta lista para pago", "INVALID_TRANSITION"))
			return
		}
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error en checkout", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessMessageResponse("Pago procesado"))
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID de orden invalido", "BAD_REQUEST"))
		return
	}

	var req struct {
		StatusID int `json:"status_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	err = h.uc.UpdateOrderStatus(c.Request.Context(), venueID.(int), orderID, req.StatusID)
	if err != nil {
		if errors.Is(err, domainOrder.ErrInvalidStatusTransition) {
			c.JSON(http.StatusUnprocessableEntity, httputil.ErrorResponse("Transicion de estado invalida", "INVALID_TRANSITION"))
			return
		}
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al actualizar estado", "INTERNAL_ERROR"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessMessageResponse("Estado de orden actualizado"))
}

func (h *Handler) ListByTable(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)
	tableIDStr := c.Query("table_id")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID de mesa invalido", "BAD_REQUEST"))
		return
	}

	orders, err := h.uc.ListOrdersByTable(c.Request.Context(), venueID.(int), tableID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httputil.ErrorResponse("Error al listar ordenes", "INTERNAL_ERROR"))
		return
	}

	result := make([]gin.H, len(orders))
	for i, o := range orders {
		result[i] = toOrdenResponse(&o)
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(result))
}

// --- Division de cuenta ---

type DivideOrderRequest struct {
	TipoDivision string           `json:"tipo_division" binding:"required"`
	NumeroPartes int              `json:"numero_partes"`
	Divisiones   []DivisionDetail `json:"divisiones"`
}

type DivisionDetail struct {
	Items []string `json:"items"`
	Monto float64  `json:"monto"`
}

// DivideOrder maneja POST /ordenes/:id/dividir
func (h *Handler) DivideOrder(c *gin.Context) {
	venueID, _ := c.Get(middleware.VenueIDKey)
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("ID de orden invalido", "BAD_REQUEST"))
		return
	}

	var req DivideOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("Datos invalidos", "BAD_REQUEST"))
		return
	}

	// Obtener la orden para calcular division
	order, err := h.uc.GetOrderByID(c.Request.Context(), venueID.(int), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, httputil.ErrorResponse("Orden no encontrada", "NOT_FOUND"))
		return
	}

	subtotal := order.TotalAmount
	impuestos := subtotal * 0.19
	total := subtotal + impuestos

	var divisions []gin.H

	switch req.TipoDivision {
	case "partes_iguales":
		parts := req.NumeroPartes
		if parts <= 0 {
			parts = 2
		}
		partSubtotal := subtotal / float64(parts)
		partImpuestos := impuestos / float64(parts)
		partTotal := total / float64(parts)

		divisions = make([]gin.H, parts)
		for i := 0; i < parts; i++ {
			divisions[i] = gin.H{
				"division_id": "div_" + strconv.Itoa(i+1),
				"subtotal":    partSubtotal,
				"impuestos":   partImpuestos,
				"total":       partTotal,
			}
		}

	case "por_monto":
		remaining := total
		divisions = make([]gin.H, 0, len(req.Divisiones))
		for i, d := range req.Divisiones {
			divTotal := d.Monto
			if divTotal > remaining {
				divTotal = remaining
			}
			divSubtotal := divTotal / 1.19
			divImpuestos := divTotal - divSubtotal

			divisions = append(divisions, gin.H{
				"division_id": "div_" + strconv.Itoa(i+1),
				"subtotal":    divSubtotal,
				"impuestos":   divImpuestos,
				"total":       divTotal,
			})
			remaining -= divTotal
		}

	case "por_item":
		divisions = make([]gin.H, 0, len(req.Divisiones))
		for i := range req.Divisiones {
			// En implementacion real se calcularia por items asignados
			divTotal := total / float64(len(req.Divisiones))
			divSubtotal := divTotal / 1.19
			divImpuestos := divTotal - divSubtotal

			divisions = append(divisions, gin.H{
				"division_id": "div_" + strconv.Itoa(i+1),
				"subtotal":    divSubtotal,
				"impuestos":   divImpuestos,
				"total":       divTotal,
			})
		}

	default:
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse("tipo_division debe ser partes_iguales, por_item o por_monto", "BAD_REQUEST"))
		return
	}

	c.JSON(http.StatusOK, httputil.SuccessResponse(divisions))
}
