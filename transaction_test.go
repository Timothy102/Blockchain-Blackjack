package main

import (
	"strconv"
	"testing"
	"time"
)

func TestNewUTXOTransaction(t *testing.T) {
	address := "done"
	bc := NewBlockchain(address)
	tt := []struct {
		from, to string
		amount   int
	}{
		{from: "tim", to: "filip", amount: 5},
		{from: "filip", to: "tim", amount: 5},
		{from: "leon", to: "filip", amount: 5},
		{from: "tim", to: "leon", amount: 5},
	}

	var times []float64
	for n, tc := range tt {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			transaction := NewUTXOTransaction(tc.from, tc.to, tc.amount, bc)
			start := time.Now()
			bc.MineBlock([]*Transaction{transaction})

			processingTime := time.Since(start).Seconds()
			if processingTime > 10*time.Minute.Minutes() {
				t.Fatalf("Took too long man! :(")
			}
			times = append(times, processingTime)
		})
	}
}
