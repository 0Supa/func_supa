package fun

import (
	"fmt"
	"time"

	. "github.com/0supa/func_supa/fun"
	. "github.com/0supa/func_supa/fun/api/twitch"
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	Fun.Register(&Cmd{
		Name: "ping",
		Handler: func(m twitch.PrivateMessage) (err error) {
			if m.Message != "`ping" {
				return
			}

			_, err = Say(m.RoomID, fmt.Sprintf("pong! %vms", time.Since(m.Time).Milliseconds()), m.ID)
			return
		},
	})
}
