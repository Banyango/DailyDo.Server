package tasks

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type SubTasksByTaskIdRequest struct {
	id string
}

func NewSubTasksByTaskIdRequestFromContext(c echo.Context) (request *SubTasksByTaskIdRequest, err error) {
	request.id = c.Param("id")
	return request, nil
}

type CreateTaskRequest struct {
	Parent    string `json:"parent"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
	Order     string `json:"order"`
}

func NewCreateTaskRequestFromContext(c echo.Context) (request *CreateTaskRequest, err error) {
	request = new(CreateTaskRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	// todo sanitize html and text!!!!!!!
	if request.Parent == "" {
		return nil, fmt.Errorf("Parent cannot be null")
	}

	return request, nil
}

type UpdateTaskRequest struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Parent    string `json:"parent"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
	Order     string `json:"order"`
}

func NewUpdateTaskRequestFromContext(c echo.Context) (request *UpdateTaskRequest, err error) {
	request = new(UpdateTaskRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	// todo sanitize html and text!!!!!!!

	return request, nil
}
