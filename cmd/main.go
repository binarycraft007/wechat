package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/binarycraft007/wechat"
	"github.com/gin-gonic/gin"
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

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	select {
	case <-ctx.Done(): // When sync returned 1101
		if err = wechatCore.Logout(); err != nil {
			log.Println(err.Error())
		}
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown: ", err)
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
				continue
			}
			if errors.As(wechat.ErrAlreadyLoggedOut, &err) {
				options.Cancel()
				return
			}
			log.Println("sync error:", err.Error())
		}
	}
}
