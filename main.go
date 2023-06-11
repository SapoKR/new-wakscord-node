package main

import (
	"log"

	"github.com/wakscord/node/config"
	"github.com/wakscord/node/discord"
	"github.com/wakscord/node/server"
)

func main() {
	err := config.Initialize()
	if err != nil {
		panic(err)
	}
	discord.Initialize()

	log.Fatal(server.Run())
}
