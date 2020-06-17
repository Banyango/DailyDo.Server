package controllers

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

type IndexController struct {
}

func NewIndexController() *IndexController {
	return &IndexController{

	}
}

// @Summary Get Indexs.
// @Description Get index
// @Accept json
// @Produce json
// @Success 200 {object} controller.IndexResponse
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/index/ [get]
func (self *IndexController) GetIndex(c echo.Context) (err error) {
	response, err := NewIndexResponse(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}
	return c.JSON(http.StatusOK, response)
}

