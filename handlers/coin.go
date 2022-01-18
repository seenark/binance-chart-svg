package handlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/seenark/binance-chart-svg/binance"
	myRedis "github.com/seenark/binance-chart-svg/redis"
	"github.com/seenark/binance-chart-svg/routine"
)

type CoinHander struct {
	repo          myRedis.ICoinRepository
	binanceClient *binance.BinanceClient
}

func SetCoinHandler(app fiber.Router, coinRepo myRedis.ICoinRepository, b *binance.BinanceClient) {
	handler := CoinHander{
		repo:          coinRepo,
		binanceClient: b,
	}
	app.Get("/coins", checkRedisStatus(coinRepo), handler.GetCoins)
	app.Get("/svg/:symbol", checkRedisStatus(coinRepo), handler.GetSVG)
	app.Get("/watch-list", checkRedisStatus(coinRepo), handler.GetWatchList)

}

func (ch CoinHander) GetCoins(ctx *fiber.Ctx) error {
	symbols := splitSymbols(ctx)
	coins, err := getCoinAndFetchBySymbols(symbols, ch.repo, *ch.binanceClient)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(coins)
}

func (ch CoinHander) GetSVG(ctx *fiber.Ctx) error {
	symbol := ctx.Params("symbol")
	coins, err := getCoinAndFetchBySymbols([]string{symbol}, ch.repo, *ch.binanceClient)
	if err != nil {
		return err
	}

	if len(coins) == 0 {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	coin := coins[0]
	return ctx.Type(".svg").Status(fiber.StatusOK).Send([]byte(coin.Svg))
}

func (ch CoinHander) GetWatchList(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(ch.repo.GetCoinList())
}

func getCoinAndFetchBySymbols(symbols []string, repo myRedis.ICoinRepository, bc binance.BinanceClient) ([]myRedis.Coin, error) {
	coins := repo.GetBySymbols(symbols)
	notFoundList := []string{}
	symbolMap := map[string]bool{}
	for _, s := range coins {
		symbolMap[s.Symbol] = true
	}
	// fmt.Printf("symbolMap: %v\n", symbolMap)
	for _, s := range symbols {
		if _, ok := symbolMap[s]; !ok {
			notFoundList = append(notFoundList, s)
		}
	}
	fmt.Printf("notFoundList: %v\n", notFoundList)
	// fetch kline for notFoundList
	for _, s := range notFoundList {

		newCoin, err := routine.FetchKline(s, bc)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			repo.DeleteBySymbols([]string{s})
			continue
			// return nil, err
		}
		repo.Create(*newCoin)
		coins = append(coins, *newCoin)
	}
	return coins, nil
}

func splitSymbols(ctx *fiber.Ctx) []string {
	symbolsString := ctx.Query("symbols")
	symbols := strings.Split(symbolsString, ",")

	for index, v := range symbols {
		if v == "" {
			continue
		}
		symbols[index] = strings.Trim(v, " ")
	}
	return symbols
}
