package product

// ── Category DTOs ──

type CreateCategoryRequest struct {
	Name       string `json:"name" binding:"required"`
	ColorClass string `json:"colorClass"`
	Icon       string `json:"icon"`
}

type CategoryResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	ColorClass string `json:"colorClass,omitempty"`
	Icon       string `json:"icon,omitempty"`
}

// ── Product DTOs ──

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Description string  `json:"description"`
	CategoryId  *int64  `json:"categoryId"`
	Image       string  `json:"image"`
	IsAvailable bool    `json:"isAvailable"`
}

type ProductResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description,omitempty"`
	CategoryId  *int64  `json:"categoryId,omitempty"`
	Image       string  `json:"image,omitempty"`
	IsAvailable bool    `json:"isAvailable"`
}

