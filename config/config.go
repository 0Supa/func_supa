package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var Meta = struct {
	Channels  []string `yaml:"channels"`
	Functions map[string]struct {
		Channels []string `yaml:"channels"`
	} `yaml:"functions"`
}{}

var Auth = struct {
	Twitch struct {
		GQL struct {
			ClientID string `yaml:"client_id"`
			Token    string `yaml:"token"`
		} `yaml:"gql"`
	} `yaml:"twitch"`
	Cloudflare struct {
		Key   string `yaml:"key"`
		AccID string `yaml:"account_id"`
	} `yaml:"cloudflare"`
}{}

func loadConfig(file string, y interface{}) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, y)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	loadConfig("meta.yml", &Meta)
	loadConfig("auth.yml", &Auth)
}
