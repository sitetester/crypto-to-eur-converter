//go:build integration

package provider

import (
	"testing"
)

// TestGetEurRate is an integration test because it calls external service (CryptoCompare API)
func TestGetEurRate_Success(t *testing.T) {
	endpointData := EndpointData{
		Endpoint: "http://localhost:8080/endpoint-btc",
		Currency: BTC,
		Amount:   0.00534315,
	}

	rp := NewCryptoCompareRatesProvider(DefaultHTTPClient, CryptoCompareBaseURL)
	rate, err := rp.GetEurRate(endpointData)
	if err != nil {
		t.Fatalf("GetEurRate failed: %v", err)
	}

	if rate == 0 {
		t.Fatalf("unexpected eur rate 0")
	}
}
