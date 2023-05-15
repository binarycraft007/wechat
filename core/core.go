package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	Uin               int    `json:"Uin"`
	UserName          string `json:"UserName"`
	NickName          string `json:"NickName"`
	HeadImgURL        string `json:"HeadImgUrl"`
	RemarkName        string `json:"RemarkName"`
	PYInitial         string `json:"PYInitial"`
	PYQuanPin         string `json:"PYQuanPin"`
	RemarkPYInitial   string `json:"RemarkPYInitial"`
	RemarkPYQuanPin   string `json:"RemarkPYQuanPin"`
	HideInputBarFlag  int    `json:"HideInputBarFlag"`
	StarFriend        int    `json:"StarFriend"`
	Sex               int    `json:"Sex"`
	Signature         string `json:"Signature"`
	AppAccountFlag    int    `json:"AppAccountFlag"`
	VerifyFlag        int    `json:"VerifyFlag"`
	ContactFlag       int    `json:"ContactFlag"`
	WebWxPluginSwitch int    `json:"WebWxPluginSwitch"`
	HeadImgFlag       int    `json:"HeadImgFlag"`
	SnsFlag           int    `json:"SnsFlag"`
}

type InitRequest struct {
	BaseRequest BaseRequest `json:"BaseRequest"`
}

type BaseRequest struct {
	Uin      int64  `json:"Uin"`
	Sid      string `json:"Sid"`
	Skey     string `json:"Skey"`
	DeviceID string `json:"DeviceID"`
}

