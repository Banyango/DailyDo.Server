package controllers

import (
	"github.com/Banyango/gifoody_server/api/services"
	"github.com/Banyango/gifoody_server/model"
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
	categoriesLink, err := services.NewLinkBuilder("categories", "get", nil, ctx)
	if err != nil {
		return IndexResponse{}, err
	}

	postLink, err := services.NewLinkBuilder("posts", "get", nil, ctx)
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
