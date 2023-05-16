package main

import (
	"fmt"
	"log"
	"time"

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

	fmt.Println(core.QrCode)
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

	err = core.StatusNotify()
	if err != nil {
		log.Fatal(err)
	}

	core.LastSyncTime = time.Now().UnixNano()
}
