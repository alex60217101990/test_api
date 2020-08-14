package models

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
