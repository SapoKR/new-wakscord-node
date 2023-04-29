package main

import (
	"log"

	"github.com/wakscord/new-wakscord-node/discord"
	"github.com/wakscord/new-wakscord-node/server"
)

func main() {
	discord.Initialize()

	log.Fatal(server.Run())
}
