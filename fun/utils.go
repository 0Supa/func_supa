package fun

import (
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/0supa/func_supa/config"
)

var QuoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

var GenericUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0"

var apiClient = http.Client{Timeout: time.Minute}

func IsPrivileged(userID string) bool {
	return slices.Contains(config.Meta.PrivilegedUsers, userID)
}
