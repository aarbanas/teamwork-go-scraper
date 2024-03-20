package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	UserId  string `json:"userId"`
	ApiKey  string `json:"apiKey"`
	Url     string `json:"url"`
	UserIds []string
}

func LoadConfiguration(file string) Config {
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

	isValidConfig, err := validateConfiguration(config)
	if !isValidConfig && err != nil {
		fmt.Printf("Invalid configuration: %s\n", err)
		os.Exit(1)
	}

	return config
}

func validateConfiguration(configuration Config) (bool, error) {
	if configuration.ApiKey == "" {
		return false, errors.New("Missing ApiKey in configuration")
	}

	if configuration.Url == "" {
		return false, errors.New("Missing Url in configuration")
	}

	if configuration.UserId == "" {
		return false, errors.New("Missing UserId in configuration")
	}

	return true, nil
}
