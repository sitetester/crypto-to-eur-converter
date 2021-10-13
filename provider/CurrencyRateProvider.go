package provider

type ExpectedOutput struct {
	Endpoint     string
	FromCurrency string
	FromAmount   float64
	ToCurrency   string
	ToAmount     float64
}

type CurrencyRate struct {
	Currency string
	Amount   string
}

type CurrencyRateProvider interface {
	ProvideRates(url string) []CurrencyRate
}
