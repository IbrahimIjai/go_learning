package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func FetchPyth() (float64, error) {
	url := "https://hermes.pyth.network/v2/updates/price/latest?ids[]=e62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43"
	resp, err := httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result struct {
		Parsed []struct {
			Price struct {
				Price string `json:"price"`
				Expo  int    `json:"expo"`
			} `json:"price"`
		} `json:"parsed"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}
	if len(result.Parsed) == 0 {
		return 0, fmt.Errorf("pyth: empty response")
	}

	raw, err := strconv.ParseFloat(result.Parsed[0].Price.Price, 64)
	if err != nil {
		return 0, err
	}

	// Pyth prices are scaled: price * 10^expo
	expo := result.Parsed[0].Price.Expo
	price := raw
	if expo < 0 {
		for i := 0; i > expo; i-- {
			price /= 10
		}
	}
	return price, nil
}
