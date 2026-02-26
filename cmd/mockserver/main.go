package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func jsonResponse(w http.ResponseWriter, json string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}

func jsonHandler(json string) http.HandlerFunc {
	// Because it must match the 2nd argument to http.HandleFunc()
	return func(w http.ResponseWriter, _ *http.Request) {
		jsonResponse(w, json)
	}
}

func main() {
	endpoints := map[string]string{
		"/endpoint-btc": `{"currency": "BTC", "amount": "0.005343150"}`,
		"/endpoint-eth": `{"curr_iso": "ETH", "price": 10}`,
		"/endpoint-sol": `{"currency_iso": "SOL", "amount": 123.45}`,
	}

	for path, json := range endpoints {
		http.HandleFunc(path, jsonHandler(json))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Mock server running on :%s\n", port)
	fmt.Println("Available endpoints:")
	for path := range endpoints {
		fmt.Printf("  http://localhost:%s%s\n", port, path)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
