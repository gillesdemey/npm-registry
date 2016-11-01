package routes

import (
  "github.com/gillesdemey/npm-registry/auth"
  "github.com/gillesdemey/npm-registry/model"
  "github.com/gillesdemey/npm-registry/storage"
  "gopkg.in/gin-gonic/gin.v1"
  "net/http"
)

func Login(c *gin.Context) {
  var login model.Login
  var err error

  if err = c.BindJSON(&login); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
    return
  }

  storage := c.Value("storage").(storage.StorageEngine)

  username := login.Username
  password := login.Password

  token, err := auth.Login(username, password)
  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
    return
  }

  err = storage.StoreUserToken(token, username)
  if err != nil {
    c.Status(http.StatusInternalServerError)
    return
  }

  c.JSON(http.StatusCreated, gin.H{"token": token})
}
