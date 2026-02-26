//go:build integration

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sitetester/crypto-to-eur-converter/cmd/provider"
)

func createTestServer(data map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}))
}

func testEndpoint(t *testing.T, data map[string]interface{}, wantCurrency string, wantAmount float64) {
	server := createTestServer(data)
	defer server.Close()

	reader := NewEndpointReader(provider.DefaultHTTPClient, server.URL)
	endpointData, err := reader.Read()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if endpointData.Currency != wantCurrency {
		t.Errorf("Expected currency %s, got %s", wantCurrency, endpointData.Currency)
	}

	if endpointData.Amount != wantAmount {
		t.Errorf("Expected amount %f, got %f", wantAmount, endpointData.Amount)
	}
}

func TestEndpointBtc(t *testing.T) {
	testEndpoint(
		t,
		map[string]interface{}{
			"currency": "BTC",
			"amount":   "0.005343150",
		},
		provider.BTC,
		0.005343150,
	)
}

func TestEndpointEth(t *testing.T) {
	testEndpoint(
		t,
		map[string]interface{}{
			"curr_iso": "ETH",
			"price":    10,
		},
		provider.ETH,
		10.0,
	)
}

func TestEndpointSol(t *testing.T) {
	testEndpoint(
		t,
		map[string]interface{}{
			"currency_iso": "SOL",
			"price":        123.45,
		},
		provider.SOL,
		123.45,
	)
}

func TestEndpointRandomFieldsOrder(t *testing.T) {
	tests := []struct {
		name string
		data map[string]interface{}
	}{
		{
			name: "currency first",
			data: map[string]interface{}{
				"some_random_name": "ABC",
				"cur_rice":         "123.45",
			},
		},
		{
			name: "currency second",
			data: map[string]interface{}{
				"currency_amount":  "123.45",
				"some_random_name": "ABC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testEndpoint(t, tt.data, "ABC", 123.45)
		})
	}
}
