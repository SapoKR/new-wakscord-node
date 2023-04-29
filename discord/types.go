package discord

import (
	"github.com/valyala/fasthttp"
)

const baseURL = "https://discord.com/api/webhooks/"

var fasthttpClient *fasthttp.Client

type Response struct {
	Key   string
	Code  int
	Error error
	Body  string
}
