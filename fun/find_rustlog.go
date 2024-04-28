package fun

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	F.Register(&Cmd{
		Name: "find",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")

			if args[0] != "`find" {
				return
			}

			if len(args) < 1 {
				_, err = Say(m.RoomID, "must specify match pattern", m.ID)
				return
			}

			pattern := strings.Join(args[1:], " ")
			rows, err := Clickhouse.Query(context.Background(), "SELECT raw FROM message WHERE channel_id=? AND match(raw, ?) ORDER BY timestamp DESC", m.RoomID, pattern)
			if err != nil {
				return
			}
			defer rows.Close()

			tt := time.Now().UnixMilli()
			f, err := os.Create(fmt.Sprintf("/var/www/fi.supa.sh/dump/%v.txt", tt))
			if err != nil {
				return
			}

			fmt.Fprintf(f, "@badges=staff/1;color=#FFFFFF;display-name=System;first-msg=0;flags=;id=0;mod=0;returning-chatter=0;room-id=%s;subscriber=0;tmi-sent-ts=%v;turbo=0;user-id=1;user-type=staff :system!system@system.tmi.twitch.tv PRIVMSG #%s :logs.supa.codes dump - query: `%s`\n", m.RoomID, tt, m.Channel, pattern)

			var len int
			var raw string
			for rows.Next() {
				rows.Scan(&raw)
				fmt.Fprintln(f, raw)
				len = len + 1
			}

			if err := f.Close(); err != nil {
				return err
			}

			if err := rows.Err(); err != nil {
				return err
			}

			_, err = Say(m.RoomID, fmt.Sprintf("found %v messages: https://logs.raccatta.cc/?url=%s?raw=1&limit=99999", len, url.QueryEscape(fmt.Sprintf("https://fi.supa.sh/dump/%v.txt", tt))), m.ID)
			return
		},
	})
}
