package fun

import (
	"fmt"
	"runtime"
	"time"

	"github.com/0supa/func_supa/config"
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

			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)

			_, err = Say(m.RoomID, fmt.Sprintf("pong! %vms - %s (%vMiB) - up:%s - channels:%v",
				time.Since(m.Time).Milliseconds(),
				runtime.Version(), mem.Alloc/1024/1024,
				time.Since(InitTime),
				len(config.Meta.Channels)), m.ID)
			return
		},
	})
}
