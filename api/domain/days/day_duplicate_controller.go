package days

import (
	"context"
	"github.com/Banyango/gifoody_server/api/domain/users"
	"github.com/Banyango/gifoody_server/api/infrastructure/time"
	"github.com/Banyango/gifoody_server/api/infrastructure/utils"
	"github.com/Banyango/gifoody_server/api/model"
	"github.com/Banyango/gifoody_server/api/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
	"net/http"
)

type DayDuplicateController struct {
	timeService    time.ITimeInterface
	authService    users.IUserAuthService
	dayRepository  repositories.IDayRepository
	taskRepository repositories.ITaskRepository
}

func NewDayDuplicateController(
	timeService time.ITimeInterface,
	dayRepository repositories.IDayRepository,
	taskRepository repositories.ITaskRepository,
	authService users.IUserAuthService) *DayDuplicateController {
	return &DayDuplicateController{
		timeService:    timeService,
		dayRepository:  dayRepository,
		taskRepository: taskRepository,
		authService:    authService}
}


// @Summary Create Day.
// @Description Create a Day.
// @Accept json
// @Produce json
// @Success 201
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/days/:id/duplicate [post]
func (self *DayDuplicateController) DuplicateDay(c echo.Context) (err error) {
	request, err := NewDuplicateDayRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad Request")
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	var dayResult model.Day
	err = self.taskRepository.Execute(c.Request().Context(), func(c context.Context) error {

		dayByIdQuery := <-self.dayRepository.GetDayByIdAsync(request.id)
		if dayByIdQuery.Err != nil {
			return utils.LogError(dayByIdQuery.Err, http.StatusNotFound, "Day Not Found")
		}
		dayToCopy := dayByIdQuery.Data.(model.Day)

		day := model.Day{
			ID:           uuid.New().String(),
			Summary:      null.NewString(dayToCopy.Summary.String, true),
			Date:         self.timeService.GetStartOfDayUTC(),
			UserID:       user.Id,
		}

		result := self.dayRepository.Save(day)
		if result.Err != nil {
			return utils.LogError(result.Err, http.StatusInternalServerError, "Error saving day")
		}

		dayResult = result.Data.(model.Day)

		tasksQuery := <-self.taskRepository.GetTasksByParentAsync(dayToCopy.ParentTaskID, c)
		if tasksQuery.Err != nil {
			return utils.LogError(tasksQuery.Err, http.StatusInternalServerError, tasksQuery.Err.Error())
		}

		for _, task := range tasksQuery.Data.([]model.Task) {

			if task.Completed {
				continue
			}

			newId, err := self.CreateTask(task.Text, dayResult.ParentTaskID, user, c)
			if err != nil  {
				return utils.LogError(err, http.StatusInternalServerError, err.Error())
			}

			childrenQuery := <-self.taskRepository.GetChildrenByTaskIdAsync(task.ID, c)
			if childrenQuery.Err != nil {
				return utils.LogError(childrenQuery.Err, http.StatusInternalServerError, childrenQuery.Err.Error())
			}

			for _, subTask := range childrenQuery.Data.([]model.Task) {
				if subTask.Completed {
					continue
				}
				_, err = self.CreateTask(subTask.Text, newId, user, c)
				if err != nil  {
					return utils.LogError(err, http.StatusInternalServerError, err.Error())
				}
			}
		}

		return nil
	})

	if err != nil {
		return utils.LogError(err, http.StatusInternalServerError, "could not duplicate day.")
	}

	return c.JSON(http.StatusCreated, &dayResult)
}

func (self *DayDuplicateController) CreateTask(text null.String, parentTaskID string, user model.User, c context.Context) (id string, err error ) {
	task := model.Task{
		ID:        uuid.New().String(),
		Type:      "Task",
		Text:      text,
		Order:     0,
		Completed: false,
		TaskID:    null.NewString(parentTaskID, true),
		UserID:    user.Id,
	}
	maxOrderQuery := <-self.taskRepository.GetMaxOrder(parentTaskID, c)
	if maxOrderQuery.Err != nil {
		return "nil", utils.LogError(maxOrderQuery.Err, http.StatusInternalServerError, "Error with max order")
	}
	if max := maxOrderQuery.Data.(*int); max != nil {
		task.Order = *max + 1
	}
	result := self.taskRepository.Save(task, c)
	if result.Err != nil {
		return "", utils.LogError(result.Err, http.StatusInternalServerError, "Error saving task")
	}
	return task.ID, nil
}