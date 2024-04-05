package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Date        string             `json:"date,omitempty" validate:"required"`
	BankId      string             `json:"bankId,omitempty" validate:"required"`
	Type        string             `json:"type,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
	Amount      float64            `json:"amount,omitempty" validate:"required"`
	Balance     float64            `json:"balance,omitempty" validate:"required"`
}
