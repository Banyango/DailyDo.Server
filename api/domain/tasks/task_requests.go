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
	Parent    string   `json:"parent"`
	Text      string   `json:"text"`
	Completed bool     `json:"completed"`
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
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	Parent    string   `json:"parent"`
	Text      string   `json:"text"`
	Completed bool     `json:"completed"`
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

type UpdateTaskOrderRequest struct {
	TaskId    string `json:"id"`
	NewParent string `json:"newParent"`
}

func NewUpdateTaskOrderRequestFromContext(c echo.Context) (request *UpdateTaskOrderRequest, err error) {
	request = new(UpdateTaskOrderRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}
