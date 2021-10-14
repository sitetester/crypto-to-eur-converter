package main

import (
	"encoding/json"
	"eth-btc-crypto-currencies-converter/helper"
	"eth-btc-crypto-currencies-converter/provider/coinranking"
	"fmt"
	"github.com/fatih/color"
	"os"
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
	if !helper.Contains(supportedCurrencies, strings.ToUpper(currency)) {
		fmt.Println("Unknown currency. Supported currencies are " + strings.Join(supportedCurrencies, ", "))
		return
	}

	var cryptoCurrencyRateProvider coinranking.CryptoCurrencyRateProvider
	var currencyRates []coinranking.CurrencyRate
	var totalCoins int

	currencyRates, totalCoins = cryptoCurrencyRateProvider.ProvideRatesForCurrency(currencyUuidMap[currency])

	color.Green("Done!")
	color.Green("Total rates found: %d\n", totalCoins)
	color.Green("Total rates parsed (excluding `0` amount currencies): %d\n", len(currencyRates))

	fmt.Printf("`FromAmount` is shown in %s\n", currency)
	euroRateStr := cryptoCurrencyRateProvider.ParseEurRate(currencyUuidMap[currency])

	displayInEuro(currencyRates, helper.ToFloat(euroRateStr), currencyUuidMap[currency])
}

func displayInEuro(currencyRates []coinranking.CurrencyRate, euroRate float64, referenceCurrencyUuid string) {
	endpoint := fmt.Sprintf("%scoins?referenceCurrencyUuid=%s", coinranking.ApiUrl, referenceCurrencyUuid)

	for _, rate := range currencyRates {
		amountFloat := helper.ToFloat(rate.Amount)

		expectedOutput := ExpectedOutput{
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

type ExpectedOutput struct {
	Endpoint     string
	FromCurrency string
	FromAmount   float64
	ToCurrency   string
	ToAmount     float64
}
