package users

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type IPasswordService interface {
	ValidatePassword(password string) error
	ComparePassword(password string, hashedPassword string) bool
	HashId(id string) (idHash string)
	HashPassword(password string) (passwordHash string, err error)
}

type PasswordService struct {
}

const PasswordCost = 10

func NewPasswordService() *PasswordService {
	m := new(PasswordService)
	return m
}

func (s *PasswordService) ValidatePassword(password string) error {
	if password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect parameters")
	}
	if len(password) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "Password must be longer than 6 characters")
	}
	return nil
}

func (self *PasswordService) ComparePassword(password string, hashedPassword string) bool {
	if len(password) == 0 || len(hashedPassword) == 0 {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (self *PasswordService) HashPassword(password string) (passwordHash string, err error) {
	var hash []byte

	hash, err = bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	return string(hash), err
}

func (self *PasswordService) HashId(id string) (idHash string) {
	sha := sha256.New()
	sha.Write([]byte(id))
	return hex.EncodeToString(sha.Sum(nil))
}
