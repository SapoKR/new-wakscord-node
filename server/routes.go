package server

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/wakscord/new-wakscord-node/config"
	"golang.org/x/exp/maps"
)

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	if path == "/" {
		handleIndex(ctx)
		return
	}

	if path == "/request" || path == "/deletedWebhooks" || path == "/environment" {
		requestKey := string(ctx.Request.Header.Peek("Authorization"))
		if requestKey == serverKey {
			switch path {
			case "/request":
				handleRequest(ctx)
			case "/deletedWebhooks":
				handleDeletedWebhooks(ctx)
			case "/environment":
				handleEnvironment(ctx)
			}
		} else {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			fmt.Fprint(ctx, "wrong key")
		}
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
	fmt.Fprint(ctx, "not found")
}

func handleIndex(ctx *fasthttp.RequestCtx) {
	status.Deleted = len(deletedWebhooks)
	status.Uptime = int(time.Since(startTime).Seconds())
	status.Goroutines = runtime.NumGoroutine()

	body, err := json.Marshal(status)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "json marshalling error: %v", err)
	}

	fmt.Fprint(ctx, string(body))
}

func handleRequest(ctx *fasthttp.RequestCtx) {
	if len(tasks) >= config.Default.MessageQueueSize {
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		fmt.Fprint(ctx, "too many tasks")
		return
	}

	payload := new(requestPayload)
	if err := json.Unmarshal(ctx.Request.Body(), payload); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, "body parsing error: %v", err)
	}

	addTask(payload.Keys, payload.Data)

	fmt.Fprint(ctx, "ok")
}

func handleDeletedWebhooks(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case fasthttp.MethodGet:
		body, err := json.Marshal(maps.Keys(deletedWebhooks))
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			fmt.Fprintf(ctx, "json marshalling error: %v", err)
		}
		fmt.Fprint(ctx, string(body))
	case fasthttp.MethodDelete:
		deletedWebhooks = map[string]struct{}{}
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}
}

func handleEnvironment(ctx *fasthttp.RequestCtx) {
	body, err := json.Marshal(config.Default)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		fmt.Fprintf(ctx, "json marshalling error: %v", err)
	}

	fmt.Fprint(ctx, string(body))
}
