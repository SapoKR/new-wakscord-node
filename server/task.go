package server

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/wakscord/new-wakscord-node/discord"
	"github.com/wakscord/new-wakscord-node/env"
	"github.com/wakscord/new-wakscord-node/utils"
)

func addTask(keys []string, data any) {
	var notDeletedKeys []string

	for _, key := range keys {
		if _, ok := deletedWebhooks[key]; !ok {
			notDeletedKeys = append(notDeletedKeys, key)
		}
	}

	chunks := utils.ChunkSlice(notDeletedKeys, env.GetInt("MAX_CONCURRENT", 500))

	go func() {
		atomic.AddInt64(&status.Pending.Messages, 1)
		atomic.AddInt64(&status.Pending.Total, 1)

		tasks <- task{
			chunks: chunks,
			data:   data,
		}
	}()
}

func chunkHandler(keys []string, data any) {
	var responseChannel = make(chan discord.Response)

	for _, key := range keys {
		go func(key string, innerChannel chan discord.Response) {
			response := discord.RequestFastHTTP(key, data, 3)

			innerChannel <- response
		}(key, responseChannel)
	}

	for range keys {
		response := <-responseChannel
		if response.Error != nil {
			log.Printf("Uncaught error occurred. Error: %v", response.Error)
		}
		if 401 <= response.Code && response.Code <= 404 {
			deletedWebhooks[response.Key] = struct{}{}
		} else if response.Code != 204 {
			log.Printf("Discord returned uncaught status code. Status Code: %d and Body: %s\n", response.Code, response.Body)
		} else {
			atomic.AddInt64(&status.Processed, 1)
		}
	}
}

func taskHandler() {
	for {
		task := <-tasks

		atomic.AddInt64(&status.Pending.Tasks, int64(len(task.chunks)))
		atomic.AddInt64(&status.Pending.Total, int64(len(task.chunks)))

		for _, chunk := range task.chunks {
			chunkHandler(chunk, task.data)
			time.Sleep(time.Second * time.Duration(env.GetInt("WAIT_CONCURRENT", 1)))

			atomic.AddInt64(&status.Pending.Tasks, -1)
			atomic.AddInt64(&status.Pending.Total, -1)
		}

		atomic.AddInt64(&status.Pending.Messages, -1)
		atomic.AddInt64(&status.Pending.Total, -1)
	}
}
