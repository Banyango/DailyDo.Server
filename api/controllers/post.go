package controllers

import (
	"net/http"

	"github.com/Banyango/gifoody_server/middleware"
	"github.com/labstack/echo"
)

type PostController struct {
	postRepository *IPostRepository
}

func NewPostController(postRepository IPostRepository) *PostController {
	return &PostController{postRepository: postRepository}
}

// @Summary List Posts.
// @Description Get a paginated list of posts.
// @Accept json
// @Produce json
// @Param limit query string false "pagination limit"
// @Param offset query string false "pagination limit"
// @Success 200 {object} model.Post
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/posts/ [get]
func (self *PostController) ListPosts(c echo.Context) (err error) {
	page := c.Get("Pagination").(middleware.Pagination)

	postRequest := <-self.postRepository.FindPost(page.Offset, page.Limit)
	if authorsRequest.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error finding posts")
	}

	page.Data = postRequest.Data
	page.SetLinks(/posts/)

	return c.JSON(http.StatusOK, page)
}
