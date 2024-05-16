package main

import (
	"fmt"

	"github.com/gt7social/internal/base"
	"github.com/gt7social/server"
)

func main() {
	if err := base.ReconsileCars(); err != nil {
		fmt.Printf("error while reconsiling cars: %v", err.Error())
	}
	server.Server()
}
