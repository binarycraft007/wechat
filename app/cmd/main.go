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

	go periodicSync(wechatCore, PeriodicSyncOption{
		Cancel: cancel,
		Period: 20,
	}) // Call as a goroutine

	select {
	case <-ctx.Done(): // When sync returned 1101
		// TODO call logout
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
