package utils

import (
	"encoding/json"
	"strconv"
)

func FetchCoinbase() (float64, error) {
	resp, err := httpClient.Get("https://api.coinbase.com/v2/prices/BTC-USD/spot")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Amount string `json:"amount"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(result.Data.Amount, 64)
}
