package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wakscord/new-wakscord-node/env"
	"golang.org/x/exp/maps"
)

func index(ctx *fiber.Ctx) error {
	return ctx.JSON(NodeStatus{
		Info: NodeInfo{
			NodeID: env.GetInt("ID", 1),
			Owner:  env.GetString("Owner", "Unknown"),
		},
		Pending:   NodePending{},
		Processed: processed,
		Deleted:   len(deletedWebhooks),
		Uptime:    int(time.Since(startTime).Seconds()),
	})
}

func request(ctx *fiber.Ctx) error {
	payload := new(RequestPayload)
	if err := ctx.BodyParser(payload); err != nil {
		return err
	}

	go createTasks(payload.Keys, payload.Data)

	return ctx.JSON(fiber.Map{"status": "ok"})
}

func getDeletedWebhooks(ctx *fiber.Ctx) error {
	return ctx.JSON(maps.Keys(deletedWebhooks))
}

func deleteDeletedWebhooks(ctx *fiber.Ctx) error {
	deletedWebhooks = map[string]struct{}{}

	return ctx.SendStatus(fiber.StatusNoContent)
}
