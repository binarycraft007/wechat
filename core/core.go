package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
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
	UUID       string
	S_Key      string
	S_Id       string
	U_In       string
	PassTicket string
	DataTicket string
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
	client := &http.Client{}
	resp, err := client.Post(core.Config.Api.JsLogin, "plain/text", http.NoBody)
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

	client := &http.Client{}
	resp, err := client.Get(urlStr)
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

func (core *Core) Login() error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", core.RedirectUri, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Add("client-version", "2.0.0")
	req.Header.Add("referer", "https://wx.qq.com/?&lang=zh_CN&target=t")
	req.Header.Add("extspam", "Go8FCIkFEokFCggwMDAwMDAwMRAGGvAESySibk50w5Wb3uTl2c2h64jVVrV7gNs06GFlWplHQbY/5FfiO++1yH4ykCyNPWKXmco+wfQzK5R98D3so7rJ5LmGFvBLjGceleySrc3SOf2Pc1gVehzJgODeS0lDL3/I/0S2SSE98YgKleq6Uqx6ndTy9yaL9qFxJL7eiA/R3SEfTaW1SBoSITIu+EEkXff+Pv8NHOk7N57rcGk1w0ZzRrQDkXTOXFN2iHYIzAAZPIOY45Lsh+A4slpgnDiaOvRtlQYCt97nmPLuTipOJ8Qc5pM7ZsOsAPPrCQL7nK0I7aPrFDF0q4ziUUKettzW8MrAaiVfmbD1/VkmLNVqqZVvBCtRblXb5FHmtS8FxnqCzYP4WFvz3T0TcrOqwLX1M/DQvcHaGGw0B0y4bZMs7lVScGBFxMj3vbFi2SRKbKhaitxHfYHAOAa0X7/MSS0RNAjdwoyGHeOepXOKY+h3iHeqCvgOH6LOifdHf/1aaZNwSkGotYnYScW8Yx63LnSwba7+hESrtPa/huRmB9KWvMCKbDThL/nne14hnL277EDCSocPu3rOSYjuB9gKSOdVmWsj9Dxb/iZIe+S6AiG29Esm+/eUacSba0k8wn5HhHg9d4tIcixrxveflc8vi2/wNQGVFNsGO6tB5WF0xf/plngOvQ1/ivGV/C1Qpdhzznh0ExAVJ6dwzNg7qIEBaw+BzTJTUuRcPk92Sn6QDn2Pu3mpONaEumacjW4w6ipPnPw+g2TfywJjeEcpSZaP4Q3YV5HG8D6UjWA4GSkBKculWpdCMadx0usMomsSS/74QgpYqcPkmamB4nVv1JxczYITIqItIKjD35IGKAUwAA==")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMovedPermanently {
		return errors.New("http status error: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(`<ret>(.*)<\/ret>`)
	if err != nil {
		return err
	}

	pm := re.FindStringSubmatch(string(body))
	if len(pm) > 1 && pm[1] == "0" {
		data := string(body)
		reSkey, err := regexp.Compile(`<skey>(.*)<\/skey>`)
		if err != nil {
			return err
		}
		core.SessionData.S_Key = reSkey.FindStringSubmatch(data)[1]

		reWxsid, err := regexp.Compile(`<wxsid>(.*)<\/wxsid>`)
		if err != nil {
			return err
		}
		core.SessionData.S_Id = reWxsid.FindStringSubmatch(data)[1]

		reWxuin, err := regexp.Compile(`<wxuin>(.*)<\/wxuin>`)
		if err != nil {
			return err
		}
		core.SessionData.U_In = reWxuin.FindStringSubmatch(data)[1]

		rePassTicket, err := regexp.Compile(`<pass_ticket>(.*)<\/pass_ticket>`)
		if err != nil {
			return err
		}
		core.SessionData.PassTicket = rePassTicket.FindStringSubmatch(data)[1]
	}

	for name, values := range resp.Header {
		// Loop over all values for the name.
		for _, value := range values {
			if name == "Set-Cookie" {
				cookie := string(value)
				re, err := regexp.Compile("=(.*?);")
				if err != nil {
					return err
				}
				if strings.Contains(cookie, "webwx") &&
					strings.Contains(cookie, "data") &&
					strings.Contains(cookie, "ticket") {
					core.SessionData.DataTicket = re.FindStringSubmatch(cookie)[1]
					log.Println(re.FindStringSubmatch(cookie)[1])
				} else if strings.Contains(cookie, "wxuin") {
					core.SessionData.U_In = re.FindStringSubmatch(cookie)[1]
					log.Println(re.FindStringSubmatch(cookie)[1])

				} else if strings.Contains(cookie, "wxsid") {
					core.SessionData.S_Id = re.FindStringSubmatch(cookie)[1]
					log.Println(re.FindStringSubmatch(cookie)[1])

				} else if strings.Contains(cookie, "pass_ticket") {
					core.SessionData.PassTicket = re.FindStringSubmatch(cookie)[1]
					log.Println(re.FindStringSubmatch(cookie)[1])
				}
			}
		}
	}

	return nil
}
