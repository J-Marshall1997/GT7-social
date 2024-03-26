package base

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func TracksGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting all tracks!!")
}

func TracksIdGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting 1 track!")
}