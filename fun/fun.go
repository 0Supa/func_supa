package fun

import "github.com/gempir/go-twitch-irc/v4"

type Fun struct {
	Name    string
	Handler func(m twitch.PrivateMessage) error
}

var Funs = []Fun{}

func RegisterFun(fun Fun) {
	Funs = append(Funs, fun)
}
