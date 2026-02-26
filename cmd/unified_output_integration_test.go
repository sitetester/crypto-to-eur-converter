//go:build integration

package main

import (
	"testing"

	"github.com/sitetester/crypto-to-eur-converter/cmd/provider"
)

// Verifies it returns the unified structure with amounts converted to EUR
func TestToUnifiedEurOutputs(t *testing.T) {
	rp := provider.NewCryptoCompareRatesProvider(
		provider.DefaultHTTPClient,
		provider.CryptoCompareBaseURL,
	)
	endpointData := provider.EndpointData{
		Endpoint: "http://localhost:8080/endpoint-btc",
		Currency: provider.BTC,
		Amount:   0.00534315,
	}
	unifiedEurOutput, err := ToUnifiedEurOutput(rp, endpointData)
	if err != nil {
		t.Fatalf("ToUnifiedEurOutput failed: %v", err)
	}

	if unifiedEurOutput.Endpoint != endpointData.Endpoint {
		t.Errorf("Expected %s, got %s", endpointData.Endpoint, unifiedEurOutput.Endpoint)
	}

	if unifiedEurOutput.Currency != provider.EUR {
		t.Errorf("Expected %s, got %s", provider.EUR, unifiedEurOutput.Currency)
	}

	if unifiedEurOutput.Amount == "" {
		t.Errorf("%s has empty amount", unifiedEurOutput.Endpoint)
	}
}
