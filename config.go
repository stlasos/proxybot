package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type AppConf struct {
	ApiUrl     string `json:"tg_api_url"`
	Token      string `json:"tg_token"`
	ConnectUrl string `json:"proxy_connect_url"`
}

func NewAppConf() *AppConf {
	conf := AppConf{}
	conf.Load()
	return &conf
}

func (c *AppConf) Load() {
	var configPath = "proxybot.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	confData, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Print("Error while loading config.")
		panic(err)
	}

	err = json.Unmarshal(confData, &c)

	if err != nil {
		log.Print("Config parse error.")
		panic(err)
	}
}
