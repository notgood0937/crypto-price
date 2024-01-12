package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron"
	"time"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})
	fmt.Println("redis conn success")
	c := cron.New()
	c.AddFunc("@every 5s", func() {
		getPrice(rdb)
	})

	c.Start()
	select {}

}

func getPrice(rdb *redis.Client) {
	fmt.Println("I am runnning task.", time.Now())
	var (
		apiKey    = ""
		secretKey = ""
	)
	////futuresClient := binance.NewFuturesClient(apiKey, secretKey) // USDT-M Futures
	//client := binance.NewFuturesClient(apiKey, secretKey)
	//res1, err1 := client.NewChangeLeverageService().Leverage(20).Symbol("BTCUSDT").Do(context.Background())
	//if err1 != nil {
	//	fmt.Println(err1)
	//	return
	//}
	//jsonData, err := json.Marshal(res1)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//println(string(jsonData))
	//res, err := client.NewCreateOrderService().Symbol("BTCUSDT").Side("buy").Type("market").Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//jsonData, err := json.Marshal(res)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//println(string(jsonData))
	//
	client := binance.NewClient(apiKey, secretKey)
	res, err := client.NewListPricesService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(jsonData))
	slice := make([]binance.SymbolPrice, 0)
	errData := json.Unmarshal([]byte(jsonData), &slice)
	if errData != nil {
		panic(errData)
	}
	pip := rdb.Pipeline()
	ctx := context.Background()
	for _, market := range slice {
		pip.Set(ctx, market.Symbol, market.Price, 0)
		//fmt.Println(market.Symbol, market.Price)
	}
	_, err1 := pip.Exec(ctx)
	if err1 != nil {
		return
	}

	//
	//wsKlineHandler := func(event binance.WsAllMiniMarketsStatEvent) {
	//	jsonData, err := json.Marshal(event)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	slice := make([]binance.WsMiniMarketsStatEvent, 0)
	//	errData := json.Unmarshal([]byte(jsonData), &slice)
	//	if errData != nil {
	//		panic(errData)
	//	}
	//
	//	// 打印切片中的元素。
	//	for _, market := range slice {
	//		jsonBytes, err := json.Marshal(market)
	//		if err != nil {
	//			// 处理错误
	//			fmt.Println(err)
	//			return
	//		}
	//		fmt.Println(string(jsonBytes))
	//		_, err = c.Do("Set", market.Symbol, string(jsonBytes))
	//		if err != nil {
	//			fmt.Println(err)
	//			return
	//		}
	//		//fmt.Println(market.Symbol, market.LastPrice)
	//	}
	//}
	//errHandler := func(err error) {
	//	fmt.Println(err)
	//}
	////binance.WebsocketKeepalive = true
	//doneC, _, err := binance.WsAllMiniMarketsStatServe(wsKlineHandler, errHandler)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//<-doneC

	//doneC, stopC, err := binance.WsAllMiniMarketsStatServe(func(event binance.WsAllMiniMarketsStatEvent) {
	//	fmt.Println(event)
	//}, func(err error) {
	//	fmt.Println(err)
	//	return
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//stopC <- struct{}{}
	//<-doneC

}
