package product

type CreateMenuItemRequest struct {
	Name       string  `json:"name" binding:"required"`
	SalesPrice float64 `json:"sales_price" binding:"required,gte=0"`
}

type UpdateMenuItemRequest struct {
	Name       string  `json:"name,omitempty"`
	SalesPrice float64 `json:"sales_price,omitempty"`
}
