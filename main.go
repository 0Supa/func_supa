package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/0supa/func_supa/config"
)

type LogChannels struct {
	Channels []struct {
		Name   string `json:"name"`
		UserID string `json:"userID"`
	} `json:"channels"`
}

type LiveChannels []struct {
	Name    string `json:"login"`
	UserID  string `json:"uid"`
	Viewers int    `json:"viewers"`
	Type    string `json:"type"`
}

type JoinPayload struct {
	Channels []string `json:"channels"`
}

var wg = &sync.WaitGroup{}

func main() {
	wg.Add(1)

	httpClient := http.Client{Timeout: time.Minute}
	go func() {
		for range time.Tick(time.Minute * 10) {
			res, err := httpClient.Get("https://logs.supa.codes/channels")
			if err != nil {
				log.Println(err)
				return
			}

			rustlog := LogChannels{}
			err = json.NewDecoder(res.Body).Decode(&rustlog)
			if err != nil {
				log.Println(err)
				return
			}
			res.Body.Close()

			ignored := []string{}
			for _, ch := range rustlog.Channels {
				ignored = append(ignored, ch.UserID)
			}

			res, err = httpClient.Get("https://tv.supa.sh/tags/ro")
			if err != nil {
				log.Println(err)
				return
			}

			liveChannels := LiveChannels{}
			err = json.NewDecoder(res.Body).Decode(&liveChannels)
			if err != nil {
				log.Println(err)
				return
			}
			res.Body.Close()

			var resMsg strings.Builder
			resMsg.WriteString("logs.supa.codes -> joining new channels:")

			joinPayload := JoinPayload{}
			for _, ch := range liveChannels {
				if ch.Viewers < 2 || slices.Contains(ignored, ch.UserID) {
					continue
				}
				resMsg.WriteString("\n@" + ch.Name)
				joinPayload.Channels = append(joinPayload.Channels, ch.UserID)
			}
			if len(joinPayload.Channels) == 0 {
				return
			}

			body, err := json.Marshal(joinPayload)
			if err != nil {
				return
			}

			req, _ := http.NewRequest(
				"POST", "https://logs.supa.codes/admin/channels",
				bytes.NewBuffer(body),
			)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Api-Key", config.Auth.Rustlog.Key)

			res, err = httpClient.Do(req)
			if err != nil {
				log.Println(err)
				return
			}

			if res.StatusCode != http.StatusOK {
				log.Println("failed joining new rustlog channels")

				b, err := io.ReadAll(res.Body)
				if err != nil {
					return
				}
				log.Println(string(b))
				return
			}
			res.Body.Close()

			// fun.Say("675052240", resMsg.String(), "")
		}
	}()

	wg.Wait()
}
