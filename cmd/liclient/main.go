package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"ireul.com/com"

	"ireul.com/yaml"
)

// Config client config
type Config struct {
	Host  string
	Token string
}

func main() {
	if len(os.Args) != 3 {
		println("Usage: liclient [NAME] [URL]")
		os.Exit(1)
		return
	}
	home := os.Getenv("HOME")
	if len(home) == 0 {
		panic("$HOME not set")
	}
	file := path.Join(home, ".config", "liclient.yaml")

	rc, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	c := Config{}

	if err = yaml.Unmarshal(rc, &c); err != nil {
		panic(err)
	}

	name, url := os.Args[1], os.Args[2]

	err = com.HttpPostJSON(http.DefaultClient, "https://"+path.Join(c.Host, "create"), map[string]string{
		"token": c.Token,
		"url":   url,
		"name":  name,
	}, nil)
	if err != nil {
		println(err.Error())
	}
}
