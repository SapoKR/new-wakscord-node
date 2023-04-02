package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/maps"
)

func index(ctx *fiber.Ctx) error {
	status.Deleted = len(deletedWebhooks)
	status.Uptime = int(time.Since(startTime).Seconds())
	return ctx.JSON(status)
}

func request(ctx *fiber.Ctx) error {
	payload := new(requestPayload)
	if err := ctx.BodyParser(payload); err != nil {
		return err
	}

	addTask(payload.Keys, payload.Data)

	return ctx.JSON(fiber.Map{"status": "ok"})
}

func getDeletedWebhooks(ctx *fiber.Ctx) error {
	return ctx.JSON(maps.Keys(deletedWebhooks))
}

func deleteDeletedWebhooks(ctx *fiber.Ctx) error {
	deletedWebhooks = map[string]struct{}{}

	return ctx.SendStatus(fiber.StatusNoContent)
}
