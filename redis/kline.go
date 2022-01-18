package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Coin struct {
	Symbol      string    `json:"symbol"`
	ClosePrices []float64 `json:"closePrices"`
	Svg         string    `json:"svg"`
}

// redis will call this method for mashal
func (coin Coin) MarshalBinary() ([]byte, error) {
	return json.Marshal(coin)
}
func (coin Coin) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &coin)
}

type CoinList struct {
	Symbols []string `json:"symbols"`
}

// redis will call this method for mashal
func (coinList CoinList) MarshalBinary() ([]byte, error) {
	return json.Marshal(coinList)
}
func (coinList CoinList) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &coinList)
}

type ICoinRepository interface {
	Create(Coin) error
	GetBySymbols([]string) []Coin
	GetCoinList() CoinList
	DeleteBySymbols(symbols []string)
	AppendCoinLlist(newSymbol string) error
	PingRedis() (string, error)
}

type CoinDB struct {
	Redis *redis.Client
	Ctx   context.Context
}

const (
	COIN_LIST_KEY = "CoinList"
)

func NewCoinRepository(r *redis.Client, ctx context.Context) ICoinRepository {
	return &CoinDB{
		Redis: r,
		Ctx:   ctx,
	}
}

func (c *CoinDB) Create(coin Coin) error {
	err := c.Redis.Set(c.Ctx, coin.Symbol, coin, 1*time.Hour).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = c.AppendCoinLlist(coin.Symbol)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// value, err := c.Redis.Get(c.Ctx, coin.Symbol).Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("value: %v\n", value)

	return nil
}

func (c *CoinDB) GetBySymbols(symbols []string) []Coin {
	coins := []Coin{}
	for _, symbol := range symbols {
		value, err := c.Redis.Get(c.Ctx, symbol).Result()
		if err != nil {
			fmt.Println(symbol)
			continue
		}
		coin := Coin{}
		err = json.Unmarshal([]byte(value), &coin)
		if err != nil {
			fmt.Println(err)
			continue
		}
		coins = append(coins, coin)
	}
	// fmt.Printf("coins: %v\n", coins)
	return coins
}

func (c *CoinDB) DeleteBySymbols(symbols []string) {
	for _, symbol := range symbols {
		err := c.Redis.Del(c.Ctx, symbol).Err()
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = c.RemoveSymbolFromCoinList(symbol)
		if err != nil {
			fmt.Sprintln("cannot remove: from CoinList", symbol)
			continue
		}
	}
}

func (c *CoinDB) AppendCoinLlist(newSymbol string) error {
	coinList := c.GetCoinList()

	symbolMap := map[string]bool{}
	for _, s := range coinList.Symbols {
		symbolMap[s] = true
	}
	// if coinList already include newSymbol
	if symbolMap[newSymbol] {
		// found
		return nil
	}

	coinList.Symbols = append(coinList.Symbols, newSymbol)
	err := c.Redis.Set(c.Ctx, COIN_LIST_KEY, coinList, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (c *CoinDB) GetCoinList() CoinList {
	coinList, err := c.Redis.Get(c.Ctx, "CoinList").Result()
	if err != nil {
		fmt.Println(err)
		return CoinList{
			Symbols: []string{},
		}
	}
	cl := CoinList{
		Symbols: []string{},
	}
	err = json.Unmarshal([]byte(coinList), &cl)
	if err != nil {
		fmt.Println(err)
		return CoinList{
			Symbols: []string{},
		}
	}
	return cl
}

func (c *CoinDB) RemoveSymbolFromCoinList(symbol string) error {
	coinList := c.GetCoinList()
	newCoinList := []string{}
	for _, c := range coinList.Symbols {
		if c == symbol {
			continue
		}
		newCoinList = append(newCoinList, c)
	}
	coinList.Symbols = newCoinList
	err := c.Redis.Set(c.Ctx, COIN_LIST_KEY, coinList, 0).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (c *CoinDB) PingRedis() (string, error) {
	return c.Redis.Ping(c.Ctx).Result()
}
