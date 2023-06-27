package main

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

func Init(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening config file")
		return false
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	_ = decoder.Decode(&Cfg)
	return true
}
