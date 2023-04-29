package discord

import (
	"github.com/valyala/fasthttp"
	"github.com/wakscord/new-wakscord-node/env"
)

func Initialize() {
	maxConns := env.GetInt("MAX_CONCURRENT", 500)
	fasthttpClient = &fasthttp.Client{
		MaxConnsPerHost: maxConns,
	}
}
