package index

import (
	"github.com/Banyango/dailydo_server/api/infrastructure/links"
	"github.com/Banyango/dailydo_server/api/model"
	"github.com/labstack/echo/v4"
)

type IndexResponse struct {
	IndexLinks IndexLinks `json:"_links"`
}

type IndexLinks struct {
	Posts        model.Link `json:"tasks"`
	Logout       model.Link `json:"logout"`
	Login        model.Link `json:"login"`
	Forgot       model.Link `json:"forgotPassword"`
	Register     model.Link `json:"register"`
	Confirm      model.Link `json:"confirm"`
	ConfirmReset model.Link `json:"confirmResetPassword"`
	Me           model.Link `json:"me"`
}

func NewIndexResponse(ctx echo.Context) (IndexResponse, error) {
	postLink, err := links.NewLinkBuilder("auth/tasks", "get", nil, ctx)
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

	meLink, err := links.NewLinkBuilder("auth/me", "get", nil, ctx)
	if err != nil {
		return IndexResponse{}, err
	}

	return IndexResponse{
		IndexLinks: IndexLinks{
			Posts:        postLink.BuildLink(),
			Logout:       logoutLink.BuildLink(),
			Login:        loginLink.BuildLink(),
			Forgot:       forgotLink.BuildLink(),
			Register:     registerLink.BuildLink(),
			Confirm:      confirmLink.BuildLink(),
			ConfirmReset: confirmResetLink.BuildLink(),
			Me:           meLink.BuildLink(),
		},
	}, nil
}
