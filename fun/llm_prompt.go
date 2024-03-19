package fun

import (
	"fmt"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

func systemPrompt(m twitch.PrivateMessage) string {
	return fmt.Sprintf(`Current time: %s

You can use Markdown syntax for formatting your response, excluding tables and images.
Do NOT add opening or closing sentences.
You are talking to user '%s' in channel '%s'.

Prompt:`,
		// prompt appended with query
		time.Now().Format("2006-01-02 15:04:05"),
		m.User.Name, m.Channel)
}

func init() {
	model := "@cf/meta/llama-2-7b-chat-fp16"
	F.Register(&Cmd{
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
