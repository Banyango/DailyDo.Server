package users

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type CreateUserRequest struct {
	Email     string
	Password  string
	Username  string
	FirstName string
	LastName  string
}

func NewCreateUserRequestFromContext(c echo.Context) (request *CreateUserRequest, err error) {
	request = new(CreateUserRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	message := "missing "
	if request.Username == "" {
		message = message + "username "
	}
	if request.Email == "" {
		message = message + "email "
	}
	if request.LastName == "" {
		message = message + "lastName "
	}
	if request.FirstName == "" {
		message = message + "firstName "
	}
	if request.FirstName == "" {
		return nil, fmt.Errorf(message)
	}

	return request, nil
}

type CreatePasswordResetRequest struct {
	Email string
}

func NewCreatePasswordResetFromContext(c echo.Context) (request *CreatePasswordResetRequest, err error) {
	request = new(CreatePasswordResetRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	if request.Email == "" {
		return nil, fmt.Errorf("Incorrect Parameters")
	}

	return request, nil
}

type UpdatePasswordRequest struct {
	Email string
	Token string
	Password string
}

func NewUpdatePasswordRequestFromContext(c echo.Context) (request *UpdatePasswordRequest, err error) {
	request = new(UpdatePasswordRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	if request.Email == "" || request.Token == "" || request.Password == "" {
		return nil, fmt.Errorf("Incorrect Parameters")
	}

	return request, nil
}

type LoginRequest struct {
	Email string
	Password string
}

func NewLoginRequestFromContext(c echo.Context) (request *LoginRequest, err error) {
	request = new(LoginRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	if request.Email == "" || request.Password == "" {
		return nil, fmt.Errorf("Incorrect Parameters")
	}

	return request, nil
}

type UpdateConfirmAccountRequest struct {
	Token string
}

func NewUpdateConfirmAccountRequestFromContext(c echo.Context) (request *UpdateConfirmAccountRequest, err error) {
	request = new(UpdateConfirmAccountRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	if request.Token == "" {
		return nil, fmt.Errorf("Incorrect Parameters")
	}

	return request, nil
}
