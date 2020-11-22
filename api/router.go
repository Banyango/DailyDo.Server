package api

import (
	"github.com/Banyango/gifoody_server/api/domain/days"
	"github.com/Banyango/gifoody_server/api/domain/index"
	"github.com/Banyango/gifoody_server/api/domain/tasks"
	"github.com/Banyango/gifoody_server/api/domain/users"
	"github.com/Banyango/gifoody_server/api/infrastructure/mail"
	"github.com/Banyango/gifoody_server/api/infrastructure/os"
	"github.com/Banyango/gifoody_server/api/infrastructure/pagination"
	"github.com/Banyango/gifoody_server/api/infrastructure/template"
	"github.com/Banyango/gifoody_server/api/infrastructure/time"
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
	timeService := time.NewTimeService()

	//controllers
	indexController := index.NewIndexController()
	dayController := days.NewDayController(timeService, store.Day(), authService)
	dayTaskController := days.NewDayTaskController(timeService, store.Day(), store.Task(), authService)
	taskController := tasks.NewTaskController(store.Task(), authService)
	taskOrderController := tasks.NewTaskOrderController(store.Task(), authService)
	userController := users.NewUserController(store.User(), mailService, templateService, authService)

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
		SigningKey:  []byte(jwtSecret),
		TokenLookup: "cookie:refresh_token",
	}))

	// users
	restrictedGroup.GET("me", userController.GetMe)

	// days
	restrictedGroup.GET("days", dayController.ListDays, pagination.Paginate())
	restrictedGroup.POST("days", dayController.CreateDay)
	restrictedGroup.PUT("days/:id", dayController.UpdateDay)
	restrictedGroup.DELETE("days/:id", dayController.DeleteDay)
	restrictedGroup.GET("days/:id/tasks", dayTaskController.ListTasksForDay)

	// tasks
	restrictedGroup.GET("tasks", taskController.ListTask, pagination.Paginate())
	restrictedGroup.GET("tasks/:id", taskController.GetTask)
	restrictedGroup.GET("tasks/:id/tasks", taskController.ListTasks)
	restrictedGroup.GET("tasks/:id/items", taskController.ListItems, pagination.Paginate())

	// todo need task_id on body
	restrictedGroup.POST("tasks", taskController.CreateTask)
	restrictedGroup.PUT("tasks/:id", taskController.UpdateTask)
	restrictedGroup.DELETE("tasks/:id", taskController.DeleteTask)
	restrictedGroup.POST("tasks/:id/summaries", taskController.CreateSummary)
	restrictedGroup.POST("tasks/:id/subtasks", taskController.CreateSubTask)

	// task order
	restrictedGroup.POST("tasks/:id/order", taskOrderController.UpdateTaskOrder)
}
