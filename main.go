package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/seenark/binance-chart-svg/binance"
	"github.com/seenark/binance-chart-svg/config"
	"github.com/seenark/binance-chart-svg/handlers"
	myRedis "github.com/seenark/binance-chart-svg/redis"
	"github.com/seenark/binance-chart-svg/routine"
)

var exitChan = make(chan bool)
var isRoutineRunning = false

func main() {
	ctx := context.TODO()
	cfg := config.GetConfig()
	app := fiber.New()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Username: "",
		Password: "",
		DB:       0,
	})

	coinRepo := myRedis.NewCoinRepository(redisClient, ctx)
	binanceClient := binance.NewBinanceClient(&http.Client{})

	api := app.Group("/api")
	handlers.SetCoinHandler(api, coinRepo, binanceClient)

	routineGroup := app.Group("/routine")
	routineGroup.Post("/start", func(c *fiber.Ctx) error {
		if !isRoutineRunning {
			go routine.FetchHourly(coinRepo, *binanceClient, exitChan)
			isRoutineRunning = true
			fmt.Println("start")
		}
		return c.SendStatus(fiber.StatusOK)
	})
	routineGroup.Post("/stop", func(c *fiber.Ctx) error {
		if isRoutineRunning {
			exitChan <- true
			isRoutineRunning = false
			fmt.Println("stop")
		}
		return c.SendStatus(fiber.StatusOK)
	})
	routineGroup.Post("/force-update", func(c *fiber.Ctx) error {
		routine.ForceUpdateAll(coinRepo, *binanceClient)
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OMG You got me")
	})
	go routine.FetchHourly(coinRepo, *binanceClient, exitChan)
	isRoutineRunning = true

	app.Listen(fmt.Sprintf(":%d", cfg.Port))
}
