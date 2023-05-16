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
	"github.com/skip2/go-qrcode"

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

type StatusNotifyRequest struct {
	BaseRequest  BaseRequest `json:"BaseRequest"`
	Code         int         `json:"Code"`
	FromUserName string      `json:"FromUserName"`
	ToUserName   string      `json:"ToUserName"`
	ClientMsgId  int64       `json:"ClientMsgId"`
}

type BatchGetContactRequest struct {
	BaseRequest BaseRequest `json:"BaseRequest"`
	Count       int         `json:"Count"`
	List        []Contact   `json:"List"`
}

type SyncRequest struct {
	BaseRequest BaseRequest `json:"BaseRequest"`
	SyncKey     SyncKey     `json:"SyncKey"`
	RR          int64       `json:"rr"`
}

type BaseRequest struct {
	Uin      int64  `json:"Uin"`
	Sid      string `json:"Sid"`
	Skey     string `json:"Skey"`
	DeviceID string `json:"DeviceID"`
}

type BaseResponse struct {
	Ret    int    `json:"Ret"`
	ErrMsg string `json:"ErrMsg"`
}

type Contact struct {
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
}

type BatchGetContactResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	Count        int          `json:"Count"`
	ContactList  []Contact    `json:"ContactList"`
}

type SyncKey struct {
	Count int `json:"Count"`
	List  []struct {
		Key int `json:"Key"`
		Val int `json:"Val"`
	} `json:"List"`
}

