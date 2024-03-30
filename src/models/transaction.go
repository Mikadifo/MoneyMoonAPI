package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	BankId      string             `bson:"bankId,omitempty" validate:"required"`
	Type        string             `bson:"type,omitempty" validate:"required"`
	Description string             `bson:"description,omitempty" validate:"required"`
	Amount      float64            `bson:"amount,omitempty" validate:"required"`
	Date        string             `bson:"date,omitempty" validate:"required"`
}
