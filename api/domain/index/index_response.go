package index

import (
	"github.com/Banyango/gifoody_server/api/infrastructure/links"
	"github.com/Banyango/gifoody_server/api/model"
	"github.com/labstack/echo/v4"
)

type IndexResponse struct {
	IndexLinks IndexLinks `json:"_links"`
}

type IndexLinks struct {
	Categories model.Link `json:"categories"`
	Posts      model.Link `json:"posts"`
}

func NewIndexResponse(ctx echo.Context) (IndexResponse, error) {
	categoriesLink, err := links.NewLinkBuilder("categories", "get", nil, ctx)
	if err != nil {
		return IndexResponse{}, err
	}

	postLink, err := links.NewLinkBuilder("posts", "get", nil, ctx)
	if err != nil {
		return IndexResponse{}, err
	}

	return IndexResponse{
		IndexLinks: IndexLinks{
			Categories: categoriesLink.BuildLink(),
			Posts:      postLink.BuildLink(),
		},
	}, nil
}
