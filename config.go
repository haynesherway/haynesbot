package haynesbot

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"path"
	"runtime"
)

// Config values
var (
	Token       string
	TestToken   string
	BotPrefix   string
	UseImages   bool
	ImageServer string
	test        bool

	config *configStruct
)

type configStruct struct {
	Token         string `json:"Token"`
	BotPrefix     string `json:"BotPrefix"`
	TestToken     string `json:"TestToken"`
	TestPrefix    string `json:"TestPrefix"`
	Images        bool   `json:"Images"`
	ImageServer   string `json:"ImageServer"`
	GuildFile     string `json:"GuildSettings"`
	TestGuildFile string `json:"TestGuildSettings"`
}

// ReadConfig reads the config file and initializes values using those configs
func ReadConfig() error {
	flag.Parse()
	log.Println("Reading from config file...")

	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("Unable to get config file location")
	}
	file, err := ioutil.ReadFile(path.Join(path.Dir(filename), "../config.json"))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if test {
		log.Println("Running test version...")

		config.Token = config.TestToken
		config.BotPrefix = config.TestPrefix
		config.GuildFile = config.TestGuildFile
	}

	log.Println("Guild file: ", config.GuildFile)
	readGuildSettings(config.GuildFile)

	TestToken = config.TestToken
	Token = config.Token
	BotPrefix = config.BotPrefix
	UseImages = config.Images
	ImageServer = config.ImageServer

	return nil
}

func init() {
	flag.BoolVar(&test, "t", false, "Run for testing")
}
