package generate

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"text/template"
)

//go:generate -in=$GOFILE -out=gen-$GOFILE

func main() {
	typeInput := "hi"
	routeInput := "/api/v1/"

	f, err := os.Create(fmt.Sprintf("%s.go", typeInput))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	controllerTemplate.Execute(f, struct {
		Type 	string
		Route   string
	}{
		Type: typeInput,
		Route: routeInput,
	})
}

var controllerTemplate = template.Must(template.New("").Parse(`
// This file was generated

type {{ .Type }}Controller struct {
	authService    users.IUserAuthService
}

func New{{ .Type }}Controller(
	authService users.IUserAuthService) *{{ .Type }}Controller {
	return &{{ .Type }}Controller{authService:authService}
}

// @Summary List {{ .Type }}s.
// @Description Get a paginated list of {{ .Type }}.
// @Accept json
// @Produce json
// @Param limit query string false "pagination limit"
// @Param offset query string false "pagination limit"
// @Success 200 {object} model.PagedResult
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router {{ .Route }} [get]
func (self *{{ .Type }}Controller) List{{ .Type }}(c echo.Context) (err error) {
	page := c.Get("Pagination").(pagination.Pagination)

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	{{ .Type | strings.ToLower }}Request := <-self.{{ .Type | strings.ToLower }}Repository.Get{{ .Type }}Async(user.Id, page.Limit, page.Offset)
	if {{ .Type | strings.ToLower }}Request.Err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, {{ .Type | strings.ToLower }Request.Err.Error())
	}

	pagedResult := pagination.NewPagedResult("{{ .Type }}", c, {{ .Type | strings.ToLower }}Request)

	return c.JSON(http.StatusOK, pagedResult)
}

// @Summary Create {{ .Type }.
// @Description Create a {{ .Type }.
// @Accept json
// @Produce json
// @Success 201
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/{{ .Type | strings.ToLower }}s/ [post]
func (self *{{ .Type }Controller) Create{{ .Type }(c echo.Context) (err error) {
	request, err := NewCreate{{ .Type }RequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad Request")
	}

	user, err := self.authService.GetLoggedInUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	{{ .Type | strings.ToLower } := model.{{ .Type }{
		ID:           uuid.New().String(),
		Summary:      null.NewString(request.Summary, true),
		Date:         self.timeService.GetStartOfDayUTC(),
		UserID:       user.Id,
	}

	result := self.{{ .Type | strings.ToLower }Repository.Save({{ .Type | strings.ToLower })
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error saving {{ .Type | strings.ToLower }")
	}

	return c.JSON(http.StatusCreated, result.Data.(model.{{ .Type }))
}

// @Summary Update {{ .Type }.
// @Description Update a {{ .Type }.
// @Accept json
// @Produce json
// @Param request body days.Update{{ .Type }RequestFromContext true "request"// @Success 200
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/day/{id} [put]
func (self *{{ .Type }Controller) Update{{ .Type }(c echo.Context) (err error) {

	request, err := NewUpdate{{ .Type }RequestFromContext(c)
	if err != nil {
		return utils.LogError(err, http.StatusBadRequest, "Bad request")
	}

	{{ .Type | strings.ToLower }ById := <-self.{{ .Type | strings.ToLower }Repository.Get{{ .Type }ByIdAsync(request.ID)
	if {{ .Type | strings.ToLower }ById.Err != nil {
		return utils.LogError({{ .Type | strings.ToLower }ById.Err, http.StatusNotFound, fmt.Sprintf("{{ .Type } id={%s} not found", request.ID))
	}

	{{ .Type | strings.ToLower } := {{ .Type | strings.ToLower }ById.Data.(model.{{ .Type })

	if {{ .Type | strings.ToLower }.Summary.String != request.Summary {
		{{ .Type | strings.ToLower }.Summary = null.NewString(request.Summary, true)
	}

	result := <-self.{{ .Type | strings.ToLower }Repository.UpdateAsync({{ .Type | strings.ToLower })
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error deleting {{ .Type | strings.ToLower }")
	}

	return c.JSON(http.StatusOK, {{ .Type | strings.ToLower })
}

// @Summary Delete {{ .Type | strings.ToLower }.
// @Description Delete a {{ .Type | strings.ToLower }.
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {string} string "bad parameters"
// @Failure 500 {string} string "Database error"
// @Router /v1/Parent/{id} [delete]
func (self *{{ .Type }Controller) Delete{{ .Type }(c echo.Context) (err error) {
	id := c.Param("id")

	{{ .Type | strings.ToLower }ById := <-self.{{ .Type | strings.ToLower }Repository.Get{{ .Type }ByIdAsync(id)
	if {{ .Type | strings.ToLower }ById.Err != nil {
		return utils.LogError({{ .Type | strings.ToLower }ById.Err, http.StatusNotFound, "Parent not found")
	}

	result := self.{{ .Type | strings.ToLower }Repository.Delete(id)
	if result.Err != nil {
		return utils.LogError(result.Err, http.StatusInternalServerError, "Error deleting {{ .Type }")
	}

	return c.NoContent(http.StatusOK)
}

`))