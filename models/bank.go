package models

type Bank struct {
	_id          string        `json:"_id"`
	userId       string        `json:"userId"`
	name         string        `json:"name"`
	transactions []Transaction `json:"transactions"`
}
