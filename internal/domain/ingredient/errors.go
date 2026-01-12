package ingredient

import "errors"

var (
	ErrNameEmpty     = errors.New("ingredient name cannot be empty")
	ErrNegativeStock = errors.New("ingredient stock value cannot be negative")
	ErrNotFound      = errors.New("ingredient not found")
	ErrAlreadyExists = errors.New("ingredient already exists")
)
