package fun

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/0supa/func_supa/config"
)

type QueryMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TextQuery struct {
	Stream   bool           `json:"stream"`
	Messages []QueryMessage `json:"messages"`
}

type ImgPrompt struct {
	Prompt string `json:"prompt"`
	Steps  int    `json:"num_steps"`
}

type Result struct {
	Response string `json:"response"`
	Error    error
}

var baseURL = "https://api.cloudflare.com/client/v4/accounts/"

func StableDiffusionImage(prompt string) (io.ReadCloser, error) {
	model := "@cf/stabilityai/stable-diffusion-xl-base-1.0"

	payload, err := json.Marshal(ImgPrompt{
		Prompt: prompt,
		Steps:  20,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(
		"POST", baseURL+config.Auth.Cloudflare.AccID+"/ai/run/"+model,
		bytes.NewBuffer(payload),
	)
	req.Header.Set("Authorization", "Bearer "+config.Auth.Cloudflare.Key)

	res, err := (&http.Client{Timeout: 5 * time.Minute}).Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK || !strings.HasPrefix(res.Header.Get("content-type"), "image/") {
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(b))
	}

	return res.Body, nil
}

func TextGeneration(c chan Result, query TextQuery, model string) {
	payload, err := json.Marshal(query)
	if err != nil {
		c <- Result{Error: err}
		return
	}

	req, _ := http.NewRequest(
		"POST", baseURL+config.Auth.Cloudflare.AccID+"/ai/run/"+model,
		bytes.NewBuffer(payload),
	)
	req.Header.Set("Authorization", "Bearer "+config.Auth.Cloudflare.Key)

	res, err := apiClient.Do(req)
	if err != nil {
		c <- Result{Error: err}
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			c <- Result{Error: err}
			return
		}

		c <- Result{Error: errors.New(string(b))}
		return
	}

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		s := strings.TrimSpace(strings.TrimPrefix(scanner.Text(), "data: "))

		if s == "" {
			continue
		}
		if s == "[DONE]" {
			break
		}

		var res Result
		err := json.Unmarshal([]byte(s), &res)
		if err != nil {
			res.Error = err
		}
		c <- res
	}
	close(c)
}
