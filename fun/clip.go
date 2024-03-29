package fun

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

type clip struct {
	File    string `json:"file"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func init() {
	F.Register(&Cmd{
		Name: "clip",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")
			if args[0] != "`clip" && args[0] != "?clip" {
				return
			}

			var channel = m.Channel
			if len(args) >= 2 {
				channel = args[1]
			}

			req, _ := http.NewRequest("POST", "http://192.168.100.31:8989/clip/"+url.PathEscape(channel), nil)
			res, err := apiClient.Do(req)
			if err != nil {
				return
			}

			decoder := json.NewDecoder(res.Body)
			var c clip
			err = decoder.Decode(&c)
			if err != nil {
				return
			}

			if c.Error != "" {
				if c.Message != "" {
					_, err = Say(m.RoomID, c.Message, m.ID)
					return
				}

				_, err = Say(m.RoomID, c.Error, m.ID)
				return
			}

			_, err = Say(m.RoomID, "https://clips.supa.sh/"+url.PathEscape(c.File), m.ID)
			return
		},
	})
}
