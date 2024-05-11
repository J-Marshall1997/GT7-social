package base

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gt7social/config"
	"github.com/labstack/echo/v4"
)

type TrackRequestFilters struct {
	Type string
	MinLength int
	MaxLength int
	Rain string
}

type Tracks struct {
	Tracks []Track `json:"tracks"`
}

type Track struct {
	Course string `json:"course"`
	Layout string `json:"layout"`
	Type string `json:"type"`
	Length int `json:"length"`
	Rain string `json:"rain"`
	Rating int `json:"rating"`
	Races int `json:"races"`
	LastRaced string `json:"last_raced"`
}

type OutputTracks struct {
	Tracks []OutputTrack
}

type OutputTrack struct {
	Course string `json:"course"`
	Layout string `json:"layout"`
	Length int
	Laps int `json:"laps"`
}

func TracksGetHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Getting all tracks!!")
}

func GetTracksWithFiltersHandler(c echo.Context) error {
	tracks, _ := GetTracksWithFilters(c.QueryParams())

	output, err := json.Marshal(tracks)
	if err != nil {
		fmt.Println("Failed to marshal track output")
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, "Failed to get filtered tracks")
	}
	return c.String(http.StatusOK, string(output))
}

func GetTracksWithFilters(params url.Values) (OutputTracks, error) {
	var minLen, maxLen int
	var trackType, rain string
	if params.Get("minLength") == "" {
		minLen = config.DEFAULT_MIN_LENGTH 
	} else {
		minLen, _ = strconv.Atoi(params.Get("minLength"))
	}
	if params.Get("maxLength") == "" {
		maxLen = config.DEFAULT_MAX_LENGTH 
	} else {
		maxLen, _ = strconv.Atoi(params.Get("maxLength"))
	}
	if params.Get("type") == "" {
		trackType = config.DEFAULT_TYPE 
	} else {
		trackType = params.Get("type") 
	}
	if params.Get("rain") == "" {
		rain = config.DEFAULT_RAIN
	} else {
		rain = params.Get("rain")
	}
	filters := TrackRequestFilters{
		Type: trackType,
		MinLength: minLen,
		MaxLength: maxLen,
		Rain: rain,
	}

	// Get all tracks from storage
	filePath := "internal/storage/tracks.json"
	var Tracks Tracks
	body, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file into byte string")
		return OutputTracks{}, err
	}
	err = json.Unmarshal(body, &Tracks)
	if err != nil {
		fmt.Println("Failed to unmarshal track data")
		fmt.Printf("%v\n", string(body))
		return OutputTracks{}, err
	}

	filteredTracks := filterType(Tracks, filters.Type)
	filteredTracks = filterLength(filteredTracks, filters.MinLength, filters.MaxLength)
	filteredTracks = filterRain(filteredTracks, filters.Rain)
	outputTracks := prepareOutputTracks(filteredTracks)
	return outputTracks, nil
}

func filterType(tracks Tracks, trackType string) (out Tracks) {
	out = Tracks{}
	for _, track := range tracks.Tracks {
		if track.Type == trackType {
			out.Tracks = append(out.Tracks, track)
		}
	}
	return out
}

func filterLength(tracks Tracks, min, max int) (out Tracks) {
	out = Tracks{}
	for _, track := range tracks.Tracks {
		if track.Length >= min && track.Length <= max {
			out.Tracks = append(out.Tracks, track)
		}
	}
	return out
}

func filterRain(tracks Tracks, rain string) (out Tracks) {
	out = Tracks{}
	if rain == "Both" {
		return tracks
	}
	for _, track := range tracks.Tracks {
		if track.Rain == rain {
			out.Tracks = append(out.Tracks, track)
		}
	}
	return out
}

func prepareOutputTracks(tracks Tracks) (out OutputTracks) {
	out = OutputTracks{}

	for _, track := range tracks.Tracks {
		outputTrack := OutputTrack{
			Course: track.Course,
			Layout: track.Layout,
			Length: track.Length,
			Laps: 0,
		}
		out.Tracks = append(out.Tracks, outputTrack)
	}
	return out
}