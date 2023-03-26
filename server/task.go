package server

import (
	"github.com/wakscord/new-wakscord-node/env"
	"github.com/wakscord/new-wakscord-node/requester"
	"github.com/wakscord/new-wakscord-node/utils"
)

func createTasks(keys []string, data any) {
	var (
		keysToRequest []string
		waiter        = make(chan struct{})
	)

	for _, key := range keys {
		if _, ok := deletedWebhooks[key]; !ok {
			keysToRequest = append(keysToRequest, key)
		}
	}

	chunks := utils.ChunkSlice(keysToRequest, env.GetInt("MAX_CONCURRENT", 500))

	for _, chunk := range chunks {
		go func(innerKeys []string, innerWaiter chan struct{}) {
			defer func() {
				innerWaiter <- struct{}{}
			}()

			startTask(innerKeys, data)
		}(chunk, waiter)
	}

	for range chunks {
		<-waiter
	}
}

func startTask(keys []string, data any) {
	waiter := make(chan struct{})

	for _, key := range keys {
		go func(innerKey string, innerWaiter chan struct{}) {
			defer func() {
				innerWaiter <- struct{}{}
			}()

			code := requester.Request(innerKey, data)
			if code == 401 || code == 403 || code == 404 {
				deletedWebhooks[innerKey] = struct{}{}
			} else {
				processed++
			}
		}(key, waiter)
	}

	for range keys {
		<-waiter
	}
}
