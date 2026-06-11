package utils

import (
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 5 * time.Second}

type PriceUpdate struct {
	Provider string
	Symbol   string
	Price    float64
	Time     time.Time
}
