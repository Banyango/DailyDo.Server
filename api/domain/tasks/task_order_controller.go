package tasks

import (
	"fmt"
	"github.com/Banyango/gifoody_server/api/domain/users"
	"github.com/Banyango/gifoody_server/api/infrastructure/collection"
	"github.com/Banyango/gifoody_server/api/infrastructure/utils"
	"github.com/Banyango/gifoody_server/api/model"
	"github.com/Banyango/gifoody_server/api/repositories"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"gopkg.in/guregu/null.v4"
	"net/http"
)

type TaskOrderController struct {
	authService    users.IUserAuthService
	taskRepository repositories.ITaskRepository
}

func NewTaskOrderController(
	taskRepository repositories.ITaskRepository,
	authService users.IUserAuthService) *TaskOrderController {
	return &TaskOrderController{
		taskRepository: taskRepository,
		authService:    authService}
}

// @Summary Get Order.
// @Description Get Children by Parent Id.
// @Accept json
// @Produce json
// @Success 200 {object} []model.Task
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Task/:id/order [get]
func (self *TaskOrderController) GetOrder(c echo.Context) (err error) {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id not defined.")
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	var tasks []model.Task
	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {
		tasksQuery := <-self.taskRepository.GetTaskByIdAsync(id, c)
		if tasksQuery.Err != nil {
			return utils.LogError(tasksQuery.Err, http.StatusNotFound, tasksQuery.Err.Error())
		}

		task := tasksQuery.Data.(model.Task)
		if task.UserID != user.Id {
			e := fmt.Errorf("bad request")
			return utils.LogError(e, http.StatusBadRequest, e.Error())
		}

		childrenQuery := <-self.taskRepository.GetChildrenByTaskIdAsync(id, c)
		if childrenQuery.Err != nil {
			return utils.LogError(childrenQuery.Err, http.StatusInternalServerError, childrenQuery.Err.Error())
		}

		tasks = childrenQuery.Data.([]model.Task)

		return nil
	})

	return c.JSON(http.StatusOK, collection.Collection{Items: tasks})
}

// @Summary Update the order of a task.
// @Description Update the order.
// @Accept json
// @Produce json
// @Success 201
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/tasks/:id/order [post]
func (self *TaskOrderController) UpdateTaskOrder(c echo.Context) (err error) {
	request, err := NewUpdateTaskOrderRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	var taskToUpdate model.Task
	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {

		if request.TaskId == request.NewParent {
			return utils.LogError(fmt.Errorf("bad request: id cannot be equal to new parent"), http.StatusBadRequest, "Bad Request: task id cannot be equal to new parent.")
		}

		// Get the task.
		taskToUpdateQuery := <-self.taskRepository.GetTaskByIdAsync(request.TaskId, c)
		if taskToUpdateQuery.Err != nil {
			return utils.LogError(taskToUpdateQuery.Err, http.StatusNotFound, "Task not found.")
		}

		// Get new parent task.
		newParentQuery := <-self.taskRepository.GetTaskByIdAsync(request.NewParent, c)
		if newParentQuery.Err != nil {
			return utils.LogError(newParentQuery.Err, http.StatusNotFound, "New Parent ID not found.")
		}

		// Get any tasks that have task as a parent.
		childrenOfTaskToUpdateQuery := <-self.taskRepository.GetTaskByOrderIdAsync(taskToUpdate.TaskID.String, request.TaskId, c)

		taskToUpdate = taskToUpdateQuery.Data.(model.Task)
		taskNewParent := newParentQuery.Data.(model.Task)

		if taskToUpdate.UserID != user.Id || taskNewParent.UserID != user.Id {
			return utils.LogError(fmt.Errorf("bad request"), http.StatusBadRequest, "Bad Request.")
		}

		if childrenOfTaskToUpdateQuery.Data != nil {
			// Update task after to current task order parent.
			taskAfter := childrenOfTaskToUpdateQuery.Data.(model.Task)
			taskAfter.Order = taskToUpdate.Order

			result := <-self.taskRepository.UpdateAsync(taskAfter, c)
			if result.Err != nil {
				return utils.LogError(result.Err, http.StatusInternalServerError, "Something went wrong.")
			}
		}

		// Get any tasks that have new parent as order parent
		childrenOfNewParentQuery := <-self.taskRepository.GetTaskByOrderIdAsync(taskToUpdate.TaskID.String, request.NewParent, c)
		if childrenOfNewParentQuery.Data != nil {
			taskToReParent := childrenOfNewParentQuery.Data.(model.Task)

			// Update them to point to task we are updating
			taskToReParent.Order = null.NewString(taskToUpdate.ID, true)

			result := <-self.taskRepository.UpdateAsync(taskToReParent, c)
			if result.Err != nil {
				return utils.LogError(result.Err, http.StatusInternalServerError, "Something went wrong.")
			}
		}

		taskToUpdate.Order = null.NewString(request.NewParent, true)

		result := <-self.taskRepository.UpdateAsync(taskToUpdate, c)
		if result.Err != nil {
			return utils.LogError(result.Err, http.StatusInternalServerError, "Something went wrong.")
		}

		return nil
	})

	return self.GetOrder(c)
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
