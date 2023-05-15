package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/binarycraft007/wechat/core/utils"

	"github.com/kataras/go-events"

	"net/http"
	"net/url"
)

type User struct {
	Avatar string
}

type SessionData struct {
	UUID string
}

type Core struct {
	Config      utils.Config
	Events      events.EventEmmiter
	SessionData SessionData
	User        User
	RedirectUri string
	QrCodeUrl   string
	QrCode      string
}

func New() (*Core, error) {
	core := Core{}

	core.Events = events.New()

	config, err := utils.NewConfig(utils.ConfigOption{})
	if err != nil {
		return nil, err
	}

	core.Config = *config

	core.Events.On("my_event", func(payload ...interface{}) {
		message := payload[0].(string)
		fmt.Println(message) // prints "this is my payload"
	})

	return &core, nil
}

func (core *Core) GetUUID() error {
	resp, err := http.Post(core.Config.Api.JsLogin, "plain/text", http.NoBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	start := strings.Index(string(body), "window.QRLogin.uuid = ")
	start += len("window.QRLogin.uuid = ") + 1
	end := len(string(body)) - 2

	uuid := string(body)[start:end]

	if resp.StatusCode != http.StatusOK {
		return errors.New("http status error: " + resp.Status)
	}

	core.QrCodeUrl = "https://login.weixin.qq.com/qrcode/" + uuid
	core.QrCode = "https://login.weixin.qq.com/l/" + uuid
	core.SessionData.UUID = uuid
	return nil
}

func (core *Core) PreLogin() error {
	ts := ^time.Now().UnixNano()
	params := url.Values{}
	params.Add("tip", "0")
	params.Add("uuid", core.SessionData.UUID)
	params.Add("loginicon", "true")
	params.Add("r", strconv.FormatInt(int64(ts), 10))

	u, err := url.ParseRequestURI(core.Config.Api.Login)
	if err != nil {
		return err
	}
	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	httpStatusSuccess := strings.Contains(string(body), "window.redirect_uri")
	httpStatusCreated := strings.Contains(string(body), "window.userAvatar")

	if !httpStatusCreated && !httpStatusSuccess {
		start := strings.Index(string(body), "window.code=")
		start += len("window.code=") + 1
		end := len(string(body)) - 2

		codeStr := string(body)[start:end]

		return errors.New("http status error: " + codeStr)
	}

	if httpStatusSuccess {
		start := strings.Index(string(body), "redirect_uri=")
		start += len("redirect_uri=") + 1
		end := len(string(body)) - 2
		redirectUri := string(body)[start:end]

		u, err := url.Parse(redirectUri)
		if err != nil {
			return err
		}

		q, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			return err
		}

		for key, _ := range q {
			q.Del(key)
		}

		u.RawQuery = q.Encode()
		urlStr := fmt.Sprintf("%v", u)

		config, err := utils.NewConfig(utils.ConfigOption{Host: urlStr})
		if err != nil {
			return err
		}

		core.Config = *config
		core.RedirectUri = redirectUri
	}

	if httpStatusCreated {
		start := strings.Index(string(body), "userAvatar = ")
		start += len("userAvatar = ") + 1
		end := len(string(body)) - 2

		core.User.Avatar = string(body)[start:end]
	}

	return nil
}
