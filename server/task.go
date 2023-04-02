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
		status.Pending.Messages++
		status.Pending.Total++

		tasks <- task{
			chunks: chunks,
			data:   data,
		}
	}()
}

func chunkHandler(keys []string, data any) {
	var codeChannel = make(chan struct{})

	for _, key := range keys {
		go func(key string, innerChannel chan struct{}) {
			code := discord.RequestHTTP(key, data, 3)
			if code == 401 || code == 403 || code == 404 {
				deletedWebhooks[key] = struct{}{}
			} else if code == 204 {
				status.Processed++
			}

			innerChannel <- struct{}{}
		}(key, codeChannel)
	}

	for range keys {
		<-codeChannel
	}

}

func taskHandler() {
	for {
		task := <-tasks

		status.Pending.Tasks += len(task.chunks)
		status.Pending.Total += len(task.chunks)

		for _, chunk := range task.chunks {
			chunkHandler(chunk, task.data)
			time.Sleep(time.Second * time.Duration(env.GetInt("WAIT_CONCURRENT", 1)))

			status.Pending.Tasks--
			status.Pending.Total--
		}

		status.Pending.Messages--
		status.Pending.Total--
	}
}
