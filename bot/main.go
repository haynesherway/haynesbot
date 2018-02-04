package main

import (
	"fmt"
	"github.com/haynesherway/haynesbot"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})

	return
}
