package server

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/wakscord/new-wakscord-node/config"
)

func Run() error {
	status = nodeStatus{
		Info: nodeInfo{
			NodeID: config.Default.ID,
			Owner:  config.Default.Owner,
		},
		Pending: nodePending{
			Total:    0,
			Messages: 0,
			Tasks:    0,
		},
		Processed: 0,
		Uptime:    0,
	}
	tasks = make(chan task, config.Default.MessageQueueSize)
	serverKey = fmt.Sprintf("Bearer %s", config.Default.Key)

	address := fmt.Sprintf("%s:%d", config.Default.Host, config.Default.Port)
	go taskHandler()
	log.Printf("listening on %s\n", address)
	return fasthttp.ListenAndServe(address, requestHandler)
}
