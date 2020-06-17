package controllers

import (
	"github.com/Banyango/gifoody_server/api/services"
	"net/http"

	"github.com/Banyango/gifoody_server/middleware"
	"github.com/Banyango/gifoody_server/model"
	. "github.com/Banyango/gifoody_server/repositories"
	"github.com/labstack/echo/v4"
)

type PostController struct {
	postRepository IPostRepository
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
// @Success 200 {object} model.PagedResult
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/posts/ [get]
func (self *PostController) ListPosts(c echo.Context) (err error) {
	page := c.Get("Pagination").(middleware.Pagination)

	postRequest := <-self.postRepository.FindPosts(model.PostQuery{Limit:page.Limit, Offset:page.Offset})
	if postRequest.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, postRequest.Err.Error())
	}

	pagedResult := services.NewPagedResult("posts", c, postRequest)

	return c.JSON(http.StatusOK, pagedResult)
}
