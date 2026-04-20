package order_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domainOrder "github.com/AXONcompany/POS/internal/domain/order"
	orderrest "github.com/AXONcompany/POS/internal/infrastructure/rest/order"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// stubUsecase implementa OrderUsecase con valores configurables por test.
type stubUsecase struct {
	createOrderResult  *domainOrder.Order
	createOrderErr     error
	addProductErr      error
	getOrderResult     *domainOrder.Order
	getOrderErr        error
	cancelItemErr      error
	divideOrderResult  []domainOrder.OrderDivision
	divideOrderErr     error
	getDivisionsResult []domainOrder.OrderDivision
	getDivisionsErr    error
	checkoutErr        error
}

func (s *stubUsecase) CreateOrderWithoutItems(_ context.Context, _, _ int, _ *int64) (*domainOrder.Order, error) {
	return s.createOrderResult, s.createOrderErr
}
func (s *stubUsecase) GetOrderByID(_ context.Context, _ int, _ int64) (*domainOrder.Order, error) {
	return s.getOrderResult, s.getOrderErr
}
func (s *stubUsecase) AddProductToOrder(_ context.Context, _ int, _ int64, _ []domainOrder.OrderItem) error {
	return s.addProductErr
}
func (s *stubUsecase) CancelOrderItem(_ context.Context, _, _ int, _, _ int64) error {
	return s.cancelItemErr
}
func (s *stubUsecase) UpdateOrderStatus(_ context.Context, _ int, _ int64, _ int) error {
	return nil
}
func (s *stubUsecase) CheckoutOrder(_ context.Context, _ int, _ int64) error {
	return s.checkoutErr
}
func (s *stubUsecase) ListOrdersByTable(_ context.Context, _ int, _ int64) ([]domainOrder.Order, error) {
	return nil, nil
}
func (s *stubUsecase) DivideOrder(_ context.Context, _ int, _ int64, _ string, _ int, _ []float64) ([]domainOrder.OrderDivision, error) {
	return s.divideOrderResult, s.divideOrderErr
}
func (s *stubUsecase) GetDivisionsByOrder(_ context.Context, _ int, _ int64) ([]domainOrder.OrderDivision, error) {
	return s.getDivisionsResult, s.getDivisionsErr
}

// setupRouter crea un engine de gin en modo test con el handler inyectado
// y los context keys de middleware pre-seteados.
func setupRouter(stub orderrest.OrderUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := orderrest.NewHandler(stub)

	injectKeys := func(c *gin.Context) {
		c.Set("venue_id", 1)
		c.Set("user_id", 1)
		c.Next()
	}

	api := r.Group("", injectKeys)
	api.POST("/ordenes", h.CreateOrder)
	api.POST("/ordenes/:id/items", h.AddItems)
	api.DELETE("/ordenes/:id/items/:item_id", h.CancelItem)
	api.POST("/ordenes/:id/dividir", h.DivideOrder)
	api.GET("/ordenes/:id/divisiones", h.GetDivisions)
	api.POST("/ordenes/:id/checkout", h.CheckoutOrder)

	return r
}

func jsonBody(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

// --- CreateOrder ---

func TestHandler_CreateOrder_201(t *testing.T) {
	stub := &stubUsecase{
		createOrderResult: &domainOrder.Order{ID: 1, StatusID: 1},
		getOrderResult:    &domainOrder.Order{ID: 1, StatusID: 1},
	}
	r := setupRouter(stub)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ordenes", jsonBody(map[string]any{"mesa_id": "1"}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var body map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Contains(t, body, "data")
}

func TestHandler_CreateOrder_500(t *testing.T) {
	stub := &stubUsecase{createOrderErr: errors.New("db error")}
	r := setupRouter(stub)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ordenes", jsonBody(map[string]any{}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- AddItems ---

func TestHandler_AddItems_409_InsufficientStock(t *testing.T) {
	stub := &stubUsecase{addProductErr: domainOrder.ErrInsufficientStock}
	r := setupRouter(stub)

	body := jsonBody(map[string]any{
		"items": []map[string]any{{"menu_item_id": "1", "cantidad": 1}},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ordenes/1/items", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "INSUFFICIENT_STOCK")
}

func TestHandler_AddItems_422_InvalidTransition(t *testing.T) {
	stub := &stubUsecase{addProductErr: domainOrder.ErrInvalidStatusTransition}
	r := setupRouter(stub)

	body := jsonBody(map[string]any{
		"items": []map[string]any{{"menu_item_id": "1", "cantidad": 1}},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ordenes/1/items", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), "INVALID_TRANSITION")
}

// --- CancelItem ---

func TestHandler_CancelItem_409_AlreadyCancelled(t *testing.T) {
	stub := &stubUsecase{cancelItemErr: domainOrder.ErrItemAlreadyCancelled}
	r := setupRouter(stub)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/ordenes/1/items/10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "ITEM_ALREADY_CANCELLED")
}

func TestHandler_CancelItem_422_InvalidTransition(t *testing.T) {
	stub := &stubUsecase{cancelItemErr: domainOrder.ErrInvalidStatusTransition}
	r := setupRouter(stub)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/ordenes/1/items/10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// --- DivideOrder ---

func TestHandler_DivideOrder_200(t *testing.T) {
	stub := &stubUsecase{
		divideOrderResult: []domainOrder.OrderDivision{
			{ID: "div_1_1", Amount: 50, Tax: 9.5, Total: 59.5},
			{ID: "div_1_2", Amount: 50, Tax: 9.5, Total: 59.5},
		},
	}
	r := setupRouter(stub)

	body := jsonBody(map[string]any{"tipo_division": "partes_iguales", "numero_partes": 2})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ordenes/1/dividir", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].([]any)
	assert.Len(t, data, 2)
}

func TestHandler_DivideOrder_409_AlreadyPaid(t *testing.T) {
	stub := &stubUsecase{divideOrderErr: domainOrder.ErrDivisionAlreadyPaid}
	r := setupRouter(stub)

	body := jsonBody(map[string]any{"tipo_division": "partes_iguales", "numero_partes": 2})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ordenes/1/dividir", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "DIVISION_ALREADY_PAID")
}

// --- GetDivisions ---

func TestHandler_GetDivisions_200(t *testing.T) {
	stub := &stubUsecase{
		getDivisionsResult: []domainOrder.OrderDivision{
			{ID: "div_1_1", Amount: 50, Tax: 9.5, Total: 59.5, IsPaid: false},
		},
	}
	r := setupRouter(stub)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ordenes/1/divisiones", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	data := resp["data"].([]any)
	assert.Len(t, data, 1)
}

// --- CheckoutOrder ---

func TestHandler_CheckoutOrder_422_InvalidTransition(t *testing.T) {
	stub := &stubUsecase{checkoutErr: domainOrder.ErrInvalidStatusTransition}
	r := setupRouter(stub)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ordenes/1/checkout", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
