package main

import "github.com/binarycraft007/wechat/core"

func main() {
	core := core.New()

	_ = core
	//fmt.Println(core.Config.Origin)
	//fmt.Println(core.Config.Api.JsLogin)

	//core.Events.Emit("my_event", "this is my payload")
}
