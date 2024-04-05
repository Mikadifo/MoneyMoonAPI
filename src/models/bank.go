package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bank struct {
	Id     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name   string             `json:"name,omitempty" validate:"required"`
	UserId string             `json:"userId,omitempty" validate:"required"`
}
