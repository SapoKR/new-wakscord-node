package server

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/wakscord/new-wakscord-node/env"
)

func Run() error {
	address := fmt.Sprintf("%s:%d", env.GetString("HOST", "0.0.0.0"), env.GetInt("PORT", 3000))
	go taskHandler()
	log.Printf("listening on %s\n", address)
	return fasthttp.ListenAndServe(address, requestHandler)
}
