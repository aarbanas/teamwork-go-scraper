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

	defer func() {
		closeErr := configFile.Close()
		if closeErr != nil {
			fmt.Println(closeErr.Error())
		}
	}()

	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	decodeErr := jsonParser.Decode(&config)
	if decodeErr != nil {
		fmt.Println(decodeErr)
	}

	return config
}
