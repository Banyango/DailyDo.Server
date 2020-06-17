package api

import (
	"database/sql"
	"github.com/Banyango/gifoody_server/api/controllers"
	"github.com/Banyango/gifoody_server/constants"
	"github.com/Banyango/gifoody_server/middleware"
	"github.com/Banyango/gifoody_server/repositories"
	"github.com/labstack/echo/v4"
)

func InitRouter(echo *echo.Echo, db *sql.DB) {
	mainGroup := echo.Group(constants.API_PATH)

	// repositories
	store := repositories.NewAppStore(db)

	//controllers
	postController := controllers.NewPostController(store.Post())
	indexController := controllers.NewIndexController()

	// index
	mainGroup.GET("index", indexController.GetIndex)

	// posts
	mainGroup.GET("posts", postController.ListPosts, middleware.Paginate())
}
