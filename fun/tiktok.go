package fun

import (
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	regexp "github.com/wasilibs/go-re2"
)

func init() {
	links := regexp.MustCompile(`(?i)\S*tiktok\.com\/\S+|\S*instagram\.com\/(reels?|p)\/\S+`)
	parentDir := "/var/www/fi.supa.sh/tiktok"

	F.Register(&Cmd{
		Name: "tiktok",
		Handler: func(m twitch.PrivateMessage) (err error) {
			link := links.FindString(m.Message)
			if link == "" {
				return
			}

			cmd := exec.Command("yt-dlp",
				"-S", "vcodec:h264",
				"--min-filesize", "50k",
				"--max-filesize", "100M",
				"--embed-metadata",
				"-P", fmt.Sprintf("%s/%s", parentDir, m.User.Name),
				"-o", fmt.Sprintf("%v.%%(ext)s", time.Now().Unix()),
				"--restrict-filenames",
				"-q", "--exec", "echo {}",
				link,
			)
			out, err := cmd.Output()
			if err != nil {
				if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
					return
				}
				return err
			}

			fileName := filepath.Base(strings.TrimSuffix(string(out), "\n"))

			_, err = Say(
				m.RoomID,
				fmt.Sprintf("mirror: https://fi.supa.sh/tiktok/%s/%s", m.User.Name, url.PathEscape(fileName)),
				m.ID,
			)

			return err
		},
	})
}
