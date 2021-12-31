package pagination

import (
	"fmt"
	"github.com/Banyango/dailydo_server/api/infrastructure/links"
	"github.com/Banyango/dailydo_server/api/model"
	"github.com/Banyango/dailydo_server/api/repositories/util"
	"github.com/labstack/echo/v4"
)

type PagedResult struct {
	Offset int              `json:"offset"`
	Limit  int              `json:"limit"`
	Items  interface{}      `json:"items"`
	Total  int              `json:"total"`
	Links  PagedResultLinks `json:"_links"`
}

type PagedResultLinks struct {
	Next     model.Link `json:"next"`
	Previous model.Link `json:"previous"`
}

func NewPagedResult(
	url string,
	ctx echo.Context,
	result util.StoreResult) PagedResult {

	page := ctx.Get("Pagination").(Pagination)

	pagedResult := PagedResult{
		Offset: page.Offset,
		Limit:  page.Limit,
		Items:  result.Data,
		Total:  result.Total,
		Links:  PagedResultLinks{},
	}

	if (page.Offset + page.Limit) < result.Total {
		offset := fmt.Sprintf("offset=%d", page.Offset+page.Limit)
		limit := fmt.Sprintf("limit=%d", page.Limit)
		next, _ := links.NewLinkBuilder(url, "get", []string{offset, limit}, ctx)
		pagedResult.Links.Next = next.BuildLink()
	}

	if page.Offset > 1 {
		offset := fmt.Sprintf("offset=%d", page.Offset-page.Limit)
		limit := fmt.Sprintf("limit=%d", page.Limit)
		previous, _ := links.NewLinkBuilder(url, "get", []string{offset, limit}, ctx)
		pagedResult.Links.Previous = previous.BuildLink()
	}

	return pagedResult
}
