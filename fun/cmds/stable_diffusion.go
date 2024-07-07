package fun

import (
	"strings"

	. "github.com/0supa/func_supa/fun"
	. "github.com/0supa/func_supa/fun/api/cloudflare"
	api_kappa "github.com/0supa/func_supa/fun/api/kappa"
	. "github.com/0supa/func_supa/fun/api/twitch"
	"github.com/gempir/go-twitch-irc/v4"
)

func init() {
	Fun.Register(&Cmd{
		Name: "stable-diffusion",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")
			if (args[0] != "`sd") || len(args) < 2 {
				return
			}

			Say(m.RoomID, "ppHop", m.ID)

			body, err := StableDiffusionImage(strings.Join(args[1:], " "))
			if err != nil {
				return
			}

			upload, err := api_kappa.UploadFile(body, "image.png", "image/png")
			if err != nil {
				return
			}

			_, err = Say(m.RoomID, upload.Link, m.ID)
			return
		},
	})
}
