package tasks

import (
	"fmt"
	"github.com/Banyango/gifoody_server/api/domain/users"
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
	authService users.IUserAuthService
}

func NewTaskController(taskRepository ITaskRepository, authService users.IUserAuthService) *TaskController {
	return &TaskController{taskRepository: taskRepository, authService:authService}
}

// @Summary List Task.
// @Description Get a paginated list of Task.
// @Accept json
// @Produce json
// @Param limit query string false "pagination limit"
// @Param offset query string false "pagination limit"
// @Success 200 {object} model.PagedResult
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Task/ [get]
func (self *TaskController) ListTask(c echo.Context) (err error) {
	page := c.Get("Pagination").(pagination.Pagination)

	taskRequest := <-self.taskRepository.GetTaskAsync(model.TaskQuery{Limit:page.Limit, Offset:page.Offset, Type:"task"})
	if taskRequest.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, taskRequest.Err.Error())
	}

	pagedResult := pagination.NewPagedResult("Task", c, taskRequest)

	return c.JSON(http.StatusOK, pagedResult)
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
// @Router /v1/Task/{id}/items/ [get]
func (self *TaskController) ListItems(c echo.Context) (err error) {

	id := c.Param("id")

	page := c.Get("Pagination").(pagination.Pagination)

	taskRequest := <-self.taskRepository.GetChildrenByTaskIdAsync(id, page.Limit, page.Offset)
	if taskRequest.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, taskRequest.Err.Error())
	}

	pagedResult := pagination.NewPagedResult("Items", c, taskRequest)

	return c.JSON(http.StatusOK, pagedResult)
}

// @Summary Create Task.
// @Description Create a Task.
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

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	task := model.Task{
		ID:uuid.New().String(),
		Type:"task",
		Text:null.NewString(request.Text,true),
		Order:request.Order,
		Completed:request.Completed,
		UserID:user.Id,
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
// @Router /v1/Task/{id}/subtasks [post]
func (self *TaskController) CreateSubTask(c echo.Context) (err error) {
	request, err := NewCreateTaskRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	id := c.Param("id")
	taskById := <- self.taskRepository.GetTaskByIdAsync(id)
	if taskById.Err != nil {
		return utils.LogError(err, http.StatusNotFound, fmt.Sprintf("Task id={%s} not found", id))
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	task := model.Task{
		ID:uuid.New().String(),
		Type:"subtask",
		Text:null.NewString(request.Text,true),
		TaskID:null.NewString(id,true),
		Order:request.Order,
		Completed:request.Completed,
		UserID:user.Id,
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
// @Router /v1/Task/{id}/summaries [post]
func (self *TaskController) CreateSummary(c echo.Context) (err error) {
	request, err := NewCreateTaskRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	id := c.Param("id")
	taskById := <- self.taskRepository.GetTaskByIdAsync(id)
	if taskById.Err != nil {
		return utils.LogError(err, http.StatusNotFound, fmt.Sprintf("Task id={%s} not found", id))
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	task := model.Task{
		ID:uuid.New().String(),
		Type:"summary",
		Text:null.NewString(request.Text,true),
		TaskID:null.NewString(id,true),
		Order:request.Order,
		Completed:request.Completed,
		UserID:user.Id,
	}

	result := self.taskRepository.Save(task)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error saving sub-task")
	}

	return c.JSON(http.StatusCreated, task)
}

// @Summary Delete Task.
// @Description Delete a Task.
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Task/{id} [delete]
func (self *TaskController) DeleteTask(c echo.Context) (err error) {
	id := c.Param("id")

	taskById := <- self.taskRepository.GetTaskByIdAsync(id)
	if taskById.Err != nil {
		return utils.LogError(taskById.Err, http.StatusNotFound, "Task not found")
	}

	result := self.taskRepository.Delete(id)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error deleting task")
	}

	return c.NoContent(http.StatusOK)
}

// @Summary Update Task.
// @Description Update a Task.
// @Accept json
// @Produce json
// @Param request body tasks.UpdateTaskRequestFromContext true "request"// @Success 200
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Task/{id} [put]
func (self *TaskController) UpdateTask(c echo.Context) (err error) {

	request, err := NewUpdateTaskRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	taskById := <- self.taskRepository.GetTaskByIdAsync(request.ID)
	if taskById.Err != nil {
		return utils.LogError(taskById.Err, http.StatusNotFound, fmt.Sprintf("Task id={%s} not found", request.ID))
	}

	task := taskById.Data.(model.Task)
	task.Text = null.NewString(request.Text,true)
	task.Order = request.Order
	task.Completed = request.Completed

	result := <- self.taskRepository.UpdateAsync(task)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error deleting task")
	}

	return c.JSON(http.StatusOK, task)
}