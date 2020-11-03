package days

import (
	"github.com/Banyango/gifoody_server/api/domain/users"
	"github.com/Banyango/gifoody_server/api/infrastructure/time"
	"github.com/Banyango/gifoody_server/api/model"
	"github.com/Banyango/gifoody_server/api/repositories"
	"github.com/labstack/echo/v4"
	"net/http"
)

type DayTaskController struct {
	timeService    time.ITimeInterface
	authService    users.IUserAuthService
	dayRepository  repositories.IDayRepository
	taskRepository repositories.ITaskRepository
}

func NewDayTaskController(
	timeService time.ITimeInterface,
	dayRepository repositories.IDayRepository,
	taskRepository repositories.ITaskRepository,
	authService users.IUserAuthService) *DayTaskController {
	return &DayTaskController{
		timeService:    timeService,
		dayRepository:  dayRepository,
		taskRepository: taskRepository,
		authService:    authService}
}

// @Summary List Tasks for Day.
// @Description Get tasks for a day.
// @Accept json
// @Produce json
// @Param limit query string false "pagination limit"
// @Param offset query string false "pagination limit"
// @Success 200 {object} model.PagedResult
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/auth/days/:id/tasks [get]
func (self *DayTaskController) ListTasksForDay(c echo.Context) (err error) {
	_, err = self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	id := c.Param("id")

	dayRequest := <-self.dayRepository.GetDayByIdAsync(id)
	if dayRequest.Err != nil {
		return echo.NewHTTPError(http.StatusNotFound, dayRequest.Err.Error())
	}

	day := dayRequest.Data.(model.Day)

	tasks := <-self.taskRepository.GetTasksByParentAsync(day.ParentTaskID)
	if tasks.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dayRequest.Err.Error())
	}

	response := TasksByDayResponse{}
	for _, task := range tasks.Data.([]model.Task) {

		children := <-self.taskRepository.GetChildrenByTaskIdAsync(task.ID)
		if children.Err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, dayRequest.Err.Error())
		}

		response.Tasks = append(response.Tasks, TaskResponse{
			Task: task,
			Children: children.Data.([]model.Task),
		})
	}

	return c.JSON(http.StatusOK, response)
}

type TasksByDayResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

type TaskResponse struct {
	Task model.Task `json:"task"`
	Children []model.Task `json:"children"`
}
