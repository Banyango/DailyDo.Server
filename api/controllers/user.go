package controllers

import (
	"github.com/Banyango/gifoody_server/api/services"
	"github.com/Banyango/gifoody_server/api/util"
	"github.com/Banyango/gifoody_server/model"
	"github.com/Banyango/gifoody_server/repositories"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"net/http"
	"time"
)

type UserController struct {
	userRepository    repositories.IUserRepository
	forgotRepository  repositories.IForgotRepository
	passwordService   services.IPasswordService
	emailTokenService services.IEmailTokenService
	jwtService        services.IJWTService
}

func NewUserController(
	repository repositories.IUserRepository,
	forgotRepository repositories.IForgotRepository,
	emailTokenService services.IEmailTokenService,
	service services.IPasswordService,
	jwtService services.IJWTService) *UserController {
	return &UserController{
		userRepository:    repository,
		forgotRepository:  forgotRepository,
		emailTokenService: emailTokenService,
		passwordService:   service,
		jwtService:        jwtService,
	}
}

// PostRegister will register a user in the system if they don't already exist.
// @Summary Register a user.
// @Description Register a new user in the system and send the confirmation email.
// @Accept json
// @Produce json
// @Param request body controllers.CreateUserRequest true "request"
// @Success 200 {object} model.User
// @Failure 401 {string} string "Email was already used to sign up"
// @Failure 500 {string} string "Token error"
// @Router /register [post]
func (self *UserController) PostRegister(c echo.Context) error {
	request, err := NewCreateUserRequestFromContext(c)
	if err != nil {
		return util.LogError(err, http.StatusBadRequest, err.Error())
	}

	user := model.User{
		Username:  request.Username,
		Email:     request.Email,
		LastName:  request.LastName,
		FirstName: request.FirstName,
	}

	userFound := <-self.userRepository.GetUserByEmail(user.Email)
	if userFound.Err == nil {
		return util.LogError(err, http.StatusUnauthorized, "Email was already used to sign up...")
	}

	hashedPassword, err := self.passwordService.HashPassword(request.Password)
	if err != nil {
		return util.LogError(err, http.StatusInternalServerError, "Internal server error")
	}

	token, id, err := self.emailTokenService.GenerateConfirmToken()
	if err != nil {
		return util.LogError(err, http.StatusInternalServerError, "Internal server error")
	}

	user.ConfirmToken = token
	user.Password = hashedPassword

	result := self.userRepository.Save(user)
	if result.Err != nil {
		return util.LogError(result.Err, http.StatusInternalServerError, "Internal server error")
	}

	err = self.emailTokenService.SendConfirmEmail(user.Email, id)
	if err != nil {
		return util.LogError(err, http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(200, user)
}

// PostResetPassword will send a token to a users email that will allow them to log back in and reset their password.
// @Summary Reset a password.
// @Description Generate a password reset token and send it to user's email.
// @Accept json
// @Produce json
// @Param request body controllers.CreatePasswordResetRequest true "request"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /reset_password [post]
func (self *UserController) PostResetPassword(c echo.Context) error {

	request, err := NewCreatePasswordResetFromContext(c)
	if err != nil {
		return util.LogError(err, http.StatusBadRequest, err.Error())
	}

	userByEmail := <-self.userRepository.GetUserByEmail(request.Email)
	if userByEmail.Err == nil {

		user := userByEmail.Data.(model.User)
		self.forgotRepository.Delete(user.Id.Hex())

		token, id, err := self.emailTokenService.GenerateConfirmToken()
		if err != nil {
			return util.LogError(err, http.StatusInternalServerError, "Could not generate token.")
		}

		resetPasswordUser := model.ForgotUser{
			Id:      user.Id,
			Token:   token,
			Created: time.Now(),
		}
		result := self.userRepository.SaveForgotUser(resetPasswordUser)
		if result.Err != nil {
			return util.LogError(err, http.StatusInternalServerError, "Internal Server error." + result.Err.Error())
		}

		self.emailTokenService.SendForgotEmail(user.Email, id)
	}

	return echo.NewHTTPError(http.StatusOK, "Password reset sent to your email.")
}

// PostConfirmResetPassword will verify the id token sent to a user's email
// and reset the password to the one provided if verified.
// @Summary verify the token and reset password.
// @Description will verify the id token sent to a user's email and reset the password to the one provided if verified
// @Accept json
// @Produce json
// @Param request body controllers.UpdatePasswordRequest true "request"
// @Success 200 {string} string
// @Failure 400 {string} string "Bad parameters"
// @Failure 500 {string} string "Token error"
// @Router /confirm_reset_password [post]
func (self *UserController) PostConfirmResetPassword(c echo.Context) error {

	request, err := NewUpdatePasswordRequestFromContext(c)
	if err != nil {
		return util.LogError(err, http.StatusBadRequest, err.Error())
	}

	idHash := self.passwordService.HashId(request.Token)

	if tokenChan := <-self.userRepository.GetUserByConfirmToken(idHash); tokenChan.Err == nil {
		forgotUser := tokenChan.Data.(model.ForgotUser)
		result := self.userRepository.DeleteForgotUser(forgotUser.Id)
		if result.Err != nil {
			return util.LogError(result.Err, http.StatusInternalServerError, "Internal server error")
		}

		userChan := <-self.userRepository.GetUserById(forgotUser.Id.Hex())
		if userChan.Err != nil {
			return util.LogError(result.Err, http.StatusInternalServerError, "Internal server error")
		}

		user := userChan.Data.(model.User)

		newPassword, err := self.passwordService.HashPassword(request.Password)
		if err != nil {
			return util.LogError(err, http.StatusInternalServerError, "Internal server error")
		}

		user.Password = newPassword
		saveResult := self.userRepository.Update(user)
		if saveResult.Err != nil {
			return util.LogError(saveResult.Err, http.StatusInternalServerError, "Internal server error")
		}

		err = self.emailTokenService.SendPasswordUpdatedEmail(user.Email)
		if err != nil {
			return util.LogError(err, http.StatusInternalServerError, "Internal server error")
		}

		return c.JSON(200, nil)
	} else {
		return echo.NewHTTPError(http.StatusInternalServerError, "Server JWT Error")
	}
}

// PostLogin will log a user in and create the jwt token.
// @Summary Login.
// @Description will log a user in and create the jwt token.
// @Accept json
// @Produce json
// @Param request body controllers.LoginRequest true "request"
// @Success 200 {string} string
// @Failure 400 {string} string "Bad parameters"
// @Failure 401 {string} string "Incorrect username password"
// @Failure 500 {string} string "Token error"
// @Router /login [post]
func (self *UserController) PostLogin(c echo.Context) error {

	request, err := NewLoginRequestFromContext(c)
	if err != nil {
		return util.LogError(err, http.StatusBadRequest, "Enter email/password")
	}

	userDb := <-self.userRepository.GetUserByEmail(request.Email)
	if userDb.Err != nil {
		return util.LogError(err, http.StatusUnauthorized, "Bad email/password")
	}

	fetchedUser := userDb.Data.(model.User)
	if self.passwordService.ComparePassword(request.Password, fetchedUser.Password) {
		self.jwtService.CreateJWTToken(c, fetchedUser.Id, fetchedUser.Username, fetchedUser.Email, fetchedUser.FirstName, fetchedUser.LastName)
		return c.JSON(200, fetchedUser)
	} else {
		return echo.NewHTTPError(http.StatusUnauthorized, "Email or Password incorrect.")
	}
}

// PostConfirmAccount will verify the id token that was sent to a user's email.
// @Summary Confirm Account.
// @Description will verify the id token that was sent to a user's email.
// @Accept json
// @Produce json
// @Param confirm body controllers.UpdateConfirmAccountRequest true "request"
// @Success 200 {string} string
// @Failure 400 {string} string "Bad parameters"
// @Failure 500 {string} string "Token error"
// @Router /confirm_account [post]
func (self *UserController) PostConfirmAccount(c echo.Context) error {

	request, err := NewUpdateConfirmAccountRequestFromContext(c)
	if err != nil {
		return util.LogError(err, http.StatusBadRequest, "Enter username/password")
	}

	token := self.passwordService.HashId(request.Token)

	if userChan := <-self.userRepository.GetUserByEmail(request.Email); userChan.Err == nil {
		user := userChan.Data.(model.User)

		if user.ConfirmToken == token {
			user.Verified = true
			user.ConfirmToken = ""
			saveResult := self.userRepository.Update(user)
			return saveResult.Err
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server JWT Error")
		}
	} else {
		return echo.NewHTTPError(http.StatusInternalServerError, "Server JWT Error")
	}
}

// PostLogout will log a user out by expiring the token
// @Summary Logout.
// @Description will expire the token and logout.
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Failure 500 {string} string "Token error"
// @Router /logout [post]
func (self *UserController) PostLogout(c echo.Context) error {
	self.jwtService.ExpireTokenImmediately(c)
	return c.JSON(http.StatusOK, "")
}

// PostMe will return ok if a user has a proper JWT Token
// or it will expire the token and return unauthorized if user is not found.
// @Summary Me.
// @Description Check the refresh_token against a user, logout if user doesn't exist.
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Failure 401 {string} string "User wasn't found from token."
// @Router /v1/me [post]
func (self *UserController) PostMe(c echo.Context) error {
	user := c.Get("refresh_token").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	userDb := <-self.userRepository.GetUserById(id)
	if userDb.Err == nil {
		fetchedUser := userDb.Data.(model.User)
		return c.JSON(http.StatusOK, fetchedUser)
	}

	// user wasn't found expire token.
	self.jwtService.ExpireTokenImmediately(c)
	return c.JSON(http.StatusUnauthorized, nil)
}
