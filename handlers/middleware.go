package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	myRedis "github.com/seenark/binance-chart-svg/redis"
)

func checkRedisStatus(coinRepo myRedis.ICoinRepository) func(*fiber.Ctx) error {
	fmt.Println("Middleware")
	return func(c *fiber.Ctx) error {
		pong, err := coinRepo.PingRedis()
		if err != nil {
			fmt.Println("Redis Ping Error:", err)
		}
		if pong != "PONG" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to Redis"})
		}
		return c.Next()
	}

}
