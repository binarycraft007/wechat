package main

import (
	"log"

	"github.com/binarycraft007/wechat/core"
)

func main() {
	core, err := core.New()
	if err != nil {
		log.Fatal(err)
	}
	_ = core
	//fmt.Println(core.Config.Origin)
	//fmt.Println(core.Config.Api.JsLogin)

	//core.Events.Emit("my_event", "this is my payload")
}
