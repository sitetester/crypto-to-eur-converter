package main

import (
	"fmt"

	"github.com/sitetester/crypto-to-eur-converter/cmd/provider"
)

// UnifiedEurOutput represents the converted EUR amount for an endpoint.
type UnifiedEurOutput struct {
	Endpoint string `json:"endpoint"` // Endpoint URL, e.g., "https://api.example.com/btc"
	Currency string `json:"currency"` // Always "EUR"
	Amount   string `json:"amount"`   // EUR amount formatted to 2 decimals
}

// ToUnifiedEurOutput converts endpoint data to EUR using the given rates provider.
func ToUnifiedEurOutput(
	ratesProvider provider.RatesProvider,
	endpointData provider.EndpointData,
) (UnifiedEurOutput, error) {
	rate, err := ratesProvider.GetEurRate(endpointData)
	if err != nil {
		return UnifiedEurOutput{}, fmt.Errorf("error getting eur rate from provider: %w", err)
	}

	eurAmount := endpointData.Amount * rate

	return UnifiedEurOutput{
		Endpoint: endpointData.Endpoint,
		Currency: provider.EUR,
		Amount:   fmt.Sprintf("%.2f", eurAmount),
	}, nil
}
