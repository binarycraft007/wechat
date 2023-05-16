package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/binarycraft007/wechat/core"
)

type PeriodicSyncOption struct {
	Cancel context.CancelFunc
	Period time.Duration
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var err error
	var wechatCore *core.Core

	if wechatCore, err = core.New(); err != nil {
		log.Fatal(err)
	}

	if err = wechatCore.GetUUID(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(wechatCore.QrCode)    // print qrcode
	fmt.Println(wechatCore.QrCodeUrl) // qrcode url

	if err = wechatCore.PreLogin(); err != nil {
		log.Fatal(err)
	}

	if err = wechatCore.Login(); err != nil {
		log.Fatal(err)
	}

	if err = wechatCore.Init(); err != nil {
		log.Fatal(err)
	}

	if err = wechatCore.StatusNotify(); err != nil {
		log.Fatal(err)
	}

	if err = wechatCore.GetContact(); err != nil {
		log.Fatal(err)
	}

	go periodicSync(wechatCore, PeriodicSyncOption{
		Cancel: cancel,
		Period: 20,
	}) // Call as a goroutine

	select {
	case <-ctx.Done(): // When sync returned 1101
		if err = wechatCore.Logout(); err != nil {
			log.Println(err.Error())
		}
		log.Println("logged out")
	}
}

func periodicSync(w *core.Core, options PeriodicSyncOption) {
	t := time.NewTicker(options.Period * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-t.C: // Activate periodically
			if err := w.SyncPolling(); err != nil {
				options.Cancel()
			}
		}
	}
}
