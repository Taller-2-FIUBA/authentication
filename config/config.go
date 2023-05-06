package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var Cfg Config

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Token struct {
		Expiration int64 `yaml:"expiration"`
	} `yaml:"token"`
	Firebase struct {
		Apikey string `yaml:"apikey"`
	} `yaml:"firebase"`
}

func Init() bool {
	f, err := os.Open("config/config.yml")
	if err != nil {
		fmt.Println("Error opening config file")
		return false
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	_ = decoder.Decode(&Cfg)
	if err != nil {
		fmt.Println("Error decoding config file")
		return false
	}
	return true
}
