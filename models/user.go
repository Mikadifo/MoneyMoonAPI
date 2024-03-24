package models

type User struct {
	_id      string `json:"_id"`
	username string `json:"username"`
	email    string `json:"email"`
	password string `json:"password"`
	banks    []Bank `json:"banks"`
	debts    []Debt `json:"debts"`
}
