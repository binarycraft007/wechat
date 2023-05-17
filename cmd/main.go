package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/binarycraft007/wechat"
)

type PeriodicSyncOption struct {
	Cancel context.CancelFunc
	Period time.Duration
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var err error
	var wechatCore *wechat.Core

	if wechatCore, err = wechat.New(wechat.CoreOption{
		SyncMsgFunc:     nil,
		SyncContactFunc: nil,
	}); err != nil {
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
		log.Println("logged out:", wechatCore.User.NickName)
	}
}

func periodicSync(w *wechat.Core, options PeriodicSyncOption) {
	t := time.NewTicker(options.Period * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-t.C: // Activate periodically
			var err error
			if err = w.SyncPolling(); err == nil {
				return
			}
			if errors.As(wechat.ErrAlreadyLoggedOut, err) {
				options.Cancel()
			}
		}
	}
}
