package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `json:"name,omitempty" validate:"required"`
	UserId       string             `bson:"userId,omitempty" validate:"required" json:"userId,omitempty"`
	Total        float64            `json:"total,omitempty" validate:"required"`
	Transactions []string           `json:"transactions"`
}
