package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const CryptoCompareBaseURL = "https://min-api.cryptocompare.com"

type CryptoCompareRatesProvider struct {
	client  *http.Client
	baseURL string
}

func NewCryptoCompareRatesProvider(client *http.Client, baseURL string) *CryptoCompareRatesProvider {
	return &CryptoCompareRatesProvider{
		client:  client,
		baseURL: baseURL,
	}
}

func (rp *CryptoCompareRatesProvider) GetEurRate(endpointData EndpointData) (float64, error) {
	url := fmt.Sprintf(
		"%s/data/pricemulti?fsyms=%s&tsyms=EUR",
		rp.baseURL,
		endpointData.Currency,
	)

	resp, err := rp.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error getting timely eur rates: %w", err)
	}
	defer resp.Body.Close()

	// e.g., {"BTC":{"EUR":57759.76}}
	var apiData map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&apiData); err != nil {
		return 0, err
	}

	if eurRateMap, ok := apiData[endpointData.Currency]; ok {
		if eurRate, exists := eurRateMap[EUR]; exists {
			return eurRate, nil
		}
		return 0, fmt.Errorf("EUR rate not found for currency %s", endpointData.Currency)
	}

	return 0, fmt.Errorf("currency %s not found in response", endpointData.Currency)
}
