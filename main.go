package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/seenark/binance-chart-svg/config"
)

func main() {
	config := config.GetConfig()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(fmt.Sprintf("localhost:%d", config.Port))
}
