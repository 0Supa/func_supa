package fun

import (
	"github.com/gempir/go-twitch-irc/v4"
)

type Cmd struct {
	Name    string
	Handler func(m twitch.PrivateMessage) error
}

type Fun struct {
	Cmds []Cmd
}

func (f *Fun) Register(c *Cmd) {
	f.Cmds = append(f.Cmds, *c)
}

var F = Fun{}
