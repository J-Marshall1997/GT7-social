package base

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RacesGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting all races!!")
}

func RacesIdGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting 1 race!")
}

func RacesPutHandler(c echo.Context) error {
	return c.String(http.StatusOK, "We added a race")
}