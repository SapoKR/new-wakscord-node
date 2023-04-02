package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/wakscord/new-wakscord-node/env"
)

func Run() error {
	app := fiber.New(fiber.Config{
		AppName: "Wakscord Node",
	})
	key := fmt.Sprintf("Bearer %s", env.GetString("KEY", "wakscord"))

	app.Get("/", index)
	app.Use(func(ctx *fiber.Ctx) error {
		if ctx.Path() == "/" || ctx.GetReqHeaders()["Authorization"] == key {
			return ctx.Next()
		}

		return fiber.ErrUnauthorized
	})
	app.Post("/request", request)
	app.Get("/deletedWebhooks", getDeletedWebhooks)
	app.Delete("/deletedWebhooks", deleteDeletedWebhooks)

	address := fmt.Sprintf("%s:%d", env.GetString("HOST", "0.0.0.0"), env.GetInt("PORT", 3000))

	go taskHandler()

	return app.Listen(address)
}
