package fun

import (
	logger "log"
	"os"
	"slices"

	"github.com/0supa/func_supa/config"
	api_twitch "github.com/0supa/func_supa/fun/api/twitch"

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
		if m.User.ID == config.Auth.Twitch.GQL.UserID {
			return
		}

		for _, cmd := range Fun.Cmds {
			go func(cmd Cmd) {
				channels := config.Meta.Functions[cmd.Name].Channels
				if len(channels) > 0 && !slices.Contains(channels, m.RoomID) {
					return
				}

				err := cmd.Handler(m)
				if err != nil {
					api_twitch.Say(m.RoomID, "🚫 "+err.Error(), m.ID)
					Log.Printf("%v:\n%v\n", m, err)
				}
			}(cmd)
		}
	})

	for _, userID := range config.Meta.Channels {
		u, err := api_twitch.GetUser("", userID)
		if err != nil {
			Log.Println("failed getting user", err)
		}
		Client.Join(u.Login)
	}

	go func() {
		err := Client.Connect()
		if err != nil {
			panic(err)
		}
	}()
}
