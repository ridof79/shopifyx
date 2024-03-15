package util

import (
	"github.com/labstack/echo/v4"
)

const ()

func ErrorHandler(c echo.Context, code int, message string) error {
	return c.JSON(code,
		map[string]string{
			"error": message,
		},
	)
}

func ResponseHandler(c echo.Context, code int, message string) error {
	return c.JSON(code,
		map[string]string{
			"message": message,
		},
	)
}
