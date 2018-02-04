package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

var (
	Token     string
	BotPrefix string
	test		bool

	config *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading from config file...")

	file, err := ioutil.ReadFile("./config.json")
	if test {
		fmt.Println("Running test version...")
		file, err = ioutil.ReadFile("./testconfig.json")
	}
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}

func init() {
	flag.BoolVar(&test, "t", false, "Run for testing")
	flag.Parse()
}
