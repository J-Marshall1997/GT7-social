package scraper

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Not all cars have a price listed. Some cars have multiple prices
// This function will handle logic to ensure that
// the database has a useful output for the car price
func handlePrice(raw_price string) string {

	// Split prices that have multiple locations
	split_pattern := "(.*?)(0)([a-zA-Z])"
	regex := regexp.MustCompile(split_pattern)

	matches := regex.FindAllStringSubmatch(raw_price, -1)

	// Extract matched parts
	var prices []string
	for _, match := range matches {
		prices = append(prices, match[1]+match[2])
	}

	// If there is a leftover unmatched portion, add it as well
	if len(matches) > 0 && len(matches[len(matches)-1]) > 3 {
		prices = append(prices, matches[len(matches)-1][2]+matches[len(matches)-1][3])
	}

	// prices := regex.FindAllStringSubmatch(raw_price, -1)
	fmt.Println(prices)
	return "0"
}

func Scrape(url string) {

	fName := "./internal/storage/cars.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatal("Failed to open cars csv")
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	c := colly.NewCollector(
		// colly.AllowedDomains(
		// 	"https://gtdb.io/",
		// 	"gtdb.io/",
		// 	"https://gtdb.io",
		// ),
		// Set a cache to stop multiple download of pages
		colly.CacheDir("./cache/car_cache"),
	)

	carCollector := c.Clone()

	// Write header row
	writer.Write([]string{"Name", "Price", "PP"})

	carCollector.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("tbody tr", func(_ int, el *colly.HTMLElement) {

			// Set the lowest price to be used for this car
			lowestPrice := 0
			store := ""

			el.ForEach("div", func(_ int, elem *colly.HTMLElement) {
				switch class := elem.Attr("class"); class {
				case "font-light text-xs":
					store = elem.Text
				case "text-right grow tabular-nums text-gray-100":
					// Very annoying string manipulation to prepare element for price comparison
					priceAsStr := strings.Replace(strings.Split(elem.Text, " ")[2], ",", "", -1)
					priceAsInt, _ := strconv.Atoi(priceAsStr)
					if priceAsInt < lowestPrice || lowestPrice == 0 {
						lowestPrice = priceAsInt
					}
				default:
				}
			})

			writer.Write([]string{
				strings.Split(el.ChildText("td:nth-of-type(1)"), "Cr.")[0],
				fmt.Sprintf("%s (%s)", fmt.Sprint(lowestPrice), store),
				el.ChildText("td:nth-of-type(3)"),
			})
		})
	})

	fmt.Printf("Visiting %s\n", url)
	carCollector.Visit("https://gtdb.io/gt7/all-cars/")
}
