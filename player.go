package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	IP                                   string
	Name                                 string
	Index                                int
	Bc                                   Blockchain
	BlockchainUnifier                    string
	CurrentAmount                        int
	Bet                                  int
	OutOfGame, Folded, Double, MakeSplit bool
	Cards                                []Card
	TimeJoined                           time.Time
	Database                             string
	Conn                                 net.Conn
	commands                             chan<- command
}

type command struct {
	player *Player
	args   []string
}

func (p *Player) Connect() error {
	if err := http.ListenAndServe(p.IP, nil); err != nil {
		return fmt.Errorf("could not serve at %s :%v", p.IP, err)
	}
	return nil
}

func (p *Player) CommunicateWithTheGame() ([]string, error) {
	var args []string
	for {
		msg, err := bufio.NewReader(p.Conn).ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("could not communicate: %v", err)
		}

		msg = strings.Trim(msg, "\r\n")

		args = strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/bet":
			p.commands <- command{
				player: p,
				args:   args,
			}
		case "/decide":
			p.commands <- command{
				player: p,
				args:   args,
			}
		case "/msg":
			p.commands <- command{
				player: p,
				args:   args,
			}
		case "/quit":
			p.commands <- command{
				player: p,
			}
		default:
			p.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
	return args, nil
}

func (p *Player) Nickname(args []string) {
	if len(args) < 2 {
		p.msg("nick is required. usage: /nick NAME")
		return
	}

	p.Name = args[1]
	p.msg(fmt.Sprintf("all right, I will call you %s", p.Name))
}

func (p *Player) Decide(args []string) string {
	if len(args) < 2 {
		p.msg("A choice is neccessary")
	}
	return args[1]
}

func (p *Player) Quit() {
	p.msg("Clearing Player #" + strconv.Itoa(p.Index))
}
func (p *Player) err(err error) {
	p.Conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (p *Player) msg(msg string) {
	p.Conn.Write([]byte("> " + msg + "\n"))
}

func (p *Player) SumInHands() int {
	var sum int
	for _, k := range p.Cards {
		sum += k.Value
	}
	return sum
}
