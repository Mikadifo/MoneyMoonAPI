package models

type Debt struct {
	id          string  `json:"id"`
	name        string  `json:"name"`
	description string  `json:"description"`
	amount      float32 `json:"amount"`
	payed       float32 `json:"payed"`
}
