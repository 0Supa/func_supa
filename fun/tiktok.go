package fun

import (
	"errors"
	"net/http"
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
			tries := 0
		retry:
			link := strings.Replace(links.FindString(m.Message), "/reels/", "/reel/", 1)
			if link == "" {
				return
			}

			out, err := exec.Command("yt-dlp",
				"-S", "vcodec:h264",
				"--get-url", link,
			).Output()
			if err != nil {
				return err
			}

			res, err := apiClient.Get(strings.TrimSuffix(string(out), "\n"))
			if err != nil {
				return err
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				if tries += 1; tries > 1 {
					return errors.New(res.Status + ": " + res.Request.RequestURI)
				}

				err := exec.Command("yt-dlp", "--update-to", "master").Run()
				if err != nil {
					return err
				}
				goto retry
			}

			upload, err := UploadFile(res.Body, "res.mp4", res.Header.Get("Content-Type"))
			if err != nil {
				return
			}

			_, err = Say(m.RoomID, "mirror: "+upload.Link+upload.Ext, m.ID)
			return err
		},
	})
}
