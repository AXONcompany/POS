package ingredient

import (
	"time"
)

type Ingredient struct {
	ID             int64
	Name           string
	UnitOfMeasure  string
	IngredientType string
	Stock          int64

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// IngredientUpdates representa actualizaciones parciales
