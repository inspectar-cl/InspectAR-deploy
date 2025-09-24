package models

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Empresa  string `json:"empresa"`
	Password string `json:"password"`
	Scope    string `json:"scope"`
}
