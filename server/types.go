package server

import (
	"time"

	"github.com/wakscord/new-wakscord-node/env"
)

var (
	status = nodeStatus{
		Info: nodeInfo{
			NodeID: env.GetInt("ID", 0),
			Owner:  env.GetString("OWNER", "Unknown"),
		},
		Pending: nodePending{
			Total:    0,
			Messages: 0,
			Tasks:    0,
		},
		Processed: 0,
		Uptime:    0,
	}

	deletedWebhooks = map[string]struct{}{}
	tasks           = make(chan task)

	startTime = time.Now()
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
	Keys []string `json:"keys"`
	Data any      `json:"data"`
}

type task struct {
	chunks [][]string
	data   any
}
