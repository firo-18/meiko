package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Token    string `json:"token"`
	ClientID string `json:"clientID"`
	GuildID  string `json:"guildID"`
}

func Load(filename string) *Configuration {
	config := Configuration{}
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("read-file:", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln("unmarshal:", err)
	}

	return &config
}
