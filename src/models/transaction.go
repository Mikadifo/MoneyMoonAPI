package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Date        string             `json:"date,omitempty" validate:"required"`
	DateObject  primitive.DateTime `bson:"dateObject,omitempty"`
	BankId      string             `bson:"bankId,omitempty" validate:"required" json:"bankId,omitempty"`
	Type        string             `json:"type,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
	Amount      float64            `json:"amount,omitempty" validate:"required"`
	Balance     *float64           `json:"balance" validate:"required"`
}
