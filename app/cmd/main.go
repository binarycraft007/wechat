package main

import (
	"fmt"
	"log"

	"github.com/binarycraft007/wechat/core"
)

func main() {
	core, err := core.New()
	if err != nil {
		log.Fatal(err)
	}

	err = core.GetUUID()
	if err != nil {
		log.Fatal(err)
	}

	err = core.PreLogin()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(core.SessionData.UUID)
	core.Events.Emit("my_event", "this is my payload")
}
