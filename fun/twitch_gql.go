package fun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/0supa/func_supa/config"
)

type TwitchGQLPayload struct {
	OperationName string `json:"operationName"`
	Query         string `json:"query"`
	Variables     any    `json:"variables"`
}

type TwitchGQLBaseResponse struct {
	Extensions struct {
		Duration      json.Number `json:"durationMilliseconds"`
		OperationName string      `json:"operationName"`
		RequestID     string      `json:"requestID"`
	} `json:"extensions"`
}

type TwitchUserResponse struct {
	*TwitchGQLBaseResponse
	Data struct {
		User TwitchUser `json:"user"`
	} `json:"data"`
}

type TwitchSendMsgResponse struct {
	*TwitchGQLBaseResponse
	Data struct {
		Mutation struct {
			DropReason *string `json:"dropReason"` // nullable hence pointer
			Message    struct {
				ID string `json:"id"`
			} `json:"message"`
		} `json:"sendChatMessage"`
	} `json:"data"`
}

type TwitchUser struct {
	ID          string `json:"id,omitempty"`
	Login       string `json:"login,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type Input struct {
	ChannelID string `json:"channelID"`
	Message   string `json:"message"`
	ParentID  string `json:"replyParentMessageID"`
}

type TwitchMsg struct {
	Input `json:"input"`
}

func GetUser(login string, id string) (user TwitchUser, err error) {
	response := TwitchUserResponse{}

	payload, err := json.Marshal(TwitchGQLPayload{
		OperationName: "User",
		Query:         "query User($login:String $id:ID) { user(login:$login id:$id) { id login displayName } }",
		Variables: TwitchUser{
			Login: login,
			ID:    id,
		},
	})
	if err != nil {
		return
	}

	req, _ := http.NewRequest("POST", "https://gql.twitch.tv/gql", bytes.NewBuffer(payload))
	req.Header.Set("Client-Id", config.Auth.Twitch.GQL.ClientID)

	res, err := apiClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return
	}

	user = response.Data.User
	return
}

var zeroWidthChar = "\U000E0000"

func Say(channelID string, message string, parentID string, ctx ...int) (response TwitchSendMsgResponse, err error) {
	if len(ctx) == 0 {
		ctx = append(ctx, 1)
	}

	og := message
	uploadMessage := func() (upload FileUpload) {
		rc := io.NopCloser(strings.NewReader(og))
		defer rc.Close()

		// TODO: someway handle err?
		upload, _ = UploadFile(rc, "msg.txt", "text/plain")
		return
	}

	if len(message) > 400 {
		message = message[:200] + "... " + uploadMessage().Link
	}

	payload, err := json.Marshal(TwitchGQLPayload{
		OperationName: "SendChatMessage",
		Query:         "mutation SendChatMessage($input: SendChatMessageInput!) {  sendChatMessage(input: $input) {  dropReason  message {  id  }  }  }",
		Variables: TwitchMsg{
			Input{
				ChannelID: channelID,
				Message:   message,
				ParentID:  parentID,
			},
		},
	})
	if err != nil {
		return
	}

	req, _ := http.NewRequest("POST", "https://gql.twitch.tv/gql", bytes.NewBuffer(payload))
	req.Header.Set("Client-Id", config.Auth.Twitch.GQL.ClientID)
	req.Header.Set("Authorization", "OAuth "+config.Auth.Twitch.GQL.Token)

	res, err := apiClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return
	}

	if dropReason := response.Data.Mutation.DropReason; dropReason != nil {
		if *dropReason == "RATE_LIMIT" || *dropReason == "MSG_DUPLICATE" {
			time.Sleep(time.Second)
			suf := " " + zeroWidthChar
			message, found := strings.CutSuffix(message, suf)
			if !found {
				message += suf
			}

			return Say(channelID, message, parentID, ctx...)
		}

		i := len(ctx) - 1
		if ctx[i] > 2 {
			return response, fmt.Errorf("message dropped after %v retries: %s", ctx[i], *dropReason)
		}

		return Say(channelID, fmt.Sprintf("(%s) failed to send reply: %s", *dropReason, uploadMessage().Link), parentID, append(ctx[:i], ctx[i]+1)...)
	}

	return
}
