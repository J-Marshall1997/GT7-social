package main

import (
	"github.com/gt7social/server"
	"github.com/gt7social/internal/base"
)

func main() {
	_ = base.ReconsileCars()
	server.Server()
}
