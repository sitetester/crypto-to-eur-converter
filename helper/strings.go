package helper

import (
	"log"
	"strconv"
)

func ToFloat(rate string) float64 {
	if rate == "" {
		return 0
	}

	float, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		log.Fatal(err)
	}

	return float
}
