package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func LogError(err error, code int, message string) error {
	log.Printf("Error: %s,", err.Error())
	return echo.NewHTTPError(code, message)
}
