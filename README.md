It's a simple Go program that takes input from different endpoints & returns the unified JSON described below
with amounts converted to EUR.

Currently latest conversion rates to EUR are provided via CryptoCompare, this can be replaced with another provider
by implementing `type RatesProvider interface`

Output structure:

```
currency String // Currency of the amount (EUR)
amount String // Amount converted to EUR
```

Example 1
"endpoint-btc" returns JSON:

```
{
 "currency": "BTC",
 "amount": "0.005343150"
}
```

Expected output:

```
{
 "endpoint": "endpoint-btc",
 "currency": "EUR",
 "amount": "534.31"
}
```

Example 2
"endpoint-eth" returns JSON:

```
{
 "curr_iso": "ETH",
 "price": 10
}
```

Expected output:

```
{
 "endpoint": "endpoint-eth",
 "currency": "EUR",
 "amount": "15000.00"
}
```

-----

## Setup

### localhost

First start mock server  
`go run cmd/mockserver/main.go`
 
Then run the main program from project root by passing an endpoint like:  
go run ./cmd http://localhost:8080/endpoint-btc  
go run ./cmd http://localhost:8080/endpoint-eth  
go run ./cmd http://localhost:8080/endpoint-sol  

### production

go run ./cmd https://your-real-api.com/endpoint-btc

**Tests**

- go test ./... -v
- go test -tags=integration ./... -v