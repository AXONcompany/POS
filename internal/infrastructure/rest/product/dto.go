package product

import "time"

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type CreateProductRequest struct {
	Name       string  `json:"name" binding:"required"`
	SalesPrice float64 `json:"sales_price" binding:"required,gte=0"`
	IsActive   bool    `json:"is_active"`
}

type ProductResponse struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	SalesPrice float64    `json:"sales_price"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}
