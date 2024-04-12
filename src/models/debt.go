package models

type Debt struct {
	Name        string   `json:"name,omitempty" validate:"required"`
	Description string   `json:"description,omitempty" validate:"required"`
	Amount      float64  `json:"amount,omitempty" validate:"required"`
	Payed       *float64 `json:"payed,omitempty" validate:"required"`
}
