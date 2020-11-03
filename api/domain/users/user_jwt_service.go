package users

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
)

type IJWTService interface {
	CreateJWTToken(c echo.Context, objectId string, username string, email string, firstName string, lastName string) error
	ExpireTokenImmediately(c echo.Context)
}

type JWTService struct {
}

func NewJWTService() *JWTService {
	return new(JWTService)
}

func (self *JWTService) CreateJWTToken(c echo.Context, objectId string, username string, email string, firstName string, lastName string) error {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = objectId
	claims["name"] = username
	claims["email"] = email
	claims["firstName"] = firstName
	claims["lastName"] = lastName
	claims["exp"] = time.Now().Add(time.Hour * 12)

	// Generate encoded token and send it as response.
	jwtSecret, exists := os.LookupEnv("API_JWT_SECRET")
	if !exists {
		return echo.NewHTTPError(http.StatusInternalServerError, "Server JWT Error")
	}
	if t, err := token.SignedString([]byte(jwtSecret)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Server Signing JWT Error")
	} else {
		cookie := new(http.Cookie)

		cookie.HttpOnly = true
		// todo set this for prod.
		//cookie.Secure = true
		cookie.Name = "refresh_token"
		cookie.Value = t
		cookie.Expires = time.Now().Add(time.Hour * 12)
		c.SetCookie(cookie)
	}

	return nil
}

func (self *JWTService) ExpireTokenImmediately(c echo.Context) {
	cookie, _ := c.Cookie("refresh_token")
	if cookie != nil {
		cookie.Expires = time.Now().Add(-100 * time.Hour)
		cookie.MaxAge = -1
		c.SetCookie(cookie)
	}
}
