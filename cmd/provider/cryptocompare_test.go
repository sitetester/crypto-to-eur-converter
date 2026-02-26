package provider

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetEurRate_WithMockServer_Timeout(t *testing.T) {
	// The mock server handles all paths by default, so it won't return 404.
	// This will be something like `http://127.0.0.1:54321/data/pricemulti?fsyms=BTC&tsyms=EUR`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond) // Just long enough to trigger timeout
	}))
	defer server.Close()

	provider := NewCryptoCompareRatesProvider(
		&http.Client{Timeout: 100 * time.Millisecond},
		server.URL,
	)

	endpointData := EndpointData{
		Currency: BTC,
		Amount:   1.0,
	}

	_, err := provider.GetEurRate(endpointData)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "timely eur rates") {
		t.Errorf("error message should mention 'timely eur rates', got: %v", err)
	}

	// Verify it's a timeout error
	var netErr interface{ Timeout() bool }
	if !errors.As(err, &netErr) || !netErr.Timeout() {
		t.Errorf("expected timeout error, got: %v", err)
	}
}

func TestGetEurRate_WithMockServer_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return valid JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"BTC":{"EUR":50000.0}}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	provider := NewCryptoCompareRatesProvider(
		&http.Client{Timeout: 5 * time.Second},
		server.URL,
	)

	endpointData := EndpointData{
		Currency: BTC,
		Amount:   1.0,
	}

	rate, err := provider.GetEurRate(endpointData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rate != 50000.0 {
		t.Errorf("expected rate 50000.0, got %f", rate)
	}
}

func TestGetEurRate_WithMockServer_CurrencyNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"BTC123":{"EUR":1.25}}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	provider := NewCryptoCompareRatesProvider(
		&http.Client{Timeout: 5 * time.Second},
		server.URL,
	)

	endpointData := EndpointData{
		Currency: BTC,
		Amount:   1.0,
	}

	rate, err := provider.GetEurRate(endpointData)
	if err == nil {
		t.Fatal("expected error for currency not found, got nil")
	}

	if rate != 0 {
		t.Errorf("expected rate 0, got %f", rate)
	}

	if !strings.Contains(err.Error(), "not found in response") {
		t.Errorf("error should mention `%s`, got: %v",
			fmt.Errorf("currency %s not found in response", endpointData.Currency),
			err.Error(),
		)
	}
}

func TestGetEurRate_WithMockServer_EURNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"BTC":{"EUR123":1.25}}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	provider := NewCryptoCompareRatesProvider(
		&http.Client{Timeout: 5 * time.Second},
		server.URL,
	)

	endpointData := EndpointData{
		Currency: BTC,
		Amount:   1.0,
	}

	rate, err := provider.GetEurRate(endpointData)
	if err == nil {
		t.Fatal("expected error for EUR not found, got nil")
	}

	if rate != 0 {
		t.Errorf("expected rate 0, got %f", rate)
	}

	if !strings.Contains(err.Error(), "EUR rate not found") {
		t.Errorf("EUR rate not found, got: %s",
			err.Error(),
		)
	}
}
