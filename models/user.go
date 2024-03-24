package models

type User struct {
	_id      string `json:"_id"`
	username string `json:"username"`
	email    string `json:"email"`
	password string `json:"password"`
	debts    []Debt `json:"debts"`
}
