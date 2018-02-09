package main

import (
	"fmt"
	"github.com/haynesherway/haynesbot"
)

func main() {
	err := haynesbot.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	haynesbot.Start()

	<-make(chan struct{})

	return
}
