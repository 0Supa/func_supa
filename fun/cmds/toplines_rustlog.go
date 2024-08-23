package fun

import (
	"context"
	"fmt"
	"io"
	"strings"

	. "github.com/0supa/func_supa/fun"
	logs_db "github.com/0supa/func_supa/fun/api/clickhouse_db"
	api_kappa "github.com/0supa/func_supa/fun/api/kappa"
	. "github.com/0supa/func_supa/fun/api/twitch"
	"github.com/0supa/func_supa/fun/utils"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/olekukonko/tablewriter"
)

func init() {
	Fun.Register(&Cmd{
		Name: "toplines",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")
			if !utils.IsPrivileged(m.User.ID) || len(args) < 2 || args[0] != "`toplines" {
				return
			}

			user, err := GetUser(args[1], "")
			if err != nil {
				Say(m.RoomID, "failed getting user: "+err.Error(), m.ID)
				return
			}

			if user.ID == "" {
				Say(m.RoomID, "user not found", m.ID)
				return
			}

			rows, err := logs_db.Clickhouse.Query(context.Background(), "SELECT channel_login, count() AS lines FROM rustlog_zonian.message_structured WHERE user_id=? GROUP BY channel_login ORDER BY lines DESC", user.ID)
			if err != nil {
				return
			}
			defer rows.Close()

			tableString := &strings.Builder{}
			tableString.WriteString("Top lines by @" + user.Login + "\n\n")
			table := tablewriter.NewWriter(tableString)
			table.SetHeader([]string{"Channel", "Lines"})

			var cLogin string
			var lineCount uint64
			for rows.Next() {
				if err := rows.Scan(&cLogin, &lineCount); err != nil {
					return err
				}
				table.Append([]string{cLogin, fmt.Sprintf("%v", lineCount)})
			}

			rc := io.NopCloser(strings.NewReader(tableString.String()))
			defer rc.Close()

			upload, err := api_kappa.UploadFile(rc, "dat.txt", "text/plain")

			_, err = Say(m.RoomID, upload.Link, m.ID)

			return err
		},
	})
}
