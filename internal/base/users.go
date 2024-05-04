package base 

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func UsersGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting all users!!")
}

func UsersIdGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting 1 user!")
}