package fun

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os/exec"
	"strings"

	"github.com/0supa/func_supa/fun"
	"github.com/gempir/go-twitch-irc/v4"
	regexp "github.com/wasilibs/go-re2"
)

type Upload struct {
	ID       string `json:"id"`
	Ext      string `json:"ext"`
	Type     string `json:"type"`
	Checksum string `json:"checksum"`
	Key      string `json:"key"`
	Link     string `json:"link"`
	Delete   string `json:"delete"`
}

func init() {
	r := regexp.MustCompile(`(?i)\S*tiktok\.com\/\S+|\S*instagram\.com\/(reels?|p)\/\S+`)
	fun.RegisterFun(fun.Fun{
		Name: "tiktok",
		Handler: func(m twitch.PrivateMessage) (err error) {
			link := strings.Replace(r.FindString(m.Message), "/reels/", "/reel/", 1)
			if link == "" {
				return
			}

			out, err := exec.Command("yt-dlp",
				"-S", "vcodec:h264",
				"--get-url", link,
			).Output()
			if err != nil {
				return
			}

			client := http.Client{}
			res, err := client.Get(strings.TrimSuffix(string(out), "\n"))
			if err != nil {
				return
			}
			defer res.Body.Close()

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", `form-data; name="file"; filename="res.mp4"`)
			h.Set("Content-Type", res.Header.Get("Content-Type"))

			part, err := writer.CreatePart(h)
			if _, err := io.Copy(part, res.Body); err != nil {
				return err
			}
			writer.Close()

			res, err = client.Post("https://kappa.lol/api/upload?skip-cd=true", writer.FormDataContentType(), body)
			if err != nil {
				return
			}

			buf, err := io.ReadAll(res.Body)
			var upload Upload
			if err := json.Unmarshal(buf, &upload); err != nil {
				return err
			}

			if res.StatusCode != http.StatusOK {
				return errors.New(string(buf))
			}

			_, err = fun.Say(m.RoomID, "mirror: "+upload.Link+upload.Ext, m.ID)
			return
		},
	})
}
