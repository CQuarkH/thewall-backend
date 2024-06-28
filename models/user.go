package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Photo    string `json:"photoUrl"`
	Token    string `json:"token,omitempty"`
}

type UserProfile struct {
	Username string `json:"username"`
	Photo    string `json:"photoUrl"`
}

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Photo    string `json:"photoUrl"`
}
