package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty" validate:"required"`
	Email    string             `bson:"email,omitempty" validate:"required"`
	Password string             `bson:"password,omitempty" validate:"required"`
	Debts    []Debt             `bson:"debts"`
}
