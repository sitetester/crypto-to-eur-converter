package coinranking

import (
	"encoding/json"
	"eth-btc-crypto-currencies-converter/provider"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CryptoCurrencyRateProvider struct{}

const BaseUrl = "https://coinranking.com/"
const ApiUrl = BaseUrl + "api/v2/"

var currencyRates []provider.CurrencyRate

func (cryptoCurrencyRateProvider *CryptoCurrencyRateProvider) ParseEurRate(currency string) string {
	toEuroUrl := fmt.Sprintf("%s/coin/%s?referenceCurrencyUuid=%s", ApiUrl, currency, "5k-_VTxqtCEI")
	body := makeRequest(toEuroUrl)

	var apiResponse EurApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		panic(err)
	}

	return apiResponse.Data.Coin.Price
}

func (cryptoCurrencyRateProvider *CryptoCurrencyRateProvider) ProvideRates(referenceCurrencyUuid string) ([]provider.CurrencyRate, int) {
	apiResponse := parseUrl(referenceCurrencyUuid, 0)
	totalCoins := apiResponse.Data.Stats.TotalCoins

	apiResponseChan := make(chan ApiResponse, totalCoins)
	offset := 0
	totalPages := 0
	for totalPages <= apiResponse.Data.Stats.TotalCoins/50 {
		totalPages += 1
		go parseUrlWithChannel(referenceCurrencyUuid, offset, apiResponseChan)
		offset += 50
	}

channelLoop:
	for {
		select {
		case apiResponse := <-apiResponseChan:
			parseRates(apiResponse.Data.Coins)

			if len(currencyRates) == totalCoins {
				break channelLoop
			}
		}
	}

	return filterRates(), totalCoins
}

// this is much faster than applying check inside `parseRates()`
func filterRates() []provider.CurrencyRate {
	var filteredRates []provider.CurrencyRate
	for _, currencyRate := range currencyRates {
		if currencyRate.Amount != "" {
			filteredRates = append(filteredRates, currencyRate)
		}
	}

	return filteredRates
}

func parseUrlWithChannel(referenceCurrencyUuid string, offset int, apiResponseChan chan ApiResponse) {
	apiResponseChan <- parseUrl(referenceCurrencyUuid, offset)
}

func parseUrl(referenceCurrencyUuid string, offset int) ApiResponse {
	url := fmt.Sprintf("%scoins?referenceCurrencyUuid=%s&timePeriod=24h&offset=%d&limit=50", ApiUrl, referenceCurrencyUuid, offset)
	fmt.Printf("Parsing URL: %s\n", url)

	body := makeRequest(url)
	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		panic(err)
	}

	return apiResponse
}

func makeRequest(url string) []byte {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func parseRates(coins []Coin) {
	for _, coin := range coins {
		currencyRates = append(currencyRates, provider.CurrencyRate{
			Currency: coin.Symbol,
			Amount:   coin.Price,
		})
	}
}

type EurApiResponse struct {
	Status string
	Data   EurData
}

type ApiResponse struct {
	Status string
	Data   Data
}

type EurData struct {
	Coin Coin
}

type Data struct {
	Stats Stats
	Coins []Coin
}

type Stats struct {
	TotalCoins int `json:"totalCoins"`
}

type Coin struct {
	Uuid   string
	Symbol string
	Name   string
	Color  string
	Price  string
}
