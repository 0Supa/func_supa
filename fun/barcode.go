package fun

import (
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
	regexp "github.com/wasilibs/go-re2"
)

func init() {
	links := regexp.MustCompile(`(?i)\S*(i\.supa\.codes|kappa\.lol|gachi\.gay|femboy\.beauty)\/(\S+)`)

	F.Register(&Cmd{
		Name: "barcode",
		Handler: func(m twitch.PrivateMessage) (err error) {
			match := links.FindStringSubmatch(m.Message)
			if len(match) < 3 {
				return
			}

			id := match[2]
			fileURL := "https://kappa.lol/" + url.PathEscape(id)
			res, err := apiClient.Head(fileURL)
			if err != nil || res.StatusCode != http.StatusOK || !strings.HasPrefix(res.Header.Get("Content-Type"), "image/") {
				return
			}

			res, err = apiClient.Get(fileURL)
			if err != nil {
				return
			}

			cmd := exec.Command("zbarimg", "-q", "-")
			cmd.Stdin = res.Body
			out, err := cmd.Output()
			if err != nil {
				if _, ok := err.(*exec.ExitError); ok {
					return nil
				}
				return
			}

			for _, str := range strings.Split(string(out), "\n") {
				dat := strings.SplitN(str, ":", 2)
				if len(dat) < 2 || (!strings.HasPrefix(dat[0], "EAN") && !strings.HasPrefix(dat[0], "UPC")) {
					continue
				}
				_, err = Say(m.RoomID, str+" https://world.openfoodfacts.org/product/"+url.PathEscape(dat[1]), m.ID)
			}

			return
		},
	})
}
