package api

import (
	"net/http"
	"time"
)

const GenericUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:104.0) Gecko/20100101 Firefox/104.0"

var Generic = http.Client{Timeout: time.Minute}
