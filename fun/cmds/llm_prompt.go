package fun

import (
	"fmt"
	"strings"
	"time"

	. "github.com/0supa/func_supa/fun"
	. "github.com/0supa/func_supa/fun/api/cloudflare"
	. "github.com/0supa/func_supa/fun/api/twitch"
	"github.com/gempir/go-twitch-irc/v4"
)

func systemPrompt(m twitch.PrivateMessage) string {
	return fmt.Sprintf(`Know that:
- The current time is: %s
- You are in the channel '%s', talking to user '%s'.
- You can use Markdown syntax for formatting your response, excluding tables and images.

Rules you must follow:
- Do NOT add opening or closing sentences.
- Keep your response concise.

Q:`,
		// prompt appended with query
		time.Now().Format(time.RFC1123),
		m.Channel, m.User.Name)
}

func init() {
	model := "@cf/meta/llama-3.1-70b-instruct"
	Fun.Register(&Cmd{
		Name: "llm",
		Handler: func(m twitch.PrivateMessage) (err error) {
			args := strings.Split(m.Message, " ")
			if (args[0] != "`ask" && args[0] != "`llm") || len(args) < 2 {
				return
			}

			query := TextQuery{
				Stream: true,
				Messages: []QueryMessage{
					{
						Role:    "system",
						Content: systemPrompt(m),
					},
					{
						Role:    "user",
						Content: strings.Join(args[1:], " "),
					},
				},
			}

			c := make(chan Result)
			go TextGeneration(c, query, model)

			var builder strings.Builder
			for data := range c {
				if err := data.Error; err != nil {
					return err
				}

				builder.WriteString(data.Response)
			}
			res := builder.String()

			_, err = Say(m.RoomID, res, m.ID)
			return
		},
	})
}
