package storage

type StorageEngine interface {
  initialize() error
  StoreTarball() error
  RetrieveTarball() ([]byte, error)
  RetrieveUser() (User, error)
  StoreUserToken(username string, token string) error
}

type User struct {
  Username, Password, Email, Token string
}
