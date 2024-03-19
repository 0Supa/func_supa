package fun

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

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

func Say(channelID string, message string, parentID string) (response TwitchSendMsgResponse, err error) {
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
		return response, errors.New(*dropReason)
	}

	return
}
