package model

import (
	"time"
)

type Token struct {
	Username  string
	Timestamp time.Time
}
