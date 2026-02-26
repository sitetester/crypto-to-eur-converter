// Package main implements a cryptocurrency to EUR converter.
// It reads crypto amounts from HTTP endpoints and converts them to EUR using latest exchange rates.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sitetester/crypto-to-eur-converter/cmd/provider"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run ./cmd <endpoint-url>")
		fmt.Println("Example: go run ./cmd http://localhost:8080/endpoint-btc")
		return
	}
	url := os.Args[1]

	reader := NewEndpointReader(provider.DefaultHTTPClient, url)
	endpointData, err := reader.Read()
	if err != nil {
		fmt.Printf("Error reading endpoint: %v\n", err)
		return
	}

	// This can be replaced with any other provider.
	ratesProvider := provider.NewCryptoCompareRatesProvider(
		provider.DefaultHTTPClient,
		provider.CryptoCompareBaseURL,
	)
	unifiedEurOutput, err := ToUnifiedEurOutput(ratesProvider, endpointData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	jsonBytes, _ := json.MarshalIndent(unifiedEurOutput, "", "  ")
	fmt.Println(string(jsonBytes))
}
