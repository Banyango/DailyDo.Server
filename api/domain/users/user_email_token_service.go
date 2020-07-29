package users

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	. "github.com/Banyango/gifoody_server/api/infrastructure/mail"
	. "github.com/Banyango/gifoody_server/api/infrastructure/template"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"math/big"
	"net/http"
	"os"
)

type IEmailTokenService interface {
	GenerateConfirmToken() (token string, id string, err error)
	GenerateRandomASCIIString(length int) (string, error)
	SendConfirmEmail(email string, id string) (err error)
	SendForgotEmail(email string, id string) (err error)
	SendPasswordUpdatedEmail(email string) (err error)
}

type EmailTokenService struct {
	mailService     MailInterface
	templateService TemplateInterface
}

func NewEmailTokenService(mailService MailInterface, templateService TemplateInterface) *EmailTokenService {
	m := new(EmailTokenService)
	m.mailService = mailService
	m.templateService = templateService
	return m
}

func (self *EmailTokenService) GenerateConfirmToken() (token string, id string, err error) {
	id, err = self.GenerateRandomASCIIString(32)
	if err != nil {
		log.Error(err)
		return "", "", echo.NewHTTPError(http.StatusInternalServerError, "Internal Server error")
	}
	sha := sha256.New()
	sha.Write([]byte(id))
	return hex.EncodeToString(sha.Sum(nil)), id, nil
}

func (self *EmailTokenService) GenerateRandomASCIIString(length int) (string, error) {
	result := ""
	for {
		if len(result) >= length {
			return result, nil
		}
		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}
		n := num.Int64()
		// Make sure that the number/byte/letter is inside
		// the range of printable ASCII characters (excluding space and DEL)
		if (n > 47 && n < 58) || (n > 64 && n < 90) || (n > 97 && n < 123) {
			result += string(n)
		}
	}
}

func (self *EmailTokenService) SendConfirmEmail(email string, id string) (err error) {
	url := os.Getenv("SERVER_URL")

	data := struct {
		Link string
	}{
		Link: fmt.Sprintf("%s/confirm_account?token=%s", url, id),
	}

	var buffer bytes.Buffer
	err = self.templateService.RenderTemplate(&buffer, "confirm_user_account_email.html", data)
	if err != nil {
		log.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Confirmation email template error.")
	}

	mailFromAddress := os.Getenv("MAIL_FROM_ADDRESS")
	err = self.mailService.SendMail(email, mailFromAddress, "Reset your Password", &buffer)
	if err != nil {
		log.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Error Sending Confirmation Email")
	}

	return nil
}

func (self *EmailTokenService) SendForgotEmail(email string, id string) (err error) {
	url := os.Getenv("SERVER_URL")

	data := struct {
		Link string
	}{
		Link: fmt.Sprintf("%s/resetPassword?token=%s", url, id),
	}

	var buffer bytes.Buffer
	err = self.templateService.RenderTemplate(&buffer, "forgot_password_email.html", data)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	mailFromAddress := os.Getenv("MAIL_FROM_ADDRESS")
	err = self.mailService.SendMail(email, mailFromAddress, "Reset your Password", &buffer)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return nil
}

func (self *EmailTokenService) SendPasswordUpdatedEmail(email string) (err error) {
	var buffer bytes.Buffer

	err = self.templateService.RenderTemplate(&buffer, "confirm_new_password_email.html", nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	mailFromAddress := os.Getenv("MAIL_FROM_ADDRESS")
	err = self.mailService.SendMail(email, mailFromAddress, "Your password was reset.", &buffer)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return nil
}