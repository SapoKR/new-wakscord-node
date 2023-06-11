package server

import (
	"errors"
	"time"

	"github.com/wakscord/node/discord"
)

var (
	status nodeStatus

	deletedWebhooks = map[string]struct{}{}
	tasks           chan task

	startTime = time.Now()

	serverKey string

	errQueueIsFull = errors.New("queue is full")
)

type nodeStatus struct {
	Info       nodeInfo    `json:"info"`
	Pending    nodePending `json:"pending"`
	Processed  int32       `json:"processed"`
	Deleted    int         `json:"deleted"`
	Uptime     int         `json:"uptime"`
	Goroutines int         `json:"goroutines"`
}

type nodeInfo struct {
	NodeID int    `json:"node_id"`
	Owner  string `json:"owner"`
}

type nodePending struct {
	Total    int32 `json:"total"`
	Messages int32 `json:"messages"`
	Tasks    int32 `json:"tasks"`
}

type requestPayload struct {
	Keys []string              `json:"keys"`
	Data discord.WebhookParams `json:"data"`
}

type task struct {
	chunks [][]string
	data   discord.WebhookParams
}
