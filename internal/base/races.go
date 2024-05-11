package base

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"

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
	raceTime, _ := strconv.Atoi(allParams.Get("raceTime"))

	randTrackIndex := rand.IntN(len(tracks.Tracks))
	randCarIndex := rand.IntN(len(cars.Cars))

	outputTrack := tracks.Tracks[randTrackIndex]
	outputCar := cars.Cars[randCarIndex]

	race := Race{
		Track: outputTrack,
		Car: outputCar,
	}
	(&race).calculateLaps(raceTime)
	output, _ := json.Marshal(race)

	return c.String(http.StatusOK, string(output))
}

func RacesPutHandler(c echo.Context) error {
	return c.String(http.StatusOK, "We added a race")
}

func (race *Race) calculateLaps(raceTime int) {
	carPP := race.Car.PP
	trackLength := race.Track.Length
	timePerLapBase := float32(trackLength)/26.0
	timePerLapCar := timePerLapBase - (carPP/7.5)
	laps := (raceTime*60) / int(timePerLapCar)
	race.Track.Laps = laps 
	fmt.Printf("timePerLapBase: %f\ntimePerLapCar: %f\n", timePerLapBase, timePerLapCar)
	fmt.Printf("raceTime: %d\n", raceTime)
	fmt.Printf("laps: %d\n", laps)
}