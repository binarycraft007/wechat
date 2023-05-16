package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/binarycraft007/wechat/core/utils"
)

type MessageType int

const (
	Text       MessageType = 1
	Image      MessageType = 3
	Voice      MessageType = 34
	Video      MessageType = 43
	MicroVideo MessageType = 62
	Emoticon   MessageType = 47
)

func (core *Core) SendText(msg string, to string) error {
	params := url.Values{}
	params.Add("pass_ticket", core.SessionData.PassTicket)
	params.Add("lang", "zh_CN")

	u, err := url.ParseRequestURI(core.Config.Api.SendMsg)
	if err != nil {
		return err
	}
	u.RawQuery = params.Encode()

	baseRequest, err := core.GetBaseRequest()
	if err != nil {
		return err
	}

	clientMsgId := utils.GetClientMsgId()

	data := SendTextRequest{
		BaseRequest: *baseRequest,
		Scene:       0,
		Message: Message{
			FromUserName: core.User.UserName,
			ToUserName:   to,
			Content:      msg,
			Type:         Text,
			ClientMsgId:  clientMsgId,
			LocalID:      clientMsgId,
		},
	}

	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(marshalled))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := core.Client.Do(req)
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

	var result SendMsgResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.BaseResponse.Ret != 0 {
		return errors.New(result.BaseResponse.ErrMsg)
	}

	return nil
}
