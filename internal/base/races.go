package base

import (
	"math/rand/v2"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Race struct {
	Track OutputTrack `json:"track"`
	Car OutputCar `json:"car"`
}

func RacesGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting all races!!")
}

func GetRaceWithFiltersHandler(c echo.Context) error {
	allParams := c.QueryParams()

	tracks, _ := GetTracksWithFilters(allParams)
	cars, _ := GetCarsWithFilters(allParams)

	randTrackIndex := rand.IntN(len(tracks.Tracks))
	randCarIndex := rand.IntN(len(cars.Cars))

	outputTrack := tracks.Tracks[randTrackIndex]
	outputCar := cars.Cars[randCarIndex]

	race := Race{
		Track: outputTrack,
		Car: outputCar,
	}
	output, _ := json.Marshal(race)

	return c.String(http.StatusOK, string(output))
}

func RacesPutHandler(c echo.Context) error {
	return c.String(http.StatusOK, "We added a race")
}