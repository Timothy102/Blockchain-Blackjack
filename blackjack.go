package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

// 4 decki kart
var deck [52 * 4]bool

var farbe = []string{"srce", "kara", "križ", "pik"}

// hiša in izvajalec
type Game struct {
	Players      []*Player
	CurrentRound int
	Deck         []bool
	CardsInFront []Card
	Score        int
	SmallestBet  int
	MaxPlayers   int
}

func NewGame() *Game {
	return &Game{[]*Player{}, 0, make([]bool, 400), make([]Card, 100), 0, 10, 10}
}

func (g *Game) SplitTheCards() {
	for _, p := range g.Players {

		// dodeliš naključne dve karte vsem igralcem
		p.Cards[0] = g.getCard()
		p.Cards[1] = g.getCard()

		// kart s temi indexi ni več v decku
		g.Deck[p.Cards[0].Index] = true
		g.Deck[p.Cards[1].Index] = true
	}
}

func (g *Game) PlayRound() error {
	// hiša dva dve karte na začetku ven
	if g.CurrentRound == 0 || g.CurrentRound == 1 {
		card := g.getCard()
		g.CardsInFront = append(g.CardsInFront, card)
		g.Score += card.Value
	}

	// vsakega igralca obdelaj
	for i, p := range g.Players {
		for !p.Folded && !p.Double {
			args, err := p.CommunicateWithTheGame()
			if err != nil {
				return fmt.Errorf("Something wrong with Player's #%d : %v", i, err)
			}
			choice := p.Decide(args)
			if choice == "call" {
				p.Cards = append(p.Cards, g.getCard())
				if p.SumInHands() > 21 {
					p.OutOfGame = true
				}
			} else if choice == "fold" {
				p.Folded = true
			} else if choice == "double" {
				p.Double = true
				p.Cards = append(p.Cards, g.getCard())
				p.Bet *= 2
			} else {
				log.Fatalf("You entered an invalid choice: %s", choice)
			}
		}
	}
	return nil
}

func (g *Game) TransactToBlockchain() {
	var transactions []*Transaction
	for _, p := range g.Players {
		transactions = append(transactions, NewUTXOTransaction(p.BlockchainUnifier, "host", 1, &p.Bc))
	}
	cbtx := NewCoinbaseTX("s", genesisCoinbaseData)
	genesis := NewGenesisBlock(cbtx)

	// make all transactions by the end of the game as included in one block
	block := NewBlock(transactions, genesis.Hash)
	pow := NewProofOfWork(block)

	if !pow.Validate() {
		log.Fatalf("Hash invalid")
	}
}

func (g *Game) EveryoneFolded() bool {
	for _, p := range g.Players {
		if p.Folded != true {
			return false
		}
	}
	return true
}

func (g *Game) getCard() Card {
	cardIndex := rand.Intn(len(g.Deck))
	value := cardIndex % 14
	if value == 12 || value == 11 || value == 10 {
		value = 10
	} else if value == 13 {
		value = 11
	}

	farba := cardIndex / 14 % 4
	return Card{cardIndex, value, farbe[farba]}
}

type Card struct {
	Index int
	Value int
	Farba string
}

func (g *Game) FinalAct() {
	for {
		card := g.getCard()
		if g.Score+card.Value > 21 {
			// hiša zgubi, vsi dobijo svoje
			for _, p := range g.Players {
				p.CurrentAmount += p.Bet
			}
		} else if g.Score >= 17 {
			for _, p := range g.Players {
				if p.SumInHands() > g.Score {
					p.CurrentAmount += p.Bet
				} else if p.SumInHands() == g.Score {
					// nothing happens
				} else {
					p.CurrentAmount -= p.Bet
				}
			}
		} else {
			continue
		}
		g.CardsInFront = append(g.CardsInFront)

	}
}
func (g *Game) Play() error {
	g.SplitTheCards()
	for {
		if g.EveryoneFolded() {
			// hiša se odpre, da vidmo kdo zmaga
			g.FinalAct()
		}
		if err := g.PlayRound(); err != nil {
			return fmt.Errorf("Something wrong with round #%d : %v", g.CurrentRound, err)
		}
		g.CurrentRound += 1
	}
	return nil
}

func (g *Game) InitConnection(port string) error {
	if err := http.ListenAndServe(port, g.handler()); err != nil {
		return fmt.Errorf("could not serve at %s :%v", port, err)
	}
	return nil
}

func (g *Game) handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/blackjack", g.doubleHandler)
	return r
}

func (g *Game) doubleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to the fuckin party!!! <h1>")
	tmpl := template.Must(template.ParseFiles("forms.html"))

	if r.Method == "GET" {
		tmpl.Execute(w, nil)
		return
	}

	count := len(g.Players)
	if count < g.MaxPlayers {
		nickname := r.FormValue("nickname")
		g.Players = append(g.Players, &Player{Name: nickname, TimeJoined: time.Unix(0, 666666)})
		fmt.Print(nickname + "has joined the game! :) ")
	} else {
		log.Printf("Number of players exceeds maximum %v", count)
	}

	tmpl.Execute(w, struct{ Success bool }{true})
}
