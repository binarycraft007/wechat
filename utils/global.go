package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type AppMessage struct {
	Name    string
	Size    int
	MediaId string
	Ext     string
}

var ErrUnknownFileType = errors.New("unknown file type")

func GetDeviceID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("e%.15s", fmt.Sprintf("%0.15f", rand.Float64())[2:17])
}

func GetClientMsgId() int64 {
	milliseconds := time.Now().UnixNano() / int64(time.Millisecond)
	return int64(float64(milliseconds) * 1e3)
}

func GetErrorMsgInt(code int) string {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d, %s error: %d", file, line, f.Name(), code)
}

func GetErrorMsgStr(str string) string {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d, %s error: %s", file, line, f.Name(), str)
}

func DetectMediaType(fileBytes []byte) (*string, error) {
	mimeType := http.DetectContentType(fileBytes)

	var mediaType string
	if strings.HasPrefix(mimeType, "image/") {
		mediaType = "pic"
	} else if strings.HasPrefix(mimeType, "video/") {
		mediaType = "video"
	} else if strings.HasPrefix(mimeType, "text/") ||
		strings.HasPrefix(mimeType, "application/") {
		mediaType = "doc"
	} else if strings.HasPrefix(mimeType, "audio/") {
		mediaType = "audio"
	} else {
		return nil, ErrUnknownFileType
	}

	return &mediaType, nil
}

func GetAttachmentContent(msg AppMessage) string {
	// Print the content
	return fmt.Sprintf(`
    <appmsg appid='wxeb7ec651dd0aefa9' sdkver=''>
        <title>%s</title>
        <des></des>
        <action></action>
        <type>6</type>
        <content></content>
        <url></url>
        <lowurl></lowurl>
        <appattach>
            <totallen>%d</totallen>
            <attachid>%s</attachid>
            <fileext>%s</fileext>
        </appattach>
        <extinfo></extinfo>
    </appmsg>
    `, msg.Name, msg.Size, msg.MediaId, msg.Ext)
}
