package api

import (
	"github.com/Banyango/dailydo_server/api/domain/days"
	"github.com/Banyango/dailydo_server/api/domain/index"
	"github.com/Banyango/dailydo_server/api/domain/tasks"
	"github.com/Banyango/dailydo_server/api/domain/users"
	"github.com/Banyango/dailydo_server/api/infrastructure/mail"
	"github.com/Banyango/dailydo_server/api/infrastructure/os"
	"github.com/Banyango/dailydo_server/api/infrastructure/pagination"
	"github.com/Banyango/dailydo_server/api/infrastructure/template"
	"github.com/Banyango/dailydo_server/api/infrastructure/time"
	"github.com/Banyango/dailydo_server/api/repositories"
	"github.com/Banyango/dailydo_server/constants"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

var (
	Unauthorized = echo.NewHTTPError(http.StatusUnauthorized, nil)
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
	timeService := time.NewTimeService()

	//controllers
	indexController := index.NewIndexController()
	dayController := days.NewDayController(timeService, store.Day(), authService)
	dayTaskController := days.NewDayTaskController(timeService, store.Day(), store.Task(), authService)
	dayDuplicateController := days.NewDayDuplicateController(timeService, store.Day(), store.Task(), authService)
	taskController := tasks.NewTaskController(store.Task(), authService)
	userController := users.NewUserController(store.User(), mailService, templateService, authService)

	// index
	mainGroup.GET("index", indexController.GetIndex)

	// users
	// todo for now dont allow registration.
	//mainGroup.POST("register", userController.PostRegister)
	mainGroup.POST("reset_password", userController.PostResetPassword)
	mainGroup.POST("confirm_reset_password", userController.PostConfirmResetPassword)
	mainGroup.POST("confirm_account", userController.PostConfirmAccount)
	mainGroup.POST("login", userController.PostLogin)
	mainGroup.POST("logout", userController.PostLogout)

	restrictedGroup := mainGroup.Group("auth/")

	jwtSecret := osService.GetEnv("API_JWT_SECRET")
	restrictedGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(jwtSecret),
		TokenLookup: "cookie:refresh_token",
		ErrorHandler: func(e error) error {
			return Unauthorized
		},
	}))

	// users
	restrictedGroup.GET("me", userController.GetMe)

	// days
	restrictedGroup.GET("days", dayController.ListDays, pagination.Paginate())
	restrictedGroup.GET("days/:id/tasks", dayTaskController.ListTasksForDay)
	restrictedGroup.PUT("days/:id", dayController.UpdateDay)

	restrictedGroup.POST("days", dayController.CreateDay)
	restrictedGroup.POST("days/:id/duplicate", dayDuplicateController.DuplicateDay)

	restrictedGroup.DELETE("days/:id", dayController.DeleteDay)

	// tasks
	restrictedGroup.GET("tasks", taskController.ListTask, pagination.Paginate())
	restrictedGroup.GET("tasks/:id", taskController.GetTask)
	restrictedGroup.GET("tasks/:id/tasks", taskController.ListTasks)
	restrictedGroup.GET("tasks/:id/items", taskController.ListItems, pagination.Paginate())

	restrictedGroup.PUT("tasks/:id", taskController.UpdateTask)
	restrictedGroup.PUT("tasks/:id/order", taskController.UpdateTaskOrder)

	restrictedGroup.POST("tasks", taskController.CreateTask)
	restrictedGroup.POST("tasks/:id/summaries", taskController.CreateSummary)
	restrictedGroup.POST("tasks/:id/subtasks", taskController.CreateSubTask)

	restrictedGroup.DELETE("tasks/:id", taskController.DeleteTask)
}
