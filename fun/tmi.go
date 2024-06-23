package fun

import (
	logger "log"
	"os"
	"slices"

	"github.com/0supa/func_supa/config"
	api_twitch "github.com/0supa/func_supa/fun/api/twitch"

	"github.com/gempir/go-twitch-irc/v4"
)

var Client = twitch.NewAnonymousClient()

var log = logger.New(os.Stdout, "TMI ", logger.LstdFlags)

func init() {
	user, err := api_twitch.GetSelf()
	if err != nil {
		logger.Panicln("failed getting current user", err)
	}

	if user.BlockedUsers != nil {
		for _, u := range *user.BlockedUsers {
			Fun.BlockedUserIDs = append(Fun.BlockedUserIDs, u.ID)
		}
		log.Printf("%v blocked Twitch users\n", len(Fun.BlockedUserIDs))
	}

	Client.OnConnect(func() {
		log.Println("connected")
	})

	Client.OnSelfJoinMessage(func(m twitch.UserJoinMessage) {
		log.Println("joined", m.Channel)
	})

	Client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		if m.User.ID == config.Auth.Twitch.GQL.UserID {
			return
		}

		if slices.Contains(Fun.BlockedUserIDs, m.User.ID) {
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
					log.Printf("[cmd error] %v: %v => %v\n", m.User.Name, m.Message, err)
					api_twitch.Say(m.RoomID, "🚫 "+err.Error(), m.ID)
				}
			}(cmd)
		}
	})

	for _, userID := range config.Meta.Channels {
		u, err := api_twitch.GetUser("", userID)
		if err != nil {
			log.Println("failed getting user", err)
		}
		Client.Join(u.Login)
	}

	go func() {
		err := Client.Connect()
		if err != nil {
			log.Panicln("failed connecting to Twitch IRC", err)
		}
	}()
}