type InitResponse struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
	Count       int `json:"Count"`
	ContactList []struct {
		Uin              int    `json:"Uin"`
		UserName         string `json:"UserName"`
		NickName         string `json:"NickName"`
		HeadImgURL       string `json:"HeadImgUrl"`
		ContactFlag      int    `json:"ContactFlag"`
		MemberCount      int    `json:"MemberCount"`
		MemberList       []any  `json:"MemberList"`
		RemarkName       string `json:"RemarkName"`
		HideInputBarFlag int    `json:"HideInputBarFlag"`
		Sex              int    `json:"Sex"`
		Signature        string `json:"Signature"`
		VerifyFlag       int    `json:"VerifyFlag"`
		OwnerUin         int    `json:"OwnerUin"`
		PYInitial        string `json:"PYInitial"`
		PYQuanPin        string `json:"PYQuanPin"`
		RemarkPYInitial  string `json:"RemarkPYInitial"`
		RemarkPYQuanPin  string `json:"RemarkPYQuanPin"`
		StarFriend       int    `json:"StarFriend"`
		AppAccountFlag   int    `json:"AppAccountFlag"`
		Statues          int    `json:"Statues"`
		AttrStatus       int    `json:"AttrStatus"`
		Province         string `json:"Province"`
		City             string `json:"City"`
		Alias            string `json:"Alias"`
		SnsFlag          int    `json:"SnsFlag"`
		UniFriend        int    `json:"UniFriend"`
		DisplayName      string `json:"DisplayName"`
		ChatRoomID       int    `json:"ChatRoomId"`
		KeyWord          string `json:"KeyWord"`
		EncryChatRoomID  string `json:"EncryChatRoomId"`
		IsOwner          int    `json:"IsOwner"`
	} `json:"ContactList"`
	SyncKey struct {
		Count int `json:"Count"`
		List  []struct {
			Key int `json:"Key"`
			Val int `json:"Val"`
		} `json:"List"`
	} `json:"SyncKey"`
	User struct {
		Uin               int    `json:"Uin"`
		UserName          string `json:"UserName"`
		NickName          string `json:"NickName"`
		HeadImgURL        string `json:"HeadImgUrl"`
		RemarkName        string `json:"RemarkName"`
		PYInitial         string `json:"PYInitial"`
		PYQuanPin         string `json:"PYQuanPin"`
		RemarkPYInitial   string `json:"RemarkPYInitial"`
		RemarkPYQuanPin   string `json:"RemarkPYQuanPin"`
		HideInputBarFlag  int    `json:"HideInputBarFlag"`
		StarFriend        int    `json:"StarFriend"`
		Sex               int    `json:"Sex"`
		Signature         string `json:"Signature"`
		AppAccountFlag    int    `json:"AppAccountFlag"`
		VerifyFlag        int    `json:"VerifyFlag"`
		ContactFlag       int    `json:"ContactFlag"`
		WebWxPluginSwitch int    `json:"WebWxPluginSwitch"`
		HeadImgFlag       int    `json:"HeadImgFlag"`
		SnsFlag           int    `json:"SnsFlag"`
	} `json:"User"`
	ChatSet             string `json:"ChatSet"`
	SKey                string `json:"SKey"`
	ClientVersion       int    `json:"ClientVersion"`
	SystemTime          int    `json:"SystemTime"`
	GrayScale           int    `json:"GrayScale"`
	InviteStartCount    int    `json:"InviteStartCount"`
	MPSubscribeMsgCount int    `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList  []any  `json:"MPSubscribeMsgList"`
	ClickReportInterval int    `json:"ClickReportInterval"`
}

type SessionData struct {
	UUID       string
	Skey       string
	Sid        string
	Uin        string
	PassTicket string
	DataTicket string
}

type Core struct {
	Config      utils.Config
	Events      events.EventEmmiter
	SessionData SessionData
	User        User
	Avatar      string
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
	params.Add("r", fmt.Sprintf("%d", int64(ts)))

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

		u.Path = ""
		u.RawQuery = ""
		u.Fragment = ""

		urlStr := fmt.Sprintf("%v", u)

		skipLen := len("https://")
		config, err := utils.NewConfig(utils.ConfigOption{Host: urlStr[skipLen:]})
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

		core.Avatar = string(body)[start:end]
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
		core.SessionData.Skey = reSkey.FindStringSubmatch(data)[1]

		reWxsid, err := regexp.Compile(`<wxsid>(.*)<\/wxsid>`)
		if err != nil {
			return err
		}
		core.SessionData.Sid = reWxsid.FindStringSubmatch(data)[1]

		reWxuin, err := regexp.Compile(`<wxuin>(.*)<\/wxuin>`)
		if err != nil {
			return err
		}
		core.SessionData.Uin = reWxuin.FindStringSubmatch(data)[1]

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
				} else if strings.Contains(cookie, "wxuin") {
					core.SessionData.Uin = re.FindStringSubmatch(cookie)[1]
				} else if strings.Contains(cookie, "wxsid") {
					core.SessionData.Sid = re.FindStringSubmatch(cookie)[1]
				} else if strings.Contains(cookie, "pass_ticket") {
					core.SessionData.PassTicket = re.FindStringSubmatch(cookie)[1]
				}
			}
		}
	}

	return nil
}

func (core *Core) Init() error {
	ts := time.Now().UnixNano()
	r := ts / -1579

	params := url.Values{}
	params.Add("pass_ticket", core.SessionData.PassTicket)
	params.Add("r", fmt.Sprintf("%d", int64(r)))

	u, err := url.ParseRequestURI(core.Config.Api.Init)
	if err != nil {
		return err
	}
	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	baseRequest, err := core.GetBaseRequest()
	if err != nil {
		return err
	}

	data := InitRequest{
		BaseRequest: *baseRequest,
	}

	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", urlStr, bytes.NewReader(marshalled))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("http status error: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result InitResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.BaseResponse.Ret == core.Config.SyncCheckRetLogout {
		return errors.New("already logged out")
	}

	if result.BaseResponse.Ret != 0 {
		return errors.New("core init error: ")
	}

	if len(result.SKey) > 0 {
		core.SessionData.Skey = result.SKey
	}

	core.User = result.User
	return nil
}

func (core *Core) GetBaseRequest() (*BaseRequest, error) {
	uin, err := strconv.ParseInt(core.SessionData.Uin, 10, 64)
	if err != nil {
		return nil, err
	}

	return &BaseRequest{
		Uin:      uin,
		Sid:      core.SessionData.Sid,
		Skey:     core.SessionData.Skey,
		DeviceID: utils.GetDeviceID(),
	}, nil
}
