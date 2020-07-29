package api

import (
	"github.com/Banyango/gifoody_server/api/domain/index"
	"github.com/Banyango/gifoody_server/api/domain/posts"
	"github.com/Banyango/gifoody_server/api/domain/users"
	"github.com/Banyango/gifoody_server/api/infrastructure/mail"
	"github.com/Banyango/gifoody_server/api/infrastructure/os"
	"github.com/Banyango/gifoody_server/api/infrastructure/pagination"
	"github.com/Banyango/gifoody_server/api/infrastructure/template"
	"github.com/Banyango/gifoody_server/constants"
	"github.com/Banyango/gifoody_server/api/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func InitRouter(echo *echo.Echo, db *sqlx.DB) {
	mainGroup := echo.Group(constants.API_PATH)

	// repositories
	store := repositories.NewAppStore(db)

	// services
	osService := os.NewOSService()
	mailService := mail.NewMailService()
	templateService := template.NewTemplateService(osService)

	//controllers
	postController := posts.NewPostController(store.Post())
	userController := users.NewUserController(store.User(), mailService, templateService)
	indexController := index.NewIndexController()

	// index
	mainGroup.GET("index", indexController.GetIndex)

	// posts
	mainGroup.GET("posts", postController.ListPosts, pagination.Paginate())

	// users
	mainGroup.POST("register", userController.PostRegister)
	mainGroup.POST("reset_password", userController.PostResetPassword)
	mainGroup.POST("confirm_reset_password", userController.PostConfirmResetPassword)
	mainGroup.POST("confirm_account", userController.PostConfirmAccount)
	mainGroup.POST("login", userController.PostLogin)
	mainGroup.POST("logout", userController.PostLogout)
}
