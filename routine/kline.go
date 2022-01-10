package routine

import (
	"fmt"
	"time"

	"github.com/seenark/binance-chart-svg/binance"
	"github.com/seenark/binance-chart-svg/helpers"
	mychart "github.com/seenark/binance-chart-svg/myChart"
	myRedis "github.com/seenark/binance-chart-svg/redis"
)

func FetchHourly(repo myRedis.ICoinRepository, bc binance.BinanceClient, exitChan chan bool) {
	for {
		select {
		case <-exitChan:
			return
		case <-time.After(getTimeRemaining()):
			fmt.Println("do")
			ForceUpdateAll(repo, bc)
		}
	}
}

func ForceUpdateAll(repo myRedis.ICoinRepository, binanceClient binance.BinanceClient) {
	coinList := repo.GetCoinList()
	repo.DeleteBySymbols(coinList.Symbols)
	for _, symbol := range coinList.Symbols {
		coin, err := FetchKline(symbol, binanceClient)
		if err != nil {
			fmt.Println(err)
			continue
		}
		repo.Create(*coin)
	}
}

func FetchKline(symbol string, binanceClient binance.BinanceClient) (*myRedis.Coin, error) {
	start, end := helpers.TimeLast24H()
	kline := binanceClient.GetKLine(symbol, "1h", start, end, 24)
	if len(kline) == 0 {
		return nil, fmt.Errorf("not found")
	}
	closedPrices := []float64{}
	for _, k := range kline {
		closedPrices = append(closedPrices, k.Close)
	}
	buff := mychart.GenerateSVG(closedPrices)
	svg := buff.String()
	newCoin := myRedis.Coin{
		Symbol:      symbol,
		ClosePrices: closedPrices,
		Svg:         svg,
	}
	return &newCoin, nil
}

func getTimeRemaining() time.Duration {
	now := time.Now()
	nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
	// nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, now.Location())
	timeToGo := nextHour.Sub(now)
	fmt.Printf("timeToGo: %v\n", timeToGo)
	return timeToGo
}
