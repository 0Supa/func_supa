package fun

import (
	logger "log"
	"os"
	"slices"

	"github.com/0supa/func_supa/config"

	"github.com/gempir/go-twitch-irc/v4"
)

var Log = logger.New(os.Stdout, "TMI ", logger.LstdFlags)
var Client = twitch.NewAnonymousClient()

func init() {
	Client.OnConnect(func() {
		Log.Println("connected")
	})

	Client.OnSelfJoinMessage(func(m twitch.UserJoinMessage) {
		Log.Println("joined", m.Channel)
	})

	Client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		for _, cmd := range F.Cmds {
			go func(cmd Cmd) {
				channels := config.Meta.Functions[cmd.Name].Channels
				if len(channels) > 0 && !slices.Contains(channels, m.RoomID) {
					return
				}

				err := cmd.Handler(m)
				if err != nil {
					Say(m.RoomID, "🚫 "+err.Error(), m.ID)
					Log.Printf("%v:\n%v\n", m, err)
				}
			}(cmd)
		}
	})

	for _, userID := range config.Meta.Channels {
		u, err := GetUser("", userID)
		if err != nil {
			Log.Println("failed getting user", err)
		}
		Client.Join(u.Login)
	}

	err := Client.Connect()
	if err != nil {
		panic(err)
	}
}
