package sale

import (
	domainSale "github.com/AXONcompany/POS/internal/domain/sale"
)

type ProcessPaymentRequest struct {
	OrderID       int64     `json:"order_id" binding:"required"`
	RestaurantID  int     `json:"restaurant_id" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Total         float64 `json:"total" binding:"required"`
}

type SplitOrderRequest struct {
	Total   float64 `json:"total" binding:"required"`
	People  int     `json:"personas" binding:"required"`
}

type SaleResponse struct {
	ID            int64   `json:"id"`
	Total         float64 `json:"total"`
	PaymentMethod string  `json:"payment_method"`
	Date          string  `json:"date"`
	OrderID       int64   `json:"order_id"`
	CreatedAt     string  `json:"created_at"`
}

type SplitResponse struct {
	Total           float64 `json:"total"`
	People          int     `json:"personas"`
	AmountPerPerson float64 `json:"monto_por_persona"`
}

func toSaleResponse(s *domainSale.Sale) SaleResponse {
	return SaleResponse{
		ID:            s.ID,
		Total:         s.Total,
		PaymentMethod: s.PaymentMethod,
		Date:          s.Date.Format("2006-01-02 15:04:05"),
		OrderID:       s.OrderID,
		CreatedAt:     s.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}