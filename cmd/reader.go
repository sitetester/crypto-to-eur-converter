package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/sitetester/crypto-to-eur-converter/cmd/provider"
)

type EndpointReader struct {
	endpoint string // URL
	client   *http.Client
}

func NewEndpointReader(client *http.Client, endpoint string) *EndpointReader {
	return &EndpointReader{
		endpoint: endpoint,
		client:   client,
	}
}

// Read fetches data from provided endpoint & returns a unified endpoint data
// It handles different JSON formats, e.g., endpoint-btc with string amount, endpoint-eth with float price
func (r *EndpointReader) Read() (provider.EndpointData, error) {
	endpointData := provider.EndpointData{}

	resp, err := r.client.Get(r.endpoint)
	if err != nil {
		return endpointData, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return endpointData, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return endpointData, err
	}

	currency, amount, err := extractCurrencyAndAmount(jsonData)
	if err != nil {
		return endpointData, err
	}

	endpointData.Endpoint = r.endpoint
	endpointData.Currency = currency
	endpointData.Amount = amount

	return endpointData, nil
}

func extractCurrencyAndAmount(jsonData map[string]interface{}) (string, float64, error) {
	var currency string
	var amount float64 = 0

	for _, value := range jsonData {
		strValue := fmt.Sprintf("%v", value)
		// e.g., BTC or USDC
		if isAlpha(strValue) && len(strValue) >= 3 && len(strValue) <= 4 {
			currency = strings.ToUpper(strValue)
		} else {
			amountTmp, err := strconv.ParseFloat(strValue, 64)
			if err != nil {
				return "", 0, err
			}
			amount = amountTmp
		}
	}

	if currency == "" || amount == 0 {
		return "", 0, fmt.Errorf("could not identify currency or amount")
	}

	return currency, amount, nil
}

func isAlpha(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return false
		}
	}
	return true
}
