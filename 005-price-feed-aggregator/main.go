package main

import (
	"fmt"
	"time"
	"go_learning/005-price-feed-aggregator/utils"
)

func average(prices map[string]float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	var total float64
	for _, price := range prices {
		total += price
	}
	return total / float64(len(prices))
}

func pollProvider(name string, interval time.Duration, fetchFn func() (float64, error), updates chan<- utils.PriceUpdate) {
	for {
		price, err := fetchFn()
		if err != nil {
			fmt.Printf("[%s] %s fetch error: %v\n", time.Now().Format("15:04:05"), name, err)
		} else {
			updates <- utils.PriceUpdate{
				Provider: name,
				Symbol:   "BTC/USD",
				Price:    price,
				Time:     time.Now(),
			}
		}
		time.Sleep(interval)
	}
}

func main() {
	updates := make(chan utils.PriceUpdate)

	go pollProvider("Binance", 3*time.Second, utils.FetchBinance, updates)
	go pollProvider("Coinbase", 5*time.Second, utils.FetchCoinbase, updates)
	go pollProvider("Pyth", 2*time.Second, utils.FetchPyth, updates)

	latestPrices := make(map[string]float64)
	timeout := time.After(30 * time.Second)

	for {
		select {
		case update := <-updates:
			latestPrices[update.Provider] = update.Price
			avg := average(latestPrices)
			fmt.Printf(
				"[%s] %-8s %s price: %10.2f | Aggregated average: %.2f\n",
				update.Time.Format("15:04:05"),
				update.Provider,
				update.Symbol,
				update.Price,
				avg,
			)

		case <-timeout:
			fmt.Println("Stopping price feed aggregator.")
			return
		}
	}
}
