package fun

import "strings"

var QuoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
