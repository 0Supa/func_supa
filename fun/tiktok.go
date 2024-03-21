package fun

import (
	"os/exec"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
	regexp "github.com/wasilibs/go-re2"
)

func init() {
	links := regexp.MustCompile(`(?i)\S*tiktok\.com\/\S+|\S*instagram\.com\/(reels?|p)\/\S+`)

	F.Register(&Cmd{
		Name: "tiktok",
		Handler: func(m twitch.PrivateMessage) (err error) {
			link := strings.Replace(links.FindString(m.Message), "/reels/", "/reel/", 1)
			if link == "" {
				return
			}

			cmd := exec.Command("yt-dlp",
				"-S", "vcodec:h264",
				link,
				"-o", "-",
			)
			out, err := cmd.StdoutPipe()
			if err != nil {
				return err
			}
			defer out.Close()

			if err = cmd.Start(); err != nil {
				return
			}

			// assuming it's always an mp4, not ideal
			upload, err := UploadFile(out, "res.mp4", "video/mp4")
			if err != nil {
				return
			}

			_, err = Say(m.RoomID, "mirror: "+upload.Link+upload.Ext, m.ID)
			return err
		},
	})
}
