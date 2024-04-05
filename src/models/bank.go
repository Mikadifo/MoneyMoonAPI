package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bank struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"name,omitempty" validate:"required"`
	UserId string             `bson:"userId,omitempty" validate:"required"`
}
