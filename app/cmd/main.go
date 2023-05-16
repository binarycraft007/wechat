package main

import (
	"fmt"
	"log"
	"time"

	"github.com/binarycraft007/wechat/core"
)

func main() {
	wechatCore, err := core.New()
	if err != nil {
		log.Fatal(err)
	}

	err = wechatCore.GetUUID()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(wechatCore.QrCode)
	fmt.Println(wechatCore.QrCodeUrl)

	err = wechatCore.PreLogin()
	if err != nil {
		log.Fatal(err)
	}

	err = wechatCore.Login()
	if err != nil {
		log.Fatal(err)
	}

	err = wechatCore.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = wechatCore.StatusNotify()
	if err != nil {
		log.Fatal(err)
	}

	err = wechatCore.GetContact()
	if err != nil {
		log.Fatal(err)
	}

	err = wechatCore.SyncCheck()
	if err != nil {
		log.Fatal(err)
	}

	wechatCore.LastSyncTime = time.Now().UnixNano()
}
