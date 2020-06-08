package api

import "github.com/labstack/echo"

func InitRouter() {
	mainGroup := echo.Group("/api")

	postController := controllers.NewPostController()

	//posts
	mainGroup.GET("posts", postController.GetPosts)
}
