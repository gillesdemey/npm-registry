package auth

type AuthProvider interface {
  Login(user, pass string) (token string, err error)
}

type HtpasswdProvider struct {
  File *HtpasswdFile
}

func NewHtpasswdProvider(file string) *HtpasswdProvider {
  return &HtpasswdProvider{
    File: NewHtpasswdFile(file),
  }
}
