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

type RaceConstants struct {
	PPSlope float32 // Leave the type name alone. This is how many m/s a car improves by per PP
	YInt float32 // This is the theoretical m/s of a car with 0PP. Yes, I get that it not being 0 doesn't make sense
}

func RacesGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting all races!!")
}

func GetRaceWithFiltersHandler(c echo.Context) error {
	allParams := c.QueryParams()

	tracks, _ := GetTracksWithFilters(allParams)
	cars, _ := GetCarsWithFilters(allParams)
	raceTime, _ := strconv.Atoi(allParams.Get("raceTime"))

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
	// The logic here needs to be moved and updated to ensure that a track whose length exceeds estimated total race distance isn't selected
	(&race).calculateLaps(raceTime)
	output, _ := json.Marshal(race)

	return c.String(http.StatusOK, string(output))
}

func RacesPutHandler(c echo.Context) error {
	return c.String(http.StatusOK, "We added a race")
}

// This uses a linear equation I've pulled from our race results.
// The y intercept was 17.xx but it wasn't giving great results so I changed it manually to 10
func (race *Race) calculateLaps(raceTime int) {

	lapsConstants := RaceConstants{
		PPSlope: 0.043,
		YInt: 10,
	}
	carMS := (lapsConstants.PPSlope * race.Car.PP) + lapsConstants.YInt
	raceDistance := (raceTime*60) * int(carMS) // raceTime is in minutes
	race.Track.Laps = raceDistance / race.Track.Length
}