package fun

import (
	"github.com/gempir/go-twitch-irc/v4"
)

type Cmd struct {
	Name    string
	Handler func(m twitch.PrivateMessage) error
}

type FunData struct {
	Cmds           []Cmd
	BlockedUserIDs []string
}

func (f *FunData) Register(c *Cmd) {
	f.Cmds = append(f.Cmds, *c)
}

var Fun = FunData{}
