package ingredient

type CreateIngredientRequest struct {
	Name          string `json:"name" binding:"required"`
	UnitOfMeasure string `json:"unit_of_measure" binding:"required"`
	Type          string `json:"type" binding:"required"`
	Stock         int64  `json:"stock"`
}

type UpdateIngredientRequest struct {
	Name          string `json:"name"`
	UnitOfMeasure string `json:"unit_of_measure"`
	Type          string `json:"type"`
	Stock         int64  `json:"stock"`
}

type IngredientResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	UnitOfMeasure string `json:"unit_of_measure"`
	Type          string `json:"type"`
	Stock         int64  `json:"stock"`
	CreatedAt     string `json:"created_at"`
}
