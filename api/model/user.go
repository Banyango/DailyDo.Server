package model

import (
	"fmt"
)

type User struct {
	Id           string `json:"id", db:"id"`
	FirstName    string `json:"firstName", db:"first_name"`
	LastName     string `json:"lastName, db:"last_name"`
	Email        string `json:"email, db:"email"`
	Username     string `json:"userName, db:"username"`
	Password     string `json:"-, db:"password"`
	ConfirmToken string `json:"-, db:"confirm_token"`
	Verified     bool   `json:"-, db:"verified"`
	Reset        bool   `json:"-, db:"reset"`
}

func New() User {
	u := User{}
	return u
}

func (u *User) ValidateForLogin() error {

	if u.Email == "" {
		return fmt.Errorf("Email must be not nil")
	}

	if u.Password == "" {
		return fmt.Errorf("Password cannot be nil")
	}

	return nil
}
