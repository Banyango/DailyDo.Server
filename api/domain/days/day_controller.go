package days

import (
	"fmt"
	"github.com/Banyango/gifoody_server/api/domain/users"
	"github.com/Banyango/gifoody_server/api/infrastructure/pagination"
	"github.com/Banyango/gifoody_server/api/infrastructure/time"
	"github.com/Banyango/gifoody_server/api/infrastructure/utils"
	"github.com/Banyango/gifoody_server/api/model"
	"github.com/Banyango/gifoody_server/api/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
	"net/http"
)

type DayController struct {
	timeService    time.ITimeInterface
	authService    users.IUserAuthService
	dayRepository  repositories.IDayRepository
}

func NewDayController(
	timeService time.ITimeInterface,
	dayRepository repositories.IDayRepository,
	authService users.IUserAuthService) *DayController {
	return &DayController{
		timeService: timeService,
		dayRepository:  dayRepository,
		authService:    authService}
}

// @Summary List Day.
// @Description Get a paginated list of Days desc from latest date.
// @Accept json
// @Produce json
// @Param limit query string false "pagination limit"
// @Param offset query string false "pagination limit"
// @Success 200 {object} model.PagedResult
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/auth/days/ [get]
func (self *DayController) ListDays(c echo.Context) (err error) {
	page := c.Get("Pagination").(pagination.Pagination)

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	dayRequest := <-self.dayRepository.GetDaysAsync(user.Id, page.Limit, page.Offset)
	if dayRequest.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dayRequest.Err.Error())
	}

	pagedResult := pagination.NewPagedResult("Day", c, dayRequest)

	return c.JSON(http.StatusOK, pagedResult)
}

// @Summary Create Day.
// @Description Create a Day.
// @Accept json
// @Produce json
// @Success 201
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/days/ [post]
func (self *DayController) CreateDay(c echo.Context) (err error) {
	request, err := NewCreateDayRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad Request")
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	day := model.Day{
		ID:           uuid.New().String(),
		Summary:      null.NewString(request.Summary, true),
		Date:         self.timeService.GetStartOfDayUTC(),
		UserID:       user.Id,
	}

	result := self.dayRepository.Save(day)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error saving day")
	}

	return c.JSON(http.StatusCreated, result.Data.(model.Day))
}

// @Summary Update Day.
// @Description Update a Day.
// @Accept json
// @Produce json
// @Param request body days.UpdateDayRequestFromContext true "request"// @Success 200
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/day/{id} [put]
func (self *DayController) UpdateDay(c echo.Context) (err error) {

	request, err := NewUpdateDayRequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	dayById := <-self.dayRepository.GetDayByIdAsync(request.ID)
	if dayById.Err != nil {
		return utils.LogError(dayById.Err, http.StatusNotFound, fmt.Sprintf("Day id={%s} not found", request.ID))
	}

	day := dayById.Data.(model.Day)

	if day.Summary.String != request.Summary {
		day.Summary = null.NewString(request.Summary, true)
	}

	result := <-self.dayRepository.UpdateAsync(day)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error deleting day")
	}

	return c.JSON(http.StatusOK, day)
}