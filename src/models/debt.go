package models

type Debt struct {
	Name        string  `bson:"name,omitempty" validate:"required"`
	Description string  `bson:"description,omitempty" validate:"required"`
	Amount      float64 `bson:"amount,omitempty" validate:"required"`
	Payed       float64 `bson:"payed,omitempty" validate:"required"`
}
