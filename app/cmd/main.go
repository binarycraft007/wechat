package main

import (
	"log"

	"os"

	"github.com/binarycraft007/wechat/core"
	"github.com/mdp/qrterminal/v3"
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

	qrterminal.Generate(core.QrCode, qrterminal.L, os.Stdout)
	log.Println(core.QrCodeUrl)

	err = core.PreLogin()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(core.Avatar)

	err = core.PreLogin()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(core.RedirectUri)

	err = core.Login()
	if err != nil {
		log.Fatal(err)
	}

	err = core.Init()
	if err != nil {
		log.Fatal(err)
	}

	//core.Events.Emit("my_event", "this is my payload")
}
