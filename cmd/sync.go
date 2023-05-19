package main

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/binarycraft007/wechat"
)

type PingMessage struct {
	FromUserName string
	ToNickName   string
}

type PeriodicSyncOption struct {
	Cancel context.CancelFunc
	Period time.Duration
}

func periodicSync(options PeriodicSyncOption) {
	t := time.NewTicker(options.Period * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-t.C: // Activate periodically
			var err error
			if err = core.SyncPolling(); err == nil {
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

func onMsgRecv(data *wechat.SyncResponse) error {
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
