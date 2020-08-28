package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var hashes []string

func getRandomHash() string {
	randomSource := rand.NewSource(int64(rand.Uint64()))
	random := rand.New(randomSource)
	rIndex := random.Intn(len(hashes))

	return hashes[rIndex]
}

func readLinks() []byte {
	content, err := ioutil.ReadFile("links.txt")
	if err != nil {
		log.Fatal(err)
	}

	return content
}

func newCustomReadTargeter() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "GET"
		tgt.URL = "http://localhost:8080/" + getRandomHash()

		return nil
	}
}

func main() {
	hashes = strings.Split(strings.ReplaceAll(string(readLinks()), "'", ""), ", ")

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
	//fmt.Println(metrics.Errors)
	fmt.Println(fmt.Sprintf("--- "))
}
