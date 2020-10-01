package pagination

import (
	"fmt"
	"github.com/Banyango/gifoody_server/api/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Page struct {
	Number int
	Url    string
}

type Pagination struct {
	Offset           int         `json:"offset"`
	Limit            int         `json:"limit"`
	HasPreviousPages bool        `json:"hasPrevious"`
	HasNextPages     bool        `json:"hasNext"`
	Data             interface{} `json:"data"`
	Next             model.Link  `json:"next"`
	Previous         model.Link  `json:"previous"`
}

// https://github.com/expressjs/express-paginate/blob/master/index.js

const PAGINATION = "Pagination"

func Paginate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			offset, err := getNumber("offset", 0, c)
			limit, err := getNumber("limit", 5, c)
			total, err := getNumber("total", 0, c)
			hasPrev := offset > 1
			hasNext := offset < total

			pagination := Pagination{
				Offset:           offset,
				Limit:            limit,
				HasPreviousPages: hasPrev,
				HasNextPages:     hasNext,
			}

			c.Set(PAGINATION, pagination)

			return next(c)
		}
	}
}

func getNumber(queryParam string, defaultValue int, c echo.Context) (number int, err error) {
	param := c.QueryParam(queryParam)

	value := defaultValue
	if param != "" {
		value, err = strconv.Atoi(param)
		if err != nil {
			return -1, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("{0} invalid.", queryParam))
		}
	}

	return value, nil
}
