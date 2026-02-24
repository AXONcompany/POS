package product

type CreateMenuItemRequest struct {
	Name        string              `json:"name" binding:"required"`
	SalesPrice  float64             `json:"sales_price" binding:"required,gte=0"`
	Ingredients []IngredientRequest `json:"ingredients" binding:"required,dive"`
}

type IngredientRequest struct {
	IngredientID int64   `json:"ingredient_id" binding:"required,gt=0"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
}

type UpdateMenuItemRequest struct {
	Name       string  `json:"name,omitempty"`
	SalesPrice float64 `json:"sales_price,omitempty"`
}
