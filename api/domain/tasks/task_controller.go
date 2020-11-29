package tasks

import (
	"context"
	"fmt"
	"github.com/Banyango/gifoody_server/api/domain/users"
	"github.com/Banyango/gifoody_server/api/infrastructure/collection"
	"github.com/Banyango/gifoody_server/api/infrastructure/pagination"
	"github.com/Banyango/gifoody_server/api/infrastructure/utils"
	"github.com/Banyango/gifoody_server/api/model"
	. "github.com/Banyango/gifoody_server/api/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
	"net/http"
)

type TaskController struct {
	taskRepository ITaskRepository
	authService    users.IUserAuthService
}

func NewTaskController(taskRepository ITaskRepository, authService users.IUserAuthService) *TaskController {
	return &TaskController{taskRepository: taskRepository, authService: authService}
}

// @Summary List Parent.
// @Description Get a paginated list of tasks.
// @Accept json
// @Produce json
// @Param limit query string false "pagination limit"
// @Param offset query string false "pagination limit"
// @Success 200 {object} model.PagedResult
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Parent/ [get]
func (self *TaskController) ListTask(ec echo.Context) (err error) {
	page := ec.Get("Pagination").(pagination.Pagination)

	var pagedResult pagination.PagedResult
	err = self.taskRepository.Execute(ec.Request().Context(), func(c context.Context) error {
		taskRequest := <-self.taskRepository.GetTaskAsync(model.TaskQuery{Limit: page.Limit, Offset: page.Offset, Type: "Parent"}, c)
		if taskRequest.Err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, taskRequest.Err.Error())
		}

		pagedResult = pagination.NewPagedResult("Parent", ec, taskRequest)

		return nil
	})

	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, pagedResult)
}

// @Summary Get Task.
// @Description Get Task by Id.
// @Accept json
// @Produce json
// @Success 200 {object} model.Task
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Task/:id [get]
func (self *TaskController) GetTask(c echo.Context) (err error) {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id not defined.")
	}

	var task model.Task
	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {
		taskRequest := <-self.taskRepository.GetTaskByIdAsync(id, c)
		if taskRequest.Err != nil {
			return taskRequest.Err
		}
		task = taskRequest.Data.(model.Task)
		return nil
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, task)
}

// @Summary List child items.
// @Description Get a paginated list of children items.
// @Accept json
// @Produce json
// @Param limit query string false "pagination limit"
// @Param offset query string false "pagination limit"
// @Success 200 {object} model.PagedResult
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Parent/{id}/items/ [get]
func (self *TaskController) ListItems(c echo.Context) (err error) {
	id := c.Param("id")

	var data interface{}
	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {
		taskRequest := <-self.taskRepository.GetChildrenByTaskIdAsync(id, c)
		if taskRequest.Err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, taskRequest.Err.Error())
		}

		data = taskRequest.Data

		return nil
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, collection.Collection{Items: data})
}

// @Summary List child items.
// @Description Get a paginated list of children items.
// @Accept json
// @Produce json
// @Param limit query string false "pagination limit"
// @Param offset query string false "pagination limit"
// @Success 200 {object} model.PagedResult
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Parent/{id}/tasks/ [get]
func (self *TaskController) ListTasks(c echo.Context) (err error) {

	id := c.Param("id")

	var data interface{}
	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {
		taskRequest := <-self.taskRepository.GetTasksByParentAsync(id, c)
		if taskRequest.Err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, taskRequest.Err.Error())
		}

		data = taskRequest.Data

		return nil
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, collection.Collection{Items: data})
}

