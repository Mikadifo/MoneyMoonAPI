package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	Id          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	BankId      string             `json:"bankId,omitempty" validate:"required"`
	Type        string             `json:"type,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
	Amount      float64            `json:"amount,omitempty" validate:"required"`
	Date        string             `json:"date,omitempty" validate:"required"`
}
