package server

import (
	"github.com/gt7social/internal/base"

	"net/http"
	"github.com/labstack/echo/v4"
)

func defaultHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

func Server() {
	e := echo.New()

	carsGroup := e.Group("/cars")
	carsGroup.GET("", base.CarsGetHandler)
	carsGroup.GET("/filter", base.GetCarsWithFilters)

	tracksGroup := e.Group("tracks")
	tracksGroup.GET("", base.TracksGetHandler)
	tracksGroup.GET("/:name", base.TracksIdGetHandler)

	racesGroup := e.Group("races")
	racesGroup.GET("", base.RacesGetHandler)
	racesGroup.GET("/:name", base.RacesIdGetHandler)
	racesGroup.PUT("", base.RacesPutHandler)

	usersGroup := e.Group("users")
	usersGroup.GET("", base.UsersGetHandler)
	usersGroup.GET("/:name", base.UsersIdGetHandler)

	e.GET("/", defaultHandler)
	e.Logger.Fatal(e.Start(":1323"))
}