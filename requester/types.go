package requester

import "github.com/valyala/fasthttp"

const baseURL = "https://discord.com/api/webhooks/"

var client = &fasthttp.Client{}
