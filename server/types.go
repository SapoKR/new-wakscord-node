package server

import (
	"time"
)

var (
	processed int

	deletedWebhooks = map[string]struct{}{}

	startTime = time.Now()
)

type NodeStatus struct {
	Info      NodeInfo    `json:"info"`
	Pending   NodePending `json:"pending"`
	Processed int         `json:"processed"`
	Deleted   int         `json:"deleted"`
	Uptime    int         `json:"uptime"`
}

type NodeInfo struct {
	NodeID int    `json:"node_id"`
	Owner  string `json:"owner"`
}

type NodePending struct {
	Total    int `json:"total"`
	Messages int `json:"messages"`
	Tasks    int `json:"tasks"`
}

type RequestPayload struct {
	Keys []string `json:"keys"`
	Data any      `json:"data"`
}
