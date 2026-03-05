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
