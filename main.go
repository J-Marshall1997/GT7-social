package main

import (
	"github.com/gt7social/internal/scraper"
	// "github.com/gt7social/server"
)

func main() {
	scraper.Scrape("https://gtdb.io/gt7/all-cars/")
	// server.Server()
}