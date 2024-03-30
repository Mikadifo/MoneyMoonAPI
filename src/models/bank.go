package models

type Bank struct {
	Id     string `bson:"_id,omitempty"`
	Name   string `bson:"name,omitempty" validate:"required"`
	UserId string `bson:"userId,omitempty" validate:"required"`
}
