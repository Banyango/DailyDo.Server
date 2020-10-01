package tasks

import (
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
// @Description Get a paginated list of Parent.
// @Accept json
// @Produce json
// @Param limit query string false "pagination limit"
// @Param offset query string false "pagination limit"
// @Success 200 {object} model.PagedResult
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Parent/ [get]
func (self *TaskController) ListTask(c echo.Context) (err error) {
	page := c.Get("Pagination").(pagination.Pagination)

	taskRequest := <-self.taskRepository.GetTaskAsync(model.TaskQuery{Limit: page.Limit, Offset: page.Offset, Type: "Parent"})
	if taskRequest.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, taskRequest.Err.Error())
	}

	pagedResult := pagination.NewPagedResult("Parent", c, taskRequest)

	return c.JSON(http.StatusOK, pagedResult)
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

	taskRequest := <-self.taskRepository.GetTaskByIdAsync(id)
	if taskRequest.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, taskRequest.Err.Error())
	}

	return c.JSON(http.StatusOK, taskRequest.Data.(model.Task))
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

	taskRequest := <-self.taskRepository.GetChildrenByTaskIdAsync(id)
	if taskRequest.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, taskRequest.Err.Error())
	}

	return c.JSON(http.StatusOK, collection.Collection{Items:taskRequest.Data})
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

	taskRequest := <-self.taskRepository.GetTasksByParentAsync(id)
	if taskRequest.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, taskRequest.Err.Error())
	}

	return c.JSON(http.StatusOK, collection.Collection{Items: taskRequest.Data})
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

	taskById := <-self.taskRepository.GetTaskByIdAsync(request.Parent)
	if taskById.Err != nil {
		return utils.LogError(err, http.StatusNotFound, fmt.Sprintf("Parent id={%s} not found", request.Parent))
	}

	task := model.Task{
		ID:        uuid.New().String(),
		Type:      "Task",
		Text:      null.NewString(request.Text, true),
		Order:     request.Order,
		Completed: request.Completed,
		TaskID:    null.NewString(request.Parent, true),
		UserID:    user.Id,
	}

	result := self.taskRepository.Save(task)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error saving task")
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
func (self *TaskController) CreateSubTask(c echo.Context) (err error) {
	request, err := NewCreateTaskRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	id := c.Param("id")
	taskById := <-self.taskRepository.GetTaskByIdAsync(id)
	if taskById.Err != nil {
		return utils.LogError(err, http.StatusNotFound, fmt.Sprintf("Parent id={%s} not found", id))
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	task := model.Task{
		ID:        uuid.New().String(),
		Type:      "SubTask",
		Text:      null.NewString(request.Text, true),
		TaskID:    null.NewString(id, true),
		Order:     request.Order,
		Completed: request.Completed,
		UserID:    user.Id,
	}

	result := self.taskRepository.Save(task)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error saving sub-task")
	}

	return c.JSON(http.StatusCreated, task)
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

	id := c.Param("id")
	taskById := <-self.taskRepository.GetTaskByIdAsync(id)
	if taskById.Err != nil {
		return utils.LogError(err, http.StatusNotFound, fmt.Sprintf("Parent id={%s} not found", id))
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	task := model.Task{
		ID:        uuid.New().String(),
		Type:      "Summary",
		Text:      null.NewString(request.Text, true),
		TaskID:    null.NewString(id, true),
		Order:     request.Order,
		Completed: request.Completed,
		UserID:    user.Id,
	}

	result := self.taskRepository.Save(task)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error saving sub-task")
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

	taskById := <-self.taskRepository.GetTaskByIdAsync(id)
	if taskById.Err != nil {
		return utils.LogError(taskById.Err, http.StatusNotFound, "Parent not found")
	}

	result := self.taskRepository.Delete(id)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error deleting task")
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

	taskById := <-self.taskRepository.GetTaskByIdAsync(request.ID)
	if taskById.Err != nil {
		return utils.LogError(taskById.Err, http.StatusNotFound, fmt.Sprintf("Parent id={%s} not found", request.ID))
	}

	task := taskById.Data.(model.Task)
	task.Text = null.NewString(request.Text, true)
	task.Order = request.Order
	task.Completed = request.Completed

	result := <-self.taskRepository.UpdateAsync(task)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error deleting task")
	}

	return c.JSON(http.StatusOK, task)
}
