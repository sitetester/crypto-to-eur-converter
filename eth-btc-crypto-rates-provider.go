package main

import (
	"encoding/json"
	"eth-btc-crypto-currencies-converter/provider"
	"eth-btc-crypto-currencies-converter/provider/coinranking"
	"fmt"
	"github.com/fatih/color"
	"log"
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
		fmt.Println("Unknown currency. Supported currencies are " + strings.Join(supportedCurrencies, ", "))
		return
	}

	var cryptoCurrencyRateProvider coinranking.CryptoCurrencyRateProvider
	var currencyRates []provider.CurrencyRate
	var totalCoins int

	currencyRates, totalCoins = cryptoCurrencyRateProvider.ProvideRates(currencyUuidMap[currency])

	color.Green("Done!")
	color.Green("Total rates found: %d\n", totalCoins)
	color.Green("Total rates parsed (excluding `0` amount currencies): %d\n", len(currencyRates))

	fmt.Printf("`FromAmount` is shown in %s\n", currency)
	euroRateStr := cryptoCurrencyRateProvider.ParseEurRate(currencyUuidMap[currency])

	displayInEuro(currencyRates, toFloat(euroRateStr), currencyUuidMap[currency])
}

func displayInEuro(currencyRates []provider.CurrencyRate, euroRate float64, referenceCurrencyUuid string) {
	endpoint := fmt.Sprintf("%scoins?referenceCurrencyUuid=%s", coinranking.ApiUrl, referenceCurrencyUuid)
	for _, rate := range currencyRates {
		amountFloat := toFloat(rate.Amount)
		expectedOutput := provider.ExpectedOutput{
			Endpoint:     endpoint,
			FromCurrency: rate.Currency,
			FromAmount:   amountFloat,
			ToCurrency:   "EUR",
			ToAmount:     amountFloat * euroRate,
		}

		bytes, err := json.Marshal(expectedOutput)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(bytes))
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
