package main

import (
	"errors"
	"io/ioutil"

	"ireul.com/web"
	"ireul.com/yaml"
)

// Config Config
type Config struct {
	DatabaseURL string `yaml:"database_url"`
	RedisURL    string `yaml:"redis_url"`
	Token       string `yaml:"token"`
	Port        int    `yaml:"port"`
	Env         string `yaml:"env"`
	Title       string `yaml:"title"`
}

// ParseConfigFile parse a config file and returns url
func ParseConfigFile(file string) (c Config, err error) {
	var s []byte
	if s, err = ioutil.ReadFile(file); err != nil {
		return
	}
	err = yaml.Unmarshal(s, &c)
	if err == nil {
		if c.Title == "" {
			c.Title = "Linkage"
		}
		if c.Env == "" {
			c.Env = web.DEV
		}
		if c.DatabaseURL == "" {
			err = errors.New("database url not set")
		} else if c.RedisURL == "" {
			err = errors.New("redis url not set")
		} else if c.Token == "" {
			err = errors.New("token not set")
		} else if c.Port == 0 {
			err = errors.New("port not set")
		}
	}
	return
}
