package model

import (
	"time"
)

type ForgotUser struct {
	Id      string    `json:"id" db:"id"`
	Token   string    `json:"token" db:"token"`
	Created time.Time `db:"created"`
}
