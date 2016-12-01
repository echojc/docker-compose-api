package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type service struct {
	Path string
}

type Config struct {
	Binary   string
	Port     int16
	Services map[string]service
}

func NewConfig(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(bytes, &config)
	return &config, err
}
