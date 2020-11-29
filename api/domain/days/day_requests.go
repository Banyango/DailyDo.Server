package days

import (
	"github.com/labstack/echo/v4"
)

type CreateDayRequest struct {
	Summary string `json:"summary"`
}

func NewCreateDayRequestFromContext(c echo.Context) (request *CreateDayRequest, err error) {
	request = new(CreateDayRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	// todo sanitize html and text!!!!!!!

	return request, nil
}

type UpdateDayRequest struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

func NewUpdateDayRequestFromContext(c echo.Context) (request *UpdateDayRequest, err error) {
	request = new(UpdateDayRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	// todo sanitize html and text!!!!!!!

	return request, nil
}

type DuplicateDayRequest struct {
	id string
}

func NewDuplicateDayRequestFromContext(c echo.Context) (request *DuplicateDayRequest, err error) {
	request = new(DuplicateDayRequest)
	request.id = c.Param("id")
	return request, nil
}