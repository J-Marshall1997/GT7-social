package base

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/gt7social/internal/utils"
	"github.com/labstack/echo/v4"
)

var lapsConstants = RaceConstants{
						PPSlope: 0.043,
						YInt: 10,
					}
type Race struct {
	Track OutputTrack `json:"track"`
	Cars []OutputCar `json:"car"`
}

type RaceConstants struct {
	PPSlope float32 // Leave the name alone. This is how many m/s a car improves by per PP
	YInt float32 // This is the theoretical m/s of a car with 0PP. Yes, I get that it not being 0 doesn't make sense
}

func GetRaceWithFiltersHandler(c echo.Context) error {
	allParams := c.QueryParams()
	raceTime, _ := strconv.Atoi(allParams.Get("raceTime"))
	maxPP, _ := strconv.ParseFloat(c.QueryParams().Get("maxPP"), 32) // Still have to float32() this anyway...

	tracks, _ := GetTracksWithFilters(allParams)
	cars, _ := GetCarsWithFilters(allParams)

	// Select the amount of cars to be in the race
	numCarsInt, _ := strconv.Atoi(allParams.Get("numCars"))
	numCars := utils.MaxInt(numCarsInt, 1) 
	numCars = utils.MinInt(numCars, 10, len(cars.Cars))
	
	var race Race
	if len(tracks.Tracks) >= numCars {
		race = initRace(cars.Cars, tracks.Tracks, numCars, raceTime, float32(maxPP)) // Using maxPP to allow for car upgrades
	} else {
		race = initRace(cars.Cars, tracks.Tracks, len(tracks.Tracks), raceTime, float32(maxPP))
	}
	
	output, _ := json.Marshal(race)

	return c.String(http.StatusOK, string(output))
}

func RacesPutHandler(c echo.Context) error {
	return c.String(http.StatusOK, "We added a race")
}

// This uses a linear equation I've pulled from our race results.
// The y intercept was 17.xx but it wasn't giving great results so I changed it manually to 10
func (race *Race) calculateLaps(raceTime int, pp float32) {

	carMS := (lapsConstants.PPSlope * pp) + lapsConstants.YInt
	raceDistance := (raceTime*60) * int(carMS) // raceTime is in minutes
	race.Track.Laps = raceDistance / race.Track.Length
}

// Taking in all possible cars, tracks and required parameters for lap calculation, create a race 
func initRace(cars []OutputCar, tracks []OutputTrack, numCars, time int, maxPP float32) (r Race) {
	
	// Need to get car first as the car's PP is needed to filter out tracks that can't be finished in time
	indexes := make(map[int]bool, numCars)
	for i := 0; i < numCars; i++ {
		index := getUniqueIndex(indexes, numCars)
		indexes[index] = true
	}
	for index := range indexes {
		r.Cars = append(r.Cars, cars[index])
	}

	// Now that we have our cars, we need to make sure that our cars can get around the track in time
	maxLength := int(float32(time*60) * (lapsConstants.PPSlope * maxPP) + lapsConstants.YInt)
	for {
		if len(tracks) == 0 {
			r.Track = OutputTrack{
				Course: "No course can be finished in time with these parameters",
			}
			return r
		}
		trackIndex := rand.IntN(len(tracks))
		if tracks[trackIndex].Length <= maxLength {
			r.Track = tracks[trackIndex]
			r.calculateLaps(time, maxPP) // May as well do this straight away
			break
		} else {
			tracks = append(tracks[:trackIndex], tracks[trackIndex+1:]...)
		}
	}

	return r
}

// Quick and dirty method to get a previously unselected index
func getUniqueIndex(i map[int]bool, max int) (index int ){
	for {
		index = rand.IntN(max)
		if !i[index] {
			return index
		}
	}
}