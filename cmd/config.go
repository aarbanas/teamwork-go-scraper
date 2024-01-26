package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	UserId string `json:"userId"`
	ApiKey string `json:"apiKey"`
	Url    string `json:"url"`
}

func loadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}
