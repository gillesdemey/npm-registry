package model

type Login struct {
	Username string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
