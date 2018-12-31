package main

import (
	"encoding/json"
)

type Config struct {
	Username     string
	Password     string
	AccessToken  string
	RefreshToken string
	Cookie       string
}

func SaveConfig(c Config) ([]byte, error) {
	data, err := json.MarshalIndent(&c, "", "  ")
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func LoadConfig(data []byte) (Config, error) {
	c := Config{}
	err := json.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
