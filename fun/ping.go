package fun

import (
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	F.Register(&Cmd{
		Name: "ping",
		Handler: func(m twitch.PrivateMessage) (err error) {
			if m.Message != "`ping" {
				return
			}

			_, err = Say(m.RoomID, "pong!", m.ID)
			return
		},
	})
}
