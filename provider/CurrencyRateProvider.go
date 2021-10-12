package provider

type EuroOutput struct {
	Endpoint string
	Currency string
	Amount   float64
}

type CurrencyRate struct {
	Currency string
	Amount   string
}

type CurrencyRateProvider interface {
	ProvideRates(url string) []CurrencyRate
}
