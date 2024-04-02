package base

import (
	// "encoding/json"

	"github.com/labstack/echo/v4"
	"net/http"
)

type Car struct {
	Name      string `json:name`
	Class     string `json:class`
	Rating    int8   `json:rating`
	Price     int32  `json:price`
	Races     int16  `json:races`
	LastRaced string `json:lastRaced`
	SetupLink string `json:setupLink`
}

func CarsGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Get all the cars!!")
}

func CarsIdGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting 1 car")
}
