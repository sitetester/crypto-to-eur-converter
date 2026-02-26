package provider

import (
	"net/http"
	"time"
)

const (
	BTC = "BTC"
	ETH = "ETH"
	SOL = "SOL"
	EUR = "EUR"
)

var DefaultHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// EndpointData keeps `Amount` as float64 for multiplication purpose
type EndpointData struct {
	Endpoint string  `json:"endpoint"` // URL
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

// RatesProvider defines the interface for fetching crypto to EUR exchange rates.
// Thus, we could have multiple implementations, e.g., CryptoCompare, CoinGecko, ...
// allowing the main program to work with various EUR exchange rates
type RatesProvider interface {
	GetEurRate(endpointData EndpointData) (float64, error)
}
