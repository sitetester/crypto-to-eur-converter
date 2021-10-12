package main

import (
	"eth-btc-crypto-currencies-converter/provider"
	"eth-btc-crypto-currencies-converter/provider/coinranking"
	"fmt"
	"github.com/fatih/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	currencyUuidMap := make(map[string]string)
	currencyUuidMap["BTC"] = "Qwsogvtv82FCd"
	currencyUuidMap["ETH"] = "razxDUgYGNAdQ"

	supportedCurrencies := make([]string, 0, len(currencyUuidMap))
	for k := range currencyUuidMap {
		supportedCurrencies = append(supportedCurrencies, k)
	}

	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please provide a crypto currency argument. e.g. ETH")
		return
	}

	currency := args[1]
	if !contains(supportedCurrencies, strings.ToUpper(currency)) {
		fmt.Println("Unknown currency provided. Supported currencies are " + strings.Join(supportedCurrencies, ", "))
		return
	}

	var cryptoCurrencyRateProvider coinranking.CryptoCurrencyRateProvider
	var currencyRates []provider.CurrencyRate
	var totalCoins int

	currencyRates, totalCoins = cryptoCurrencyRateProvider.ProvideRates(currencyUuidMap[currency])

	color.Green("Done!")
	color.Green("Total rates found: %d\n", totalCoins)
	color.Green("Total rates parsed: %d\n", len(currencyRates))

	euroRateStr := cryptoCurrencyRateProvider.ParseEurRate(currencyUuidMap[currency])
	println("euroRateStr", euroRateStr)
	displayInEuro(currencyRates, toFloat(euroRateStr))
}

// https://yourbasic.org/golang/round-float-2-decimal-places/
func displayInEuro(currencyRates []provider.CurrencyRate, euroRate float64) {
	endpoint := coinranking.BaseUrl
	i := 0
	for _, rate := range currencyRates {
		i += 1
		if i > 5 {
			break
		}

		euroOutput := provider.EuroOutput{
			Endpoint: endpoint,
			Currency: "EUR",
			Amount:   math.Round(toFloat(rate.Amount)*euroRate*100) / 100, // round to nearest
		}

		fmt.Printf("%s, %+v\n", rate.Currency, euroOutput)
	}
}

func toFloat(strRate string) float64 {
	if strRate == "" {
		return 0
	}

	euroRate, err := strconv.ParseFloat(strRate, 64)
	if err != nil {
		log.Fatal(err)
	}

	return euroRate
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
