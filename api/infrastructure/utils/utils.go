package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"runtime/debug"
)

func LogError(err error, code int, message string) error {
	log.Printf("Error: %s", err.Error())
	debug.PrintStack()
	return echo.NewHTTPError(code, message)
}
