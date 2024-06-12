package fun

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	. "github.com/0supa/func_supa/fun"
	logs_db "github.com/0supa/func_supa/fun/api/clickhouse_db"
	. "github.com/0supa/func_supa/fun/api/twitch"
	"github.com/dustin/go-humanize"
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	Fun.Register(&Cmd{
		Name: "global_rl",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")

			if args[0] != "`grl" && args[0] != "`grq" {
				return
			}

			target := m.User.ID
			if len(args) > 1 {
				x := args[1]

				u, err := GetUser(x, "")
				if err != nil {
					return err
				}
				target = u.ID
			}

			var count uint64
			err = logs_db.Clickhouse.QueryRow(context.Background(), "SELECT count() FROM message WHERE user_id = ?", target).Scan(&count)
			if err != nil {
				return
			}

			if count == 0 {
				_, err = Say(m.RoomID, "no messages from user", m.ID)
				return
			}

			offset := rand.Intn(int(count))

			var raw string
			err = logs_db.Clickhouse.QueryRow(context.Background(), "SELECT raw FROM message WHERE user_id = ? LIMIT 1 OFFSET ?", target, offset).Scan(raw)
			if err != nil {
				return
			}

			rMsg := twitch.ParseMessage(raw)

			switch msg := rMsg.(type) {
			case *twitch.PrivateMessage:
				_, err = Say(m.RoomID, fmt.Sprintf("[%s] #%s %s: %s", humanize.Time(msg.Time), msg.Channel, msg.User.Name, msg.Message), m.ID)
			default:
				_, err = Say(m.RoomID, raw, m.ID)
			}
			return
		},
	})
}
