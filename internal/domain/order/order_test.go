package order_test

import (
	"testing"

	"github.com/AXONcompany/POS/internal/domain/order"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	tableID := int64(1)

	t.Run("success", func(t *testing.T) {
		items := []order.OrderItem{
			{ProductID: 1, Quantity: 2, UnitPrice: 15.0},
			{ProductID: 2, Quantity: 1, UnitPrice: 10.0},
		}

		o, err := order.NewOrder(1, 1, &tableID, items)

		assert.NoError(t, err)
		assert.NotNil(t, o)
		assert.Equal(t, float64(40.0), o.TotalAmount)
		assert.Equal(t, 1, o.StatusID)
	})

	t.Run("empty items", func(t *testing.T) {
		o, err := order.NewOrder(1, 1, &tableID, []order.OrderItem{})

		assert.ErrorIs(t, err, order.ErrInvalidOrderItems)
		assert.Nil(t, o)
	})
}

func TestCanTransitionTo(t *testing.T) {
	cases := []struct {
		current int
		next    int
		valid   bool
		label   string
	}{
		// Flujo normal paso a paso
		{1, 2, true, "PENDING → SENT"},
		{2, 3, true, "SENT → PREPARING"},
		{3, 4, true, "PREPARING → READY"},
		{4, 5, true, "READY → PAID"},
		// Saltos hacia adelante (MVP sin KDS)
		{1, 3, true, "PENDING → PREPARING (salto)"},
		{1, 4, true, "PENDING → READY (salto)"},
		{1, 5, true, "PENDING → PAID (express)"},
		{2, 5, true, "SENT → PAID (salto)"},
		{3, 5, true, "PREPARING → PAID (salto)"},
		// Cancelación desde cualquier estado no terminal
		{1, 6, true, "PENDING → CANCELLED"},
		{2, 6, true, "SENT → CANCELLED"},
		{3, 6, true, "PREPARING → CANCELLED"},
		{4, 6, true, "READY → CANCELLED"},
		// Retrocesos y terminales — inválidos
		{5, 1, false, "PAID → PENDING (reverso)"},
		{6, 1, false, "CANCELLED → PENDING (reverso)"},
		{4, 3, false, "READY → PREPARING (reverso)"},
		{5, 6, false, "PAID → CANCELLED (terminal)"},
		{6, 5, false, "CANCELLED → PAID (terminal)"},
		{0, 1, false, "estado desconocido"},
	}

	for _, tc := range cases {
		got := order.CanTransitionTo(tc.current, tc.next)
		if got != tc.valid {
			t.Errorf("%s: CanTransitionTo(%d, %d) = %v, want %v", tc.label, tc.current, tc.next, got, tc.valid)
		}
	}
}

func TestNewOrder_EdgeCases(t *testing.T) {
	tableID := int64(42)

	t.Run("nil tableID (delivery order)", func(t *testing.T) {
		items := []order.OrderItem{{ProductID: 1, Quantity: 1, UnitPrice: 20.0}}
		o, err := order.NewOrder(1, 1, nil, items)

		assert.NoError(t, err)
		assert.NotNil(t, o)
		assert.Nil(t, o.TableID)
		assert.Equal(t, float64(20.0), o.TotalAmount)
	})

	t.Run("single item total calculation", func(t *testing.T) {
		items := []order.OrderItem{{ProductID: 1, Quantity: 3, UnitPrice: 7.5}}
		o, err := order.NewOrder(1, 1, &tableID, items)

		assert.NoError(t, err)
		assert.InDelta(t, 22.5, o.TotalAmount, 0.001)
	})

	t.Run("free item with zero price", func(t *testing.T) {
		items := []order.OrderItem{{ProductID: 1, Quantity: 5, UnitPrice: 0.0}}
		o, err := order.NewOrder(1, 1, &tableID, items)

		assert.NoError(t, err)
		assert.Equal(t, float64(0), o.TotalAmount)
	})

	t.Run("multiple items total is sum of all", func(t *testing.T) {
		items := []order.OrderItem{
			{ProductID: 1, Quantity: 1, UnitPrice: 10.0},
			{ProductID: 2, Quantity: 2, UnitPrice: 5.0},
			{ProductID: 3, Quantity: 3, UnitPrice: 3.0},
		}
		o, err := order.NewOrder(1, 1, &tableID, items)

		assert.NoError(t, err)
		// 1*10 + 2*5 + 3*3 = 10 + 10 + 9 = 29
		assert.InDelta(t, 29.0, o.TotalAmount, 0.001)
	})

	t.Run("status always starts as pending (1)", func(t *testing.T) {
		items := []order.OrderItem{{ProductID: 1, Quantity: 1, UnitPrice: 5.0}}
		o, err := order.NewOrder(99, 99, &tableID, items)

		assert.NoError(t, err)
		assert.Equal(t, 1, o.StatusID)
		assert.Equal(t, 99, o.VenueID)
		assert.Equal(t, 99, o.UserID)
	})

	t.Run("nil items slice treated as empty", func(t *testing.T) {
		o, err := order.NewOrder(1, 1, &tableID, nil)

		assert.ErrorIs(t, err, order.ErrInvalidOrderItems)
		assert.Nil(t, o)
	})
}
