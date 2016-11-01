package model

type User struct {
	Username, Password, Email, Token string
}

type Login struct {
	Username string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
