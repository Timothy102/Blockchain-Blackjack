package main

import (
	"fmt"
	"testing"
	"time"
)

var miningTimes []float64

func TestMine(t *testing.T) {
	bc := CreateBlockchain("e")
	transaction := NewUTXOTransaction("s", "ad", 10, bc)
	cbtx := NewCoinbaseTX("s", genesisCoinbaseData)
	genesis := NewGenesisBlock(cbtx)

	t.Run("nothing", func(t *testing.T) {
		block := NewBlock([]*Transaction{transaction}, genesis.Hash)
		pow := NewProofOfWork(block)

		start := time.Now()
		hashes, _ := pow.Run()
		if !pow.Validate() {
			t.Fatalf("Hash invalid")
		}

		miningTime := time.Since(start).Seconds()
		miningTimes = append(miningTimes, miningTime)
		fmt.Println(miningTime)

		hashRate := hashes / int(miningTime)
		fmt.Println("Hash rate stands at: ", hashRate)
	})
}

// hash rate: as the price of the coin is higher, more hashrate joins the network as less efficient miners can remain profitable due to fatter margins.
// minerji dobijo award
