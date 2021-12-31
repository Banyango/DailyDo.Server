package users

import (
	"fmt"
	"github.com/Banyango/dailydo_server/api/model"
	"github.com/Banyango/dailydo_server/api/repositories"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type IUserAuthService interface {
	GetLoggedInUser(c echo.Context) (model.User, error)
	GetUserIdFromContext(c echo.Context) (string, error)
}

type UserAuthService struct {
	UserRepository repositories.IUserRepository
}

func NewUserAuthService(userRepository repositories.IUserRepository) *UserAuthService {
	m := new(UserAuthService)
	m.UserRepository = userRepository
	return m
}

func (s *UserAuthService) GetLoggedInUser(c echo.Context) (model.User, error) {
	userInterface := c.Get("user")
	if userInterface == nil {
		return model.User{}, fmt.Errorf("token not found")
	}

	user := userInterface.(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	userById := <-s.UserRepository.GetUserByIdAsync(id)
	if userById.Err != nil {
		return model.User{}, userById.Err
	}
	return userById.Data.(model.User), nil
}

func (s *UserAuthService) GetUserIdFromContext(c echo.Context) (string, error) {
	userInterface := c.Get("user")
	if userInterface == nil {
		return "", fmt.Errorf("token not found")
	}
	user := userInterface.(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	return id, nil
}
