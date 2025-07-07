package utils

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Whatsapp struct {
		WhiteListedChats []string `yaml:"WhiteListedChats"`
		MessageReceiver string `yaml:"MessageReceiver"`
	} `yaml:"Whatsapp"`

	AI struct {
		Controls struct {
			PollingInterval int `yaml:"PollingInterval"` // in seconds
			WorkerCount     int `yaml:"WorkerCount"`     // number of AI workers
		} `yaml:"Controls"`

		BenchmarkMessage string `yaml:"BenchmarkMessage"`
	} `yaml:"AI"`
}

func ReadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Error parsing config file: %v", err)
		return nil, err
	}

	return &config, nil
}



