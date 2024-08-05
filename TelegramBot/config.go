package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type config struct {
	Bot struct {
		Token string `json:"token"`
	} `json:"bot"`
	Seeker struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"seeker"`
}

func GetConfiguration() (conf config, err error) {
	configFilename := os.Getenv("PATH_TO_CONFIG")
	if configFilename == "" {
		configFilename = "config.json"
	}

	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return config{}, errors.New("error: Couldn't read config file")
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return config{}, errors.New("error: Couldn't unmarshal config file")
	}
	return
}
