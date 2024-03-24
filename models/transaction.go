// Package models provides ...
package models

type Transaction struct {
	_id         string  `json:"_id"`
	bankId      string  `json:"bankId"`
	transType   string  `json:"type"`
	description string  `json:"description"`
	amount      float32 `json:"amount"`
	date        string  `json:"date"`
}