type InitResponse struct {
	BaseResponse        BaseResponse `json:"BaseResponse"`
	Count               int          `json:"Count"`
	ContactList         []Contact    `json:"ContactList"`
	SyncKey             SyncKey      `json:"SyncKey"`
	User                User         `json:"User"`
	ChatSet             string       `json:"ChatSet"`
	SKey                string       `json:"SKey"`
	ClientVersion       int          `json:"ClientVersion"`
	SystemTime          int          `json:"SystemTime"`
	GrayScale           int          `json:"GrayScale"`
	InviteStartCount    int          `json:"InviteStartCount"`
	MPSubscribeMsgCount int          `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList  []any        `json:"MPSubscribeMsgList"`
	ClickReportInterval int          `json:"ClickReportInterval"`
}

type StatusNotifyResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MsgID        string       `json:"MsgID"`
}

type SyncResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	AddMsgCount  int          `json:"AddMsgCount"`
	AddMsgList   []struct {
		MsgID                string `json:"MsgId"`
		FromUserName         string `json:"FromUserName"`
		ToUserName           string `json:"ToUserName"`
		MsgType              int    `json:"MsgType"`
		Content              string `json:"Content"`
		Status               int    `json:"Status"`
		ImgStatus            int    `json:"ImgStatus"`
		CreateTime           int    `json:"CreateTime"`
		VoiceLength          int    `json:"VoiceLength"`
		PlayLength           int    `json:"PlayLength"`
		FileName             string `json:"FileName"`
		FileSize             string `json:"FileSize"`
		MediaID              string `json:"MediaId"`
		URL                  string `json:"Url"`
		AppMsgType           int    `json:"AppMsgType"`
		StatusNotifyCode     int    `json:"StatusNotifyCode"`
		StatusNotifyUserName string `json:"StatusNotifyUserName"`
		RecommendInfo        struct {
			UserName   string `json:"UserName"`
			NickName   string `json:"NickName"`
			QQNum      int    `json:"QQNum"`
			Province   string `json:"Province"`
			City       string `json:"City"`
			Content    string `json:"Content"`
			Signature  string `json:"Signature"`
			Alias      string `json:"Alias"`
			Scene      int    `json:"Scene"`
			VerifyFlag int    `json:"VerifyFlag"`
			AttrStatus int    `json:"AttrStatus"`
			Sex        int    `json:"Sex"`
			Ticket     string `json:"Ticket"`
			OpCode     int    `json:"OpCode"`
		} `json:"RecommendInfo"`
		ForwardFlag int `json:"ForwardFlag"`
		AppInfo     struct {
			AppID string `json:"AppID"`
			Type  int    `json:"Type"`
		} `json:"AppInfo"`
		HasProductID  int    `json:"HasProductId"`
		Ticket        string `json:"Ticket"`
		ImgHeight     int    `json:"ImgHeight"`
		ImgWidth      int    `json:"ImgWidth"`
		SubMsgType    int    `json:"SubMsgType"`
		NewMsgID      int64  `json:"NewMsgId"`
		OriContent    string `json:"OriContent"`
		EncryFileName string `json:"EncryFileName"`
	} `json:"AddMsgList"`
	ModContactCount        int   `json:"ModContactCount"`
	ModContactList         []any `json:"ModContactList"`
	DelContactCount        int   `json:"DelContactCount"`
	DelContactList         []any `json:"DelContactList"`
	ModChatRoomMemberCount int   `json:"ModChatRoomMemberCount"`
	ModChatRoomMemberList  []any `json:"ModChatRoomMemberList"`
	Profile                struct {
		BitFlag  int `json:"BitFlag"`
		UserName struct {
			Buff string `json:"Buff"`
		} `json:"UserName"`
		NickName struct {
			Buff string `json:"Buff"`
		} `json:"NickName"`
		BindUin   int `json:"BindUin"`
		BindEmail struct {
			Buff string `json:"Buff"`
		} `json:"BindEmail"`
		BindMobile struct {
			Buff string `json:"Buff"`
		} `json:"BindMobile"`
		Status            int    `json:"Status"`
		Sex               int    `json:"Sex"`
		PersonalCard      int    `json:"PersonalCard"`
		Alias             string `json:"Alias"`
		HeadImgUpdateFlag int    `json:"HeadImgUpdateFlag"`
		HeadImgURL        string `json:"HeadImgUrl"`
		Signature         string `json:"Signature"`
	} `json:"Profile"`
	ContinueFlag int     `json:"ContinueFlag"`
	SyncKey      SyncKey `json:"SyncKey"`
	SKey         string  `json:"SKey"`
	SyncCheckKey SyncKey `json:"SyncCheckKey"`
}

type GetContactResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MemberCount  int          `json:"MemberCount"`
	MemberList   []any        `json:"MemberList"`
	Seq          int          `json:"Seq"`
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
	Config          utils.Config
	SessionData     SessionData
	User            User
	Avatar          string
	RedirectUri     string
	QrCodeUrl       string
	QrCode          string
	NotifyUserName  string
	ContactMap      map[string]Contact
	LastSyncTime    int64
	SyncKey         SyncKey
	SyncSelector    int
	FormatedSyncKey string
	ContactSeq      int
}

func New() (*Core, error) {
	core := Core{}
	config, err := utils.NewConfig(utils.ConfigOption{})
	if err != nil {
		return nil, err
	}

	core.Config = *config
	return &core, nil
}

func (core *Core) GetUUID() error {
	client := &http.Client{}
	resp, err := client.Post(core.Config.Api.JsLogin, "", http.NoBody)
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
	qrCodeContent := "https://login.weixin.qq.com/l/" + uuid

	qrCode, err := qrcode.New(qrCodeContent, qrcode.Medium)
	if err != nil {
		return err
	}

	core.QrCode = qrCode.ToSmallString(false)
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

		config, err := utils.NewConfig(utils.ConfigOption{Host: u.Hostname()})
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
		core.PreLogin()
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
		return errors.New(result.BaseResponse.ErrMsg)
	}

	if len(result.SKey) > 0 {
		core.SessionData.Skey = result.SKey
	}

	core.SyncKey = result.SyncKey
	core.SetFormatedSyncKey(result.SyncKey)

	core.User = result.User
	core.ContactMap = make(map[string]Contact)

	for _, contact := range result.ContactList {
		core.ContactMap[contact.UserName] = contact
	}

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

func (core *Core) StatusNotify() error {
	params := url.Values{}
	params.Add("pass_ticket", core.SessionData.PassTicket)
	params.Add("lang", "zh_CN")

	u, err := url.ParseRequestURI(core.Config.Api.StatusNotify)
	if err != nil {
		return err
	}
	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	baseRequest, err := core.GetBaseRequest()
	if err != nil {
		return err
	}

	var code int
	var userName string
	if len(core.NotifyUserName) > 0 {
		code = 1
		userName = core.NotifyUserName
	} else {
		code = 3
		userName = core.User.UserName
	}
	core.NotifyUserName = ""

	data := StatusNotifyRequest{
		BaseRequest:  *baseRequest,
		Code:         code,
		FromUserName: core.User.UserName,
		ToUserName:   userName,
		ClientMsgId:  time.Now().UnixNano(),
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

	var result StatusNotifyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.BaseResponse.Ret != 0 {
		return errors.New(result.BaseResponse.ErrMsg)
	}

	return nil
}

func (core *Core) GetContact() error {
	ts := time.Now().UnixNano() / int64(time.Millisecond)

	params := url.Values{}
	params.Add("seq", fmt.Sprintf("%d", int(core.ContactSeq)))
	params.Add("skey", core.SessionData.Skey)
	params.Add("r", fmt.Sprintf("%d", int64(ts)))

	core.ContactSeq = 0

	u, err := url.ParseRequestURI(core.Config.Api.GetContact)
	if err != nil {
		return err
	}
	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result GetContactResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.Seq > 0 {
		core.ContactSeq = result.Seq
		if err = core.GetContact(); err != nil {
			return err
		}
		return nil
	}

	if result.Seq == 0 {
		contacts := make([]Contact, len(core.ContactMap))
		i := 0
		for _, contact := range core.ContactMap {
			if strings.HasPrefix(contact.UserName, "@@") &&
				contact.MemberCount == 0 {
				contacts[i] = contact
				i++
			}
		}
		err := core.BatchGetContact(contacts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (core *Core) BatchGetContact(contacts []Contact) error {
	ts := time.Now().UnixNano() / int64(time.Millisecond)

	params := url.Values{}
	params.Add("pass_ticket", core.SessionData.PassTicket)
	params.Add("type", "ex")
	params.Add("r", fmt.Sprintf("%d", int64(ts)))
	params.Add("lang", "zh_CN")

	u, err := url.ParseRequestURI(core.Config.Api.BatchGetContact)
	if err != nil {
		return err
	}
	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	baseRequest, err := core.GetBaseRequest()
	if err != nil {
		return err
	}

	data := BatchGetContactRequest{
		BaseRequest: *baseRequest,
		Count:       len(contacts),
		List:        contacts,
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

	var result BatchGetContactResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.BaseResponse.Ret != 0 {
		return errors.New(result.BaseResponse.ErrMsg)
	}

	for _, contact := range result.ContactList {
		core.ContactMap[contact.UserName] = contact
	}

	return nil
}

func (core *Core) SyncCheck() error {
	ts := time.Now().UnixNano() / int64(time.Millisecond)

	params := url.Values{}
	params.Add("r", fmt.Sprintf("%d", int64(ts)))
	params.Add("sid", core.SessionData.Sid)
	params.Add("uid", core.SessionData.Uin)
	params.Add("skey", core.SessionData.Skey)
	params.Add("deviceid", utils.GetDeviceID())
	params.Add("synckey", core.FormatedSyncKey)

	u, err := url.ParseRequestURI(core.Config.Api.SyncCheck)
	if err != nil {
		return err
	}

	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(body), "retcode:\"1101\"") {
		return errors.New("already logged out")
	}

	start := strings.Index(string(body), "selector:")
	start += len("selector:") + 1
	end := len(string(body)) - 2

	selectorStr := string(body)[start:end]
	selector, err := strconv.ParseInt(selectorStr, 10, 64)
	if err != nil {
		return err
	}
	core.SyncSelector = int(selector)

	return nil
}

func (core *Core) Sync() error {
	ts := ^time.Now().UnixNano()

	params := url.Values{}
	params.Add("sid", core.SessionData.Sid)
	params.Add("skey", core.SessionData.Skey)
	params.Add("pass_ticket", core.SessionData.PassTicket)
	params.Add("lang", "zh_CN")

	u, err := url.ParseRequestURI(core.Config.Api.Sync)
	if err != nil {
		return err
	}

	u.RawQuery = params.Encode()
	urlStr := fmt.Sprintf("%v", u)

	baseRequest, err := core.GetBaseRequest()
	if err != nil {
		return err
	}

	data := SyncRequest{
		BaseRequest: *baseRequest,
		SyncKey:     core.SyncKey,
		RR:          int64(ts),
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

	var result SyncResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.BaseResponse.Ret != 0 {
		return errors.New(result.BaseResponse.ErrMsg)
	}

	core.SyncKey = result.SyncCheckKey
	core.SetFormatedSyncKey(core.SyncKey)

	return nil
}

func (core *Core) SetFormatedSyncKey(syncKey SyncKey) {
	syncKeyList := make([]string, len(syncKey.List))
	for i, item := range syncKey.List {
		syncKeyList[i] = strconv.Itoa(item.Key) + "_" +
			strconv.Itoa(item.Val)
	}
	core.FormatedSyncKey = strings.Join(syncKeyList, "|")
}
