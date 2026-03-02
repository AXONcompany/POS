package sale

import (
	"time"
)

type Sale struct {
	ID 		      int64 `json:"id" db:"id"`
	Total 	 	  float64 `json:"total" db:"total"`
	PaymentMethod string `json:"payment_method" db:"payment_method"`
	Date 		  time.Time `json:"date" db:"date"`
	OrderID 	  int64 `json:"order_id" db:"order_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time 
}
