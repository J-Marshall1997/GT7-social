package base

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func CarsGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Get all the cars!!")
}

func CarsIdGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting 1 car")
}