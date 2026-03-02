package sale

import "errors"

var (
	ErrSaleNotFound 	  = errors.New("Sale not found")
	ErrInvalidID 		  = errors.New("invalid id")
	ErrInvalidOrderID 	  = errors.New("invalid order id")
	ErrPaymentMethodEmpty = errors.New("order already paid")
	ErrInvalidTotal       = errors.New("total must be greater than 0")
)