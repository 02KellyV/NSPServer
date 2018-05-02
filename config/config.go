package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port     int `json:"port"`
	Database struct {
		Driver string `json:"driver"`
		Host   string `json:"host"`
		Name   string `json:"name"`
		Port   int    `json:"port"`
		User   string `json:"user"`
		Pass   string `json:"pass"`
	} `json:"database"`
}

func LoadConfig() (Config, error) {
	f, err := os.Open("./config/config.json")
	conf := Config{}
	if err != nil {
		return conf, err
	}
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
