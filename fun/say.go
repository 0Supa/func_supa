package fun

import (
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	F.Register(&Cmd{
		Name: "say",
		Handler: func(m twitch.PrivateMessage) (err error) {
			if m.User.ID != "675052240" { // 8supa
				return
			}

			args := strings.Split(m.Message, " ")
			if args[0] != "`say" || len(args) < 3 {
				return
			}

			user, err := GetUser(args[1], "")
			if err != nil {
				_, err = Say(m.RoomID, "user not found: "+err.Error(), m.ID)
				return
			}

			res, err := Say(user.ID, strings.Join(args[2:], " "), "")
			if err != nil {
				_, err = Say(m.RoomID, "failed sending message: "+err.Error(), m.ID)
				return
			}

			_, err = Say(m.RoomID, res.Data.Mutation.Message.ID, m.ID)
			return
		},
	})
}
