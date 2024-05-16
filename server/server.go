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
	carsGroup.GET("/filter", base.GetCarsWithFiltersHandler)

	tracksGroup := e.Group("tracks")
	tracksGroup.GET("/filter", base.GetTracksWithFiltersHandler)

	racesGroup := e.Group("races")
	racesGroup.GET("/filter", base.GetRaceWithFiltersHandler)
	racesGroup.PUT("", base.RacesPutHandler)

	usersGroup := e.Group("users")
	usersGroup.GET("", base.UsersGetHandler)
	usersGroup.GET("/:name", base.UsersIdGetHandler)

	e.GET("/", defaultHandler)
	e.Logger.Fatal(e.Start(":1323"))
}