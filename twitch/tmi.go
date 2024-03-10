package twitch

import (
	logger "log"
	"os"
	"slices"

	"github.com/0supa/func_supa/config"
	"github.com/0supa/func_supa/fun"

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
		for _, cmd := range fun.Funs {
			go func(cmd fun.Fun) {
				channels := config.Meta.Functions[cmd.Name].Channels
				if channels[0] != "*" && !slices.Contains(channels, m.RoomID) {
					return
				}

				err := cmd.Handler(m)
				if err != nil {
					// fun.Say(m.RoomID, "upaS a message was about to be sent but something broke", m.ID)
					Log.Printf("%v:\n%v\n", m, err)
				}
			}(cmd)
		}
	})

	for _, userID := range config.Meta.Channels {
		u, err := fun.GetUser("", userID)
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
