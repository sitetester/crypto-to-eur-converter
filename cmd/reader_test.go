package main

import "testing"

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		// true
		{"BTC", true},
		{"btc", true},
		{"ETH", true},
		{"Solana", true},
		// false
		{"", false},
		{"123", false},
		{"BTC123", false},
		{"BTC-USD", false},
		{"BTC_USD", false},
		{"123BTC", false},
	}

	for _, tt := range tests {
		got := isAlpha(tt.input)
		if got != tt.want {
			t.Errorf("isAlpha(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestExtractCurrencyAndAmount(t *testing.T) {
	tests := []struct {
		name         string
		input        map[string]interface{}
		wantCurrency string
		wantAmount   float64
		wantErr      bool
	}{
		{
			name:         "BTC with string amount",
			input:        map[string]interface{}{"currency": "BTC", "amount": "0.005343150"},
			wantCurrency: "BTC",
			wantAmount:   0.005343150,
			wantErr:      false,
		},
		{
			name:         "BTC with string amount (in reverse order)",
			input:        map[string]interface{}{"amount": "0.005343150", "currency": "BTC"},
			wantCurrency: "BTC",
			wantAmount:   0.005343150,
			wantErr:      false,
		},
		{
			name:         "ETH with float amount",
			input:        map[string]interface{}{"curr_iso": "ETH", "price": 10.0},
			wantCurrency: "ETH",
			wantAmount:   10.0,
			wantErr:      false,
		},
		{
			name:         "SOL with decimal",
			input:        map[string]interface{}{"currency_iso": "SOL", "amount": 123.45},
			wantCurrency: "SOL",
			wantAmount:   123.45,
			wantErr:      false,
		},
		{
			name:         "lowercase currency",
			input:        map[string]interface{}{"currency": "btc", "amount": "100"},
			wantCurrency: "BTC",
			wantAmount:   100.0,
			wantErr:      false,
		},
		{
			name:         "Mixed case currency",
			input:        map[string]interface{}{"currency": "bTc", "amount": "100"},
			wantCurrency: "BTC",
			wantAmount:   100.0,
			wantErr:      false,
		},
		{
			name:         "4-letter currency",
			input:        map[string]interface{}{"currency": "USDT", "amount": 1000.0},
			wantCurrency: "USDT",
			wantAmount:   1000.0,
			wantErr:      false,
		},
		{
			name:    "no currency",
			input:   map[string]interface{}{"amount": 123, "price": 456},
			wantErr: true,
		},
		{
			name:    "no amount",
			input:   map[string]interface{}{"currency": "BTC", "coin": "ETH"},
			wantErr: true,
		},
		{
			name:    "currency too short",
			input:   map[string]interface{}{"currency": "BT", "amount": 100},
			wantErr: true,
		},
		{
			name:    "currency too long",
			input:   map[string]interface{}{"currency": "BITCOIN", "amount": 100},
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   map[string]interface{}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currency, amount, err := extractCurrencyAndAmount(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if currency != tt.wantCurrency {
				t.Errorf("currency = %q, want %q", currency, tt.wantCurrency)
			}

			if amount != tt.wantAmount {
				t.Errorf("amount = %f, want %f", amount, tt.wantAmount)
			}
		})
	}
}
