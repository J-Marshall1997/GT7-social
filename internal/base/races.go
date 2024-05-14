package base

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/gt7social/internal/utils"
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

	numCars, _ := strconv.Atoi(allParams.Get("numCars"))
	numCars = utils.MinInt(numCars, 10, len(cars.Cars))
	var indexes map[int]bool

	for i := 0; i < numCars; i++ {
		randCarIndex := rand.IntN(numCars)
		if !indexes[randCarIndex] {
			indexes[randCarIndex] = true
		}
	}
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