package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/binarycraft007/wechat"
	"github.com/gin-gonic/gin"
)

type PingMessage struct {
	FromUserName string
	ToNickName   string
}

type PeriodicSyncOption struct {
	Cancel context.CancelFunc
	Period time.Duration
}

func main() {
	var err error
	var wechatCore *wechat.Core

	if wechatCore, err = wechat.New(wechat.CoreOption{
		SyncMsgFunc:     onMsgRecv,
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

	interruptContext, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	ctx, cancel := context.WithCancel(interruptContext)

	defer cancel()
	defer stop()

	go periodicSync(wechatCore, PeriodicSyncOption{
		Cancel: cancel,
		Period: 20,
	}) // Call as a goroutine

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		to := "filehelper"
		msg := "Welcome Gin Server"
		if err := wechatCore.SendMsg(msg, to); err != nil {
			log.Println("Send message error:", err)
		}

		pngBytes, err := ioutil.ReadFile("media/zero.png")
		if err != nil {
			log.Println("Read file error:", err)
		}

		msgPng := wechat.MediaMessage{
			Name:      "zero.png",
			FileBytes: pngBytes,
		}
		if err := wechatCore.SendMsg(msgPng, to); err != nil {
			log.Println("Send message error:", err)
		}

		mp4Bytes, err := ioutil.ReadFile("media/gopher.mp4")
		if err != nil {
			log.Println("Read file error:", err)
		}

		msgMp4 := wechat.MediaMessage{
			Name:      "gopher.mp4",
			FileBytes: mp4Bytes,
		}
		if err := wechatCore.SendMsg(msgMp4, to); err != nil {
			log.Println("Send message error:", err)
		}

		txtBytes, err := ioutil.ReadFile("media/hello.txt")
		if err != nil {
			log.Println("Read file error:", err)
		}

		msgTxt := wechat.MediaMessage{
			Name:      "hello.txt",
			FileBytes: txtBytes,
		}
		if err := wechatCore.SendMsg(msgTxt, to); err != nil {
			log.Println("Send message error:", err)
		}

		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Println("listen error:", err)
			cancel()
		}
	}()

	select {
	case <-ctx.Done(): // When sync returned 1101
		if err = wechatCore.Logout(); err != nil {
			log.Println("logout error:", err.Error())
		}
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("shutdown error:", err.Error())
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
			} else {
				log.Println("sync error:", err.Error())
			}
		}
	}
}

func onMsgRecv(data *wechat.SyncResponse, core *wechat.Core) error {
	for _, message := range data.AddMsgList {
		if len(message.Content) == 0 {
			return nil
		}

		pingMsg := extractPingMessage(message.Content)
		userNick := core.User.NickName

		if pingMsg != nil &&
			strings.HasPrefix(pingMsg.ToNickName, userNick) {
			to := message.FromUserName
			msg := "What can I do for you?"
			if err := core.SendMsg(msg, to); err != nil {
				log.Println("Send message error:", err)
			}
		}
	}
	return nil
}

func extractPingMessage(message string) *PingMessage {
	if strings.HasPrefix(message, "@") {
		fromIdxEnd := strings.Index(message, ":")
		toIdxStart := strings.Index(message, "<br/>@")

		if fromIdxEnd == -1 || toIdxStart == -1 {
			return nil
		}

		padding := len("<br/>@")

		return &PingMessage{
			FromUserName: message[0:fromIdxEnd],
			ToNickName:   message[toIdxStart+padding:],
		}
	}
	return nil
}
