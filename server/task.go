package server

import (
	"time"

	"github.com/wakscord/new-wakscord-node/discord"
	"github.com/wakscord/new-wakscord-node/env"
	"github.com/wakscord/new-wakscord-node/utils"
)

func addTask(keys []string, data any) {
	chunks := utils.ChunkSlice(keys, env.GetInt("MAX_CONCURRENT", 500))

	go func() {
		tasks <- task{
			chunks: chunks,
			data:   data,
		}
	}()
}

func chunkHandler(keys []string, data any) {
	var codeChannel = make(chan int)
	status.Pending.Tasks++
	status.Pending.Total++

	for _, key := range keys {
		go func(key string, innerChannel chan int) {
			code := discord.Request(key, data, 3)
			if code == 401 || code == 403 || code == 404 {
				deletedWebhooks[key] = struct{}{}
			} else if code == 204 {
				status.Processed++
			}

			innerChannel <- code
		}(key, codeChannel)
	}

	for range keys {
		<-codeChannel
	}

	status.Pending.Tasks--
	status.Pending.Total--
}

func taskHandler() {
	for {
		task := <-tasks
		status.Pending.Messages++
		status.Pending.Total++

		for _, chunk := range task.chunks {
			chunkHandler(chunk, task.data)
			time.Sleep(time.Second * time.Duration(env.GetInt("WAIT_CONCURRENT", 0)))
		}

		status.Pending.Messages--
		status.Pending.Total--
	}
}
