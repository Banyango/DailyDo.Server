package api

import (
	"github.com/Banyango/gifoody_server/api/domain/index"
	"github.com/Banyango/gifoody_server/api/domain/tasks"
	"github.com/Banyango/gifoody_server/api/domain/users"
	"github.com/Banyango/gifoody_server/api/infrastructure/mail"
	"github.com/Banyango/gifoody_server/api/infrastructure/os"
	"github.com/Banyango/gifoody_server/api/infrastructure/pagination"
	"github.com/Banyango/gifoody_server/api/infrastructure/template"
	"github.com/Banyango/gifoody_server/api/repositories"
	"github.com/Banyango/gifoody_server/constants"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRouter(echo *echo.Echo, db *sqlx.DB) {
	mainGroup := echo.Group(constants.API_PATH)

	// repositories
	store := repositories.NewAppStore(db)

	// services
	osService := os.NewOSService()
	mailService := mail.NewMailService()
	templateService := template.NewTemplateService(osService)
	authService := users.NewUserAuthService(store.User())

	//controllers
	taskController := tasks.NewTaskController(store.Task(), authService)
	userController := users.NewUserController(store.User(), mailService, templateService, authService)
	indexController := index.NewIndexController()

	// index
	mainGroup.GET("index", indexController.GetIndex)

	// users
	mainGroup.POST("register", userController.PostRegister)
	mainGroup.POST("reset_password", userController.PostResetPassword)
	mainGroup.POST("confirm_reset_password", userController.PostConfirmResetPassword)
	mainGroup.POST("confirm_account", userController.PostConfirmAccount)
	mainGroup.POST("login", userController.PostLogin)
	mainGroup.POST("logout", userController.PostLogout)

	restrictedGroup := mainGroup.Group("auth/")

	jwtSecret := osService.GetEnv("API_JWT_SECRET")
	restrictedGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:[]byte(jwtSecret),
		TokenLookup:"cookie:refresh_token",
	}))

	// users
	restrictedGroup.GET("me", userController.GetMe)

	// todo api needs to be
	// days
	// days/{id}/tasks
	// days/{id}/tasks etc...


	// tasks
	restrictedGroup.GET("tasks", taskController.ListTask, pagination.Paginate())
	restrictedGroup.GET("tasks/:id/items", taskController.ListItems, pagination.Paginate())
	// todo need task_id on body
	restrictedGroup.POST("tasks", taskController.CreateTask)
	restrictedGroup.PUT("tasks/:id", taskController.UpdateTask)
	restrictedGroup.DELETE("tasks/:id", taskController.DeleteTask)
	restrictedGroup.POST("tasks/:id/summaries", taskController.CreateSummary)
	restrictedGroup.POST("tasks/:id/subtasks", taskController.CreateSubTask)
}
