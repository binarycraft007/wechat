package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/binarycraft007/wechat"
	"github.com/gin-gonic/gin"
)

func demoHandler(c *gin.Context) {
	to := "filehelper"
	msg := "Welcome Gin Server"
	if err := core.SendMsg(msg, to); err != nil {
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
	if err := core.SendMsg(msgPng, to); err != nil {
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
	if err := core.SendMsg(msgMp4, to); err != nil {
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
	if err := core.SendMsg(msgTxt, to); err != nil {
		log.Println("Send message error:", err)
	}

	c.String(http.StatusOK, "Welcome Gin Server")
}
