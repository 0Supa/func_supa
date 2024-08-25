package utils

import (
	"slices"
	"strings"

	"github.com/0supa/func_supa/config"
)

var QuoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func IsPrivileged(userID string) bool {
	return slices.Contains(config.Meta.PrivilegedUsers, userID)
}

func StringPtr(str string) *string {
	return &str
}
