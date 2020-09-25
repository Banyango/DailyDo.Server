package tasks

import (
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
	Type      string `json:"type" db:"type"`
	Task      string `json:"task" db:"task_id"`
	Text      string `json:"text" db:"text"`
	Completed bool   `json:"completed" db:"completed"`
	Order     int    `json:"order" db:"order"`
}

func NewCreateTaskRequestFromContext(c echo.Context) (request *CreateTaskRequest, err error) {
	request = new(CreateTaskRequest)
	err = c.Bind(request)
	if err != nil {
		return nil, err
	}

	// todo sanitize html and text!!!!!!!

	return request, nil
}

type UpdateTaskRequest struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Task      string `json:"task"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
	Order     int    `json:"order"`
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