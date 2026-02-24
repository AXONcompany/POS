package product

type AddIngredientRequest struct {
	IngredientID int64   `json:"ingredient_id" binding:"required,gt=0"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
}

type RecipeItemResponse struct {
	ID             int64   `json:"id"`
	ProductID      int64   `json:"product_id"`
	IngredientID   int64   `json:"ingredient_id"`
	IngredientName string  `json:"ingredient_name,omitempty"`
	UnitOfMeasure  string  `json:"unit_of_measure,omitempty"`
	Quantity       float64 `json:"quantity"`
}
