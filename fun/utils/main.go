package utils

import (
	"strings"

	"github.com/0supa/func_supa/config"
)

var QuoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

var privilegedUsers = make(map[string]struct{})
var botUsers = make(map[string]struct{})

func SliceToSet(slice []string, set map[string]struct{}) {
	for _, el := range slice {
		set[el] = struct{}{}
	}
}

func init() {
	SliceToSet(config.Meta.PrivilegedUsers, privilegedUsers)
	SliceToSet(config.Meta.BotUsers, botUsers)
}

func IsPrivileged(userID string) bool {
	_, exists := privilegedUsers[userID]
	return exists
}

func IsBot(userID string) bool {
	_, exists := botUsers[userID]
	return exists
}

func StringPtr(str string) *string {
	return &str
}
