package product

import (
	"errors"
)

var (
	ErrNameEmpty        = errors.New("name cannot be empty")
	ErrPriceNegative    = errors.New("price cannot be negative")
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrInvalidID        = errors.New("id cannot be zero or negative")
	ErrCategoryNotFound = errors.New("category not found")
	ErrProductNotFound  = errors.New("product not found")
)
