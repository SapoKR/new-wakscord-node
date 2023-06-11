package server

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/wakscord/node/config"
	"github.com/wakscord/node/discord"
	"github.com/wakscord/node/utils"
)

func addTask(keys []string, data discord.WebhookParams) error {
	if len(tasks) >= config.Default.MessageQueueSize {
		return errQueueIsFull
	}

	var notDeletedKeys []string

	for _, key := range keys {
		if _, ok := deletedWebhooks[key]; !ok {
			notDeletedKeys = append(notDeletedKeys, key)
		}
	}

	chunks := utils.ChunkSlice(notDeletedKeys, config.Default.MaxConcurrent)

	atomic.AddInt32(&status.Pending.Messages, 1)
	atomic.AddInt32(&status.Pending.Total, 1)

	tasks <- task{
		chunks: chunks,
		data:   data,
	}

	return nil
}

func chunkHandler(keys []string, data discord.WebhookParams) {
	var responseChannel = make(chan discord.Response)

	for _, key := range keys {
		go func(key string, innerChannel chan discord.Response) {
			response := discord.RequestFastHTTP(key, data, 5)

			innerChannel <- response
		}(key, responseChannel)
	}

	for range keys {
		response := <-responseChannel
		if response.Error != nil {
			log.Printf("Uncaught error occurred. Error: %v\n", response.Error)
		}
		if 401 <= response.Code && response.Code <= 404 {
			deletedWebhooks[response.Key] = struct{}{}
		} else if response.Code != 204 {
			log.Printf("Discord returned uncaught status code. Status Code: %d and Body: %s\n", response.Code, response.Body)
		} else {
			status.Processed++
		}
	}

	close(responseChannel)
}

func taskHandler() {
	log.Println("Task handler started")

	for {
		task := <-tasks
		status.Pending.Tasks += int32(len(task.chunks))
		status.Pending.Total += int32(len(task.chunks))

		for _, chunk := range task.chunks {
			chunkHandler(chunk, task.data)
			time.Sleep(time.Second * time.Duration(config.Default.WaitConcurrent))

			status.Pending.Tasks--
			status.Pending.Total--
		}

		status.Pending.Messages--
		status.Pending.Total--
	}
}
