package users

import (
	"github.com/Banyango/gifoody_server/api/infrastructure/mail"
	"github.com/Banyango/gifoody_server/api/infrastructure/template"
	"github.com/Banyango/gifoody_server/api/infrastructure/utils"
	. "github.com/Banyango/gifoody_server/api/model"
	"github.com/Banyango/gifoody_server/api/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type UserController struct {
	userRepository    repositories.IUserRepository
	emailTokenService IEmailTokenService
	passwordService   IPasswordService
	jwtService        IJWTService
	userAuthService   IUserAuthService
}

func NewUserController(
	repository repositories.IUserRepository,
	mailService mail.MailInterface,
	templateService template.TemplateInterface,
	userAuthService IUserAuthService) *UserController {

	emailTokenService := NewEmailTokenService(mailService, templateService)
	passwordService := NewPasswordService()
	jwtService := NewJWTService()

	return &UserController{
		userRepository:    repository,
		emailTokenService: emailTokenService,
		passwordService:   passwordService,
		jwtService:        jwtService,
		userAuthService:   userAuthService,
	}
}

// PostRegister will register a user in the system if they don't already exist.
// @Summary Register a user.
// @Description Register a new user in the system and send the confirmation email.
// @Accept json
// @Produce json
// @Param request body controllers.CreateUserRequest true "request"
// @Success 200 {object} User
// @Failure 401 {string} string "Email was already used to sign up"
// @Failure 500 {string} string "Token error"
// @Router /register [post]
func (self *UserController) PostRegister(c echo.Context) error {
	request, err := NewCreateUserRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, err.Error())
	}

	user := User{
		Username:  request.Username,
		Email:     request.Email,
		LastName:  request.LastName,
		FirstName: request.FirstName,
	}

	userFound := <-self.userRepository.GetUserByEmailOrUsernameAsync(user.Email, user.Username)
	if userFound.Err == nil {
		return utils.LogError(err, http.StatusBadRequest, "Username/Email was already used to sign up...")
	}

	hashedPassword, err := self.passwordService.HashPassword(request.Password)
	if err != nil {
		return utils.LogError(err, http.StatusInternalServerError, "Internal server error")
	}

	token, id, err := self.emailTokenService.GenerateConfirmToken()
	if err != nil {
		return utils.LogError(err, http.StatusInternalServerError, "Internal server error")
	}

	user.ConfirmToken = token
	user.Password = hashedPassword

	user.Id = uuid.New().String()

	result := self.userRepository.Save(user)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Internal server error")
	}

	err = self.emailTokenService.SendConfirmEmail(user.Email, id)
	if err != nil {
		return utils.LogError(err, http.StatusInternalServerError, "Internal server error")
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
		return utils.LogError(err, http.StatusBadRequest, err.Error())
	}

	userByEmail := <-self.userRepository.GetUserByEmailOrUsernameAsync(request.Email, "")
	if userByEmail.Err == nil {

		user := userByEmail.Data.(User)
		self.userRepository.DeleteForgotUser(user.Id)

		token, id, err := self.emailTokenService.GenerateConfirmToken()
		if err != nil {
			return utils.LogError(err, http.StatusInternalServerError, "Could not generate token.")
		}

		resetPasswordUser := ForgotUser{
			Id:      user.Id,
			Token:   token,
			Created: time.Now(),
		}
		result := self.userRepository.SaveForgotUser(resetPasswordUser)
		if result.Err != nil {
			return utils.LogError(err, http.StatusInternalServerError, "Internal Server error."+result.Err.Error())
		}

		self.emailTokenService.SendForgotEmail(user.Email, id)

		return echo.NewHTTPError(http.StatusCreated, "Password reset sent to your email.")
	} else {
		return echo.NewHTTPError(http.StatusConflict, "Email exists.")
	}

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
		return utils.LogError(err, http.StatusBadRequest, err.Error())
	}

	idHash := self.passwordService.HashId(request.Token)

	if tokenChan := <-self.userRepository.GetForgotUserByTokenAsync(idHash); tokenChan.Err == nil {
		forgotUser := tokenChan.Data.(ForgotUser)
		result := self.userRepository.DeleteForgotUser(forgotUser.Id)
		if result.Err != nil {
			return utils.LogError(result.Err, http.StatusInternalServerError, "Internal server error")
		}

		userChan := <-self.userRepository.GetUserByIdAsync(forgotUser.Id)
		if userChan.Err != nil {
			return utils.LogError(result.Err, http.StatusInternalServerError, "Internal server error")
		}

		user := userChan.Data.(User)

		newPassword, err := self.passwordService.HashPassword(request.Password)
		if err != nil {
			return utils.LogError(err, http.StatusInternalServerError, "Internal server error")
		}

		user.Password = newPassword
		saveResult := <-self.userRepository.UpdateAsync(user)
		if saveResult.Err != nil {
			return utils.LogError(saveResult.Err, http.StatusInternalServerError, "Internal server error")
		}

		err = self.emailTokenService.SendPasswordUpdatedEmail(user.Email)
		if err != nil {
			return utils.LogError(err, http.StatusInternalServerError, "Internal server error")
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
		return utils.LogError(err, http.StatusBadRequest, "Enter email/password")
	}

	userDb := <-self.userRepository.GetUserByEmailOrUsernameAsync(request.Email, "")
	if userDb.Err != nil {
		return utils.LogError(err, http.StatusUnauthorized, "Bad email/password")
	}

	fetchedUser := userDb.Data.(User)
	if self.passwordService.ComparePassword(request.Password, fetchedUser.Password) {
		err := self.jwtService.CreateJWTToken(c, fetchedUser.Id, fetchedUser.Username, fetchedUser.Email, fetchedUser.FirstName, fetchedUser.LastName)
		if err != nil {
			return utils.LogError(err, http.StatusInternalServerError, "JWT error")
		}

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
		return utils.LogError(err, http.StatusBadRequest, "Bad Parameters")
	}

	token := self.passwordService.HashId(request.Token)

	if userChan := <-self.userRepository.GetUserByConfirmTokenAsync(token); userChan.Err == nil {
		user := userChan.Data.(User)
		user.Verified = true
		user.ConfirmToken = ""
		saveResult := <-self.userRepository.UpdateAsync(user)
		return saveResult.Err
	} else {
		return echo.NewHTTPError(http.StatusInternalServerError, "Confirmation error")
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

// GetMe will return a user if a user has a proper JWT Token
// or it will expire the token and return unauthorized if user is not found.
// @Summary Me.
// @Description Check the refresh_token against a user, logout if user doesn't exist.
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 401 {string} string "User wasn't found from token."
// @Router /v1/me [get]
func (self *UserController) GetMe(c echo.Context) error {

	user, err := self.userAuthService.GetLoggedInUser(c)
	if err == nil {
		return c.JSON(http.StatusOK, user)
	}

	// user wasn't found expire token.
	self.jwtService.ExpireTokenImmediately(c)
	return c.JSON(http.StatusUnauthorized, nil)
}
