package routes

import (
  "github.com/gillesdemey/npm-registry/auth"
  "github.com/gillesdemey/npm-registry/model"
  "github.com/gillesdemey/npm-registry/storage"
  "gopkg.in/gin-gonic/gin.v1"
  "net/http"
  "log"
  "regexp"
)

// Create or verify a user named <username>
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

// Return the username associated with the NPM token
func Whoami(c *gin.Context) {
  re := regexp.MustCompile("(?i)Bearer ")
  authHeader := c.Request.Header.Get("Authorization")
  token := re.ReplaceAllString(authHeader, "")

  log.Printf("Whoami request for token '%s'", token)

  storage := c.Value("storage").(storage.StorageEngine)
  username, err := storage.RetrieveUsernameFromToken(token)

  if err != nil {
    c.Status(http.StatusUnauthorized)
    return
  }

  c.JSON(http.StatusOK, gin.H{"username": username})
}