// @Summary Create Parent.
// @Description Create a Parent.
// @Accept json
// @Produce json
// @Param limit query string false "pagination limit"
// @Param offset query string false "pagination limit"
// @Success 201
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Tasks/ [post]
func (self *TaskController) CreateTask(c echo.Context) (err error) {
	request, err := NewCreateTaskRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad reqest")
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	task := model.Task{
		ID:        uuid.New().String(),
		Type:      "Task",
		Text:      null.NewString(request.Text, true),
		Order:     0,
		Completed: request.Completed,
		TaskID:    null.NewString(request.Parent, true),
		UserID:    user.Id,
	}

	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {
		parentTaskQuery := <-self.taskRepository.GetTaskByIdAsync(request.Parent, c)
		if parentTaskQuery.Err != nil {
			return utils.LogError(parentTaskQuery.Err, http.StatusNotFound, fmt.Sprintf("Parent id={%s} not found", request.Parent))
		}

		maxOrderQuery := <-self.taskRepository.GetMaxOrder(request.Parent, c)
		if maxOrderQuery.Err != nil {
			return utils.LogError(maxOrderQuery.Err, http.StatusInternalServerError, "Error with max order")
		}

		if max := maxOrderQuery.Data.(*int); max != nil {
			task.Order = *max + 1
		}

		result := self.taskRepository.Save(task, c)
		if result.Err != nil {
			return utils.LogError(result.Err, http.StatusInternalServerError, "Error saving task")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, task)
}

// @Summary Create SubTask.
// @Description Create a SubTask.
// @Accept json
// @Param request body tasks.CreateTaskRequestFromContext true "request"
// @Produce json
// @Success 201
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/tasks/{id}/subtasks [post]
func (self *TaskController) CreateSubTask(ec echo.Context) (err error) {
	request, err := NewCreateTaskRequestFromContext(ec)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	user, err := self.authService.GetLoggedInUser(ec)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	var task model.Task
	err = self.taskRepository.Execute(ec.Request().Context(), func(c context.Context) error {
		id := ec.Param("id")
		taskById := <-self.taskRepository.GetTaskByIdAsync(id, c)
		if taskById.Err != nil {
			return utils.LogError(err, http.StatusNotFound, fmt.Sprintf("Parent id={%s} not found", id))
		}

		task = model.Task{
			ID:        uuid.New().String(),
			Type:      "SubTask",
			Text:      null.NewString(request.Text, true),
			TaskID:    null.NewString(id, true),
			Order:     0,
			Completed: request.Completed,
			UserID:    user.Id,
		}

		maxOrderQuery := <-self.taskRepository.GetMaxOrder(id, c)
		if maxOrderQuery.Err != nil {
			return utils.LogError(maxOrderQuery.Err, http.StatusInternalServerError, "Error with max order")
		}

		if max := maxOrderQuery.Data.(*int); max != nil {
			task.Order = *max + 1
		}

		result := self.taskRepository.Save(task, c)
		if result.Err != nil {
			return utils.LogError(result.Err, http.StatusInternalServerError, "Error saving sub-task")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusCreated, task)
}

// @Summary Create Summary.
// @Description Create a Summary.
// @Accept json
// @Param request body tasks.CreateTaskRequestFromContext true "request"
// @Produce json
// @Success 201
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Parent/{id}/summaries [post]
func (self *TaskController) CreateSummary(c echo.Context) (err error) {
	request, err := NewCreateTaskRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	id := c.Param("id")

	task := model.Task{
		ID:        uuid.New().String(),
		Type:      "Summary",
		Text:      null.NewString(request.Text, true),
		TaskID:    null.NewString(id, true),
		Order:     0,
		Completed: request.Completed,
		UserID:    user.Id,
	}

	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {
		taskById := <-self.taskRepository.GetTaskByIdAsync(id, c)
		if taskById.Err != nil {
			return utils.LogError(err, http.StatusNotFound, fmt.Sprintf("Parent id={%s} not found", id))
		}

		maxOrderQuery := <-self.taskRepository.GetMaxOrder(id, c)
		if maxOrderQuery.Err != nil {
			return utils.LogError(maxOrderQuery.Err, http.StatusInternalServerError, "Error with max order")
		}

		if max := maxOrderQuery.Data.(*int); max != nil {
			task.Order = *max
		}

		result := self.taskRepository.Save(task, c)
		if result.Err != nil {
			return utils.LogError(result.Err, http.StatusInternalServerError, "Error saving sub-task")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, task)
}

// @Summary Delete Parent.
// @Description Delete a Parent.
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Parent/{id} [delete]
func (self *TaskController) DeleteTask(c echo.Context) (err error) {
	id := c.Param("id")

	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {
		taskById := <-self.taskRepository.GetTaskByIdAsync(id, c)
		if taskById.Err != nil {
			return utils.LogError(taskById.Err, http.StatusNotFound, "Parent not found")
		}

		result := self.taskRepository.Delete(id, c)
		if result.Err != nil {
			return utils.LogError(result.Err, http.StatusInternalServerError, "Error deleting task")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// @Summary Update Parent.
// @Description Update a Parent.
// @Accept json
// @Produce json
// @Param request body tasks.UpdateTaskRequestFromContext true "request"// @Success 200
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/tasks/{id} [put]
func (self *TaskController) UpdateTask(c echo.Context) (err error) {

	request, err := NewUpdateTaskRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	var task model.Task
	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {
		taskById := <-self.taskRepository.GetTaskByIdAsync(request.ID, c)
		if taskById.Err != nil {
			return utils.LogError(taskById.Err, http.StatusNotFound, fmt.Sprintf("Parent id={%s} not found", request.ID))
		}

		task = taskById.Data.(model.Task)
		task.Text = null.NewString(request.Text, true)

		task.Completed = request.Completed

		result := <-self.taskRepository.UpdateAsync(task, c)
		if result.Err != nil {
			return utils.LogError(result.Err, http.StatusInternalServerError, "Error deleting task")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, task)
}

// @Summary Update the order of a task.
// @Description Update the order.
// @Accept json
// @Produce json
// @Success 201
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/tasks/:id/order [post]
func (self *TaskController) UpdateTaskOrder(c echo.Context) (err error) {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Id not defined.")
	}

	request, err := NewUpdateTaskOrderRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	if request.TaskId == request.NewParent {
		return echo.NewHTTPError(http.StatusBadRequest, "TaskId is equal to NewParent")
	}

	var taskToUpdate model.Task
	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {

		taskToUpdateQuery := <-self.taskRepository.GetTaskByIdAsync(id, c)
		if taskToUpdateQuery.Err != nil {
			return utils.LogError(err, http.StatusNotFound, "Not Found.")
		}

		taskToUpdate = taskToUpdateQuery.Data.(model.Task)

		// Get all tasks for parent.
		var tasks model.Tasks
		if taskToUpdate.Type == "Task" {
			childrenOfNewParentQuery := <-self.taskRepository.GetTasksByParentAsync(taskToUpdate.TaskID.String, c)
			if childrenOfNewParentQuery.Err != nil {
				return utils.LogError(err, http.StatusInternalServerError, err.Error())
			}
			tasks = model.Tasks{}.Append(childrenOfNewParentQuery.Data.([]model.Task)...)
		} else {
			childrenOfNewParentQuery := <-self.taskRepository.GetChildrenByTaskIdAsync(taskToUpdate.TaskID.String, c)
			if childrenOfNewParentQuery.Err != nil {
				return utils.LogError(err, http.StatusInternalServerError, err.Error())
			}
			tasks = model.Tasks{}.Append(childrenOfNewParentQuery.Data.([]model.Task)...)
		}

		tasks = tasks.Filter(func(task model.Task) bool {
			return task.ID != request.TaskId
		})

		index := -1
		for i, t := range tasks {
			if t.ID == request.NewParent {
				index = i
				break
			}
		}

		if taskToUpdate.Order <= index {
			tasks = tasks.Insert(index+1, taskToUpdate)
		} else {
			tasks = tasks.Insert(index, taskToUpdate)
		}

		for i, t := range tasks {
			t.Order = i
			err := <-self.taskRepository.UpdateAsync(t, c)
			if err.Err != nil {
				return utils.LogError(err.Err, http.StatusInternalServerError, "error updating order")
			}
		}

		return nil
	})

	return c.JSON(http.StatusOK, taskToUpdate)
}
