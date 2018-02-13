package haynesbot

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
)

var (
	Token     string
	TestToken string
	BotPrefix string
	ManagedGuilds []string
	test      bool

	config *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
	TestToken string `json:"TestToken"`
	TestPrefix string `json:"TestPrefix"`
	ManagedGuilds []string `json:"ManagedGuilds"`
}

func ReadConfig() error {
	fmt.Println("Reading from config file...")

	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("Unable to get config file location")
	}
	file, err := ioutil.ReadFile(path.Join(path.Dir(filename), "../config.json"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if test {
		fmt.Println("Running test version...")

		config.Token = config.TestToken
		config.BotPrefix = config.TestPrefix
	}

	TestToken = config.TestToken
	Token = config.Token
	BotPrefix = config.BotPrefix
	ManagedGuilds = config.ManagedGuilds

	return nil
}

func init() {
	flag.BoolVar(&test, "t", false, "Run for testing")
	flag.Parse()
}
