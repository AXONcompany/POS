package ingredient

import "github.com/AXONcompany/POS/internal/domain/ingredient"

func toIngredientResponse(ing *ingredient.Ingredient) IngredientResponse {
	return IngredientResponse{
		ID:            ing.ID,
		Name:          ing.Name,
		UnitOfMeasure: ing.UnitOfMeasure,
		Type:          ing.IngredientType,
		Stock:         ing.Stock,
		CreatedAt:     ing.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toIngredientArrayResponse(ingredients []ingredient.Ingredient) []IngredientResponse {
	res := make([]IngredientResponse, len(ingredients))
	for i, ing := range ingredients {
		res[i] = toIngredientResponse(&ing)
	}
	return res
}
