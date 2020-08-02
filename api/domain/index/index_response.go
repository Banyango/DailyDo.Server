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
	Categories   model.Link `json:"categories"`
	Posts        model.Link `json:"posts"`
	Logout       model.Link `json:"logout"`
	Login        model.Link `json:"login"`
	Forgot       model.Link `json:"forgotPassword"`
	Register     model.Link `json:"register"`
	Confirm      model.Link `json:"confirm"`
	ConfirmReset model.Link `json:"confirmResetPassword"`
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

	logoutLink, err := links.NewLinkBuilder("logout", "post", nil, ctx)
	if err != nil {
		return IndexResponse{}, err
	}

	loginLink, err := links.NewLinkBuilder("login", "post", nil, ctx)
	if err != nil {
		return IndexResponse{}, err
	}

	forgotLink, err := links.NewLinkBuilder("reset_password", "post", nil, ctx)
	if err != nil {
		return IndexResponse{}, err
	}

	registerLink, err := links.NewLinkBuilder("register", "post", nil, ctx)
	if err != nil {
		return IndexResponse{}, err
	}

	confirmLink, err := links.NewLinkBuilder("confirm_account", "post", nil, ctx)
	if err != nil {
		return IndexResponse{}, err
	}

	confirmResetLink, err := links.NewLinkBuilder("confirm_reset_password", "post", nil, ctx)
	if err != nil {
		return IndexResponse{}, err
	}

	return IndexResponse{
		IndexLinks: IndexLinks{
			Categories:   categoriesLink.BuildLink(),
			Posts:        postLink.BuildLink(),
			Logout:       logoutLink.BuildLink(),
			Login:        loginLink.BuildLink(),
			Forgot:       forgotLink.BuildLink(),
			Register:     registerLink.BuildLink(),
			Confirm:      confirmLink.BuildLink(),
			ConfirmReset: confirmResetLink.BuildLink(),
		},
	}, nil
}
