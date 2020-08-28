package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v5"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func getRandomLink() string {
	return gofakeit.URL()
}

func newCustomReadTargeter() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		type linkT struct {
			Link string `json:"link"`
		}

		tgt.Method = "POST"
		tgt.URL = "http://localhost:8080/api/v1/link"
		tgt.Body, _ = json.Marshal(&linkT{
			Link: getRandomLink(),
		})

		return nil
	}
}

func main() {
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 30 * time.Second

	targeter := newCustomReadTargeter()
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Latencies: \n - 99: %v \n - 95: %v \n - 50: %v\n Total: %v\n Max: %v\n Mean: %v\n",
		metrics.Latencies.P99,
		metrics.Latencies.P95,
		metrics.Latencies.P50,
		metrics.Latencies.Total,
		metrics.Latencies.Max,
		metrics.Latencies.Mean,
	)

	fmt.Println(metrics.StatusCodes)
	fmt.Println(metrics.Errors)
	fmt.Println(fmt.Sprintf("--- "))
}
