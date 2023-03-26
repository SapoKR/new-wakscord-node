package main

import (
	"log"

	"github.com/wakscord/new-wakscord-node/server"
)

func main() {
	log.Fatal(server.Run())
}
