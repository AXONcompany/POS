package product

type RecipeItem struct {
	ID               int64
	ProductID        int64
	IngredientID     int64
	IngredientName   string
	UnitOfMeasure    string
	QuantityRequired float64
}
