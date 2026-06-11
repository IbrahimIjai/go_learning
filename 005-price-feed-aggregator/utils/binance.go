package utils

import (
	"encoding/json"
	"strconv"
)

func FetchBinance() (float64, error) {
	resp, err := httpClient.Get("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(result.Price, 64)
}
