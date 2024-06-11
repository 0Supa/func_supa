package fun

import (
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/0supa/func_supa/fun"
	. "github.com/0supa/func_supa/fun/api/twitch"
	"github.com/gempir/go-twitch-irc/v4"
	regexp "github.com/wasilibs/go-re2"
)

var cooldown = map[string]struct{}{}

func init() {
	links := regexp.MustCompile(`(?i)\S*tiktok\.com\/\S+|\S*instagram\.com\/(reels?|p)\/\S+`)
	parentDir := "/var/www/fi.supa.sh/tiktok"

	Fun.Register(&Cmd{
		Name: "tiktok",
		Handler: func(m twitch.PrivateMessage) (err error) {
			link := strings.Replace(links.FindString(m.Message), "/reels/", "/reel/", 1)
			if link == "" {
				return
			}

			if _, found := cooldown[m.User.ID]; found {
				return
			}

			cooldown[m.User.ID] = struct{}{}

			defer func() {
				time.Sleep(10 * time.Second)
				delete(cooldown, m.User.ID)
			}()

			cmd := exec.Command("./bin/yt-dlp",
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
					fmt.Printf("yt-dlp error: %s:\n%s\n", err.Error(), exit.Stderr)
					return nil
				}
				return err
			}

			fileName := filepath.Base(strings.TrimSuffix(string(out), "\n"))

			if fileName == "." {
				return
			}

			_, err = Say(
				m.RoomID,
				fmt.Sprintf("mirror: https://fi.supa.sh/tiktok/%s/%s", m.User.Name, url.PathEscape(fileName)),
				m.ID,
			)

			return err
		},
	})
}
