package main

import (
	"flag"
	"log"
	"strconv"
)

func main() {

	var port, maxPlayers, smallestBet int
	flag.IntVar(&port, "port", 8080, "input port for establishing connection")
	flag.IntVar(&maxPlayers, "maxP", 10, "maximum number of players")
	flag.IntVar(&smallestBet, "bet", 50, "set the initial bet")

	flag.Parse()

	game := NewGame()
	game.MaxPlayers = maxPlayers
	game.SmallestBet = smallestBet

	if err := game.InitConnection(strconv.Itoa(port)); err != nil {
		log.Fatalf("could not init connection %d: %v", port, err)
	}
	if err := game.Play(); err != nil {
		log.Fatalf("Having trouble playing: %v", err)
	}
	game.TransactToBlockchain()
}

// fix the player code
// external IP: 172.17.0.1

// flagi: port, maxPlayers,
