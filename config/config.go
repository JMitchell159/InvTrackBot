package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Token     string `json:"token"`
	BotPrefix string `json:"botPrefix"`
	DBUrl     string `json:"dbURL"`
}

func ReadConfig() (*Config, error) {
	fmt.Println("Reading config.json")
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return nil, err
	}

	fmt.Println("Unmarshalling config.json...")
	var cfg Config
	err = json.Unmarshal([]byte(data), &cfg)
	if err != nil {
		fmt.Println("Error unmarshaling config.json")
		return nil, err
	}

	return &cfg, nil
}
