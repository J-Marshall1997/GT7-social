package base

import (
	// "encoding/json"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"log"

	"github.com/gt7social/config"
	"github.com/labstack/echo/v4"
)

type CarRequestFilters struct {
	Groups []string
	MinPP  float32
	MaxPP  float32
	MinPrice int
	MaxPrice int
}

type APIData struct {
	Cars Cars `json:"data"`
}

type Cars struct {
	Cars []Car `json:"gt_car"`
}

type Car struct {
	Id            int          `json:"id"`
	Manufacturer  Manufacturer `json:"manufacturer"`
	Name          string       `json:"name"`
	ShortName     string       `json:"short_name"`
	Slug          string       `json:"slug"`
	Group         string       `json:"group"`
	Price         int          `json:"price"`
	UsedPrice     int          `json:"used_price"`
	HaggertyPrice int          `json:"haggerty_price"`
	PP            PP           `json:"pp"`
	Tags          []string     `json:"tags"`
	MaxPower      string       `json:"max_power"`
	Update        interface{}  `json:"update"`
}

type OutputCars struct {
	Cars []OutputCar `json:"cars"`
}

type OutputCar struct {
	Manufacturer string  `json:"manufacturer"`
	ShortName 	 string  `json:"short_name"`
	Group		 string  `json:"group"`
	Price		 int 	 `json:"price"`
	Shop		 string  `json:"shop"`
	PP			 float32 `json:"pp"`
}

type Manufacturer struct {
	Name string `json:"name"`
}

type PP float32

func (pp *PP) UnmarshalJSON(d []byte) error {
	var v float32
	err := json.Unmarshal(bytes.Trim(d, `"`), &v)
	*pp = PP(v)
	return err
}

func GetCarsWithFiltersHandler(c echo.Context) error {
	log.Println("Received call to /cars")
	cars, _ := GetCarsWithFilters(c.QueryParams())
	
	output, err := json.Marshal(cars)
	if err != nil {
		fmt.Println("Failed to marshal car output")
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, "Failed to get filtered cars")
	}
	log.Println("Returning message from /cars")
	return c.String(http.StatusOK, string(output))
}

func GetCarsWithFilters(params url.Values) (OutputCars, error) {	
	var minPP, maxPP float64
	var minPrice, maxPrice int
	var groups []string
	// Setup filter variables
	if params.Get("minPP") == "" {
		minPP = config.DEFAULT_MIN_PP
	} else {
		minPP, _ = strconv.ParseFloat(params.Get("minPP"), 32)
	}
	if params.Get("maxPP") == "" {
		maxPP = config.DEFAULT_MAX_PP
	} else {
		maxPP, _ = strconv.ParseFloat(params.Get("maxPP"), 32)
	}
	if params["groups"] == nil {
		groups = config.DEFAULT_GROUPS
	} else {
		groups = params["groups"]
	}
	if params.Get("minPrice") == "" {
		minPrice = config.DEFAULT_MIN_PRICE
	} else {
		minPrice, _ = strconv.Atoi(params.Get("minPrice"))
	}	
	if params.Get("maxPrice") == "" {
		maxPrice = config.DEFAULT_MAX_PRICE
	} else {
		maxPrice, _ = strconv.Atoi(params.Get("maxPrice"))
	}
	filters := CarRequestFilters{
		Groups: groups,
		MinPP: float32(minPP),
		MaxPP: float32(maxPP),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}
	
	// Get all of the cars from storage
	filePath := "internal/storage/cars.json"
	var Data APIData
	body, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file into byte string")
		fmt.Println(err.Error())
		return OutputCars{}, err
	}
	err = json.Unmarshal(body, &Data)
	if err != nil {
		fmt.Println("Failed to unmarshal car data")
		fmt.Printf("%v\n", string(body))
		return OutputCars{}, err
	}

	filteredCars := filterPP(Data.Cars, filters.MinPP, filters.MaxPP)
	filteredCars = filterGroup(filteredCars, filters.Groups)
	filteredCars = filterPrice(filteredCars, filters.MinPrice, filters.MaxPrice)
	outputCars := prepareOutputCars(filteredCars)
	return outputCars, nil
}

func ReconsileCars() error {
	// Need to hit the gt7 API endpoint and build a map where the id of the car is the key
	fmt.Println("Reconsiling Cars")
	resp, err := http.Get("https://gtdb.io/api/graphql_middleware/query/AllCars")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Have hit the API")
	body, _ := io.ReadAll(resp.Body)

	filePath := "internal/storage/cars.json"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println("Error opening cars.json")
		fmt.Println(err.Error())
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		fmt.Println("Error writing to file")
	}

	return nil
}

func filterPP(c Cars, min, max float32) (out Cars) {
	out = Cars{}
	for _, car := range c.Cars {
		if float32(car.PP) >= min && float32(car.PP) <= max {
			out.Cars = append(out.Cars, car)
		}
	}
	return out
}

func filterGroup(c Cars, groups []string) (out Cars) {
	out = Cars{}
	for _, car := range c.Cars {
		for _, group := range groups {
			if car.Group == group || (car.Group == "" && group == "open"){
				out.Cars = append(out.Cars, car)
			}
		}
	}
	return out
}

func filterPrice(c Cars, minPrice, maxPrice int) (out Cars) {
	out = Cars{}
	for _, car := range c.Cars {
		// Get the lowest price from all of the possible scores
		lowestPrice, _ := getLowestPrice(car)
		if lowestPrice >= minPrice && lowestPrice <= maxPrice {
			out.Cars = append(out.Cars, car)
		}
	}
	return out
}

func getLowestPrice(c Car) (lowestPrice int, store string) {
	lowestPrice = config.DEFAULT_MAX_PRICE
	store = "not on sale"

	if c.Price != 0 && c.Price < lowestPrice {
		lowestPrice = c.Price
		store = "brand central"
	}
	if c.HaggertyPrice != 0 && c.HaggertyPrice < lowestPrice {
		lowestPrice = c.HaggertyPrice
		store = "legend cars"
	}
	if c.UsedPrice != 0 && c.UsedPrice < lowestPrice {
		lowestPrice = c.UsedPrice
		store = "used cars"
	}
	return lowestPrice, store
}

func prepareOutputCars(cars Cars) (out OutputCars) {
	out = OutputCars{}

	for _, car := range cars.Cars {
		lowestPrice, store := getLowestPrice(car)
		outputCar := OutputCar{
			Manufacturer: car.Manufacturer.Name,
			ShortName: car.ShortName,
			Group: car.Group,
			Price: lowestPrice,
			Shop: store,
			PP: float32(car.PP),
		}
		out.Cars = append(out.Cars, outputCar)
	}
	return out
}