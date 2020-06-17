package services

import (
	"fmt"
	"github.com/Banyango/gifoody_server/constants"
	"github.com/Banyango/gifoody_server/model"
	"github.com/labstack/echo/v4"
	url2 "net/url"
	"strings"
)

type LinkBuilder struct {
	Host   string
	Scheme string
	Path   string
	Method string
	Param  string
}

func NewLinkBuilder(path string, method string, params []string, ctx echo.Context) (LinkBuilder, error) {
	host, scheme := ParseRequest(ctx)
	err := ValidateHost(host)
	if err != nil {
		return LinkBuilder{}, err
	}

	builder := LinkBuilder{
		Host:   host,
		Scheme: scheme,
		Path:   path,
		Method: method,
		Param:  strings.Join(params, "&"),
	}

	return builder, nil
}

func ParseRequest(ctx echo.Context) (host string, scheme string) {
	return ctx.Request().Host, ctx.Scheme()
}

func ValidateHost(host string) error {
	if !strings.Contains(host, "api.gifoody.com") && !strings.Contains(host, "localhost") {
		return fmt.Errorf("Bad Request")
	}
	return nil
}

func (l *LinkBuilder) BuildLink() model.Link {
	url := url2.URL{
		Scheme:   l.Scheme,
		Host:     l.Host,
		Path:     strings.Join([]string{constants.API_PATH, l.Path}, ""),
		RawQuery: l.Param,
	}

	return model.Link{
		Href:   url.String(),
		Method: l.Method,
	}
}
