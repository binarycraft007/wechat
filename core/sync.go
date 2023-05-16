package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/binarycraft007/wechat/core/utils"
)

func (core *Core) StatusNotify() error {
	params := url.Values{}
	params.Add("pass_ticket", core.SessionData.PassTicket)
	params.Add("lang", "zh_CN")

	u, err := url.ParseRequestURI(core.Config.Api.StatusNotify)
	if err != nil {
		return err
	}
	u.RawQuery = params.Encode()

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

	var result StatusNotifyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.BaseResponse.Ret != 0 {
		return errors.New(result.BaseResponse.ErrMsg)
	}

	return nil
}

func (core *Core) SyncCheck() error {
	ts := time.Now().UnixNano() / int64(time.Millisecond)

	params := url.Values{}
	params.Add("r", fmt.Sprintf("%d", int64(ts)))
	params.Add("sid", core.SessionData.Sid)
	params.Add("uin", core.SessionData.Uin)
	params.Add("skey", core.SessionData.Skey)
	params.Add("deviceid", utils.GetDeviceID())
	params.Add("synckey", core.FormatedSyncKey)

	u, err := url.ParseRequestURI(core.Config.Api.SyncCheck)
	if err != nil {
		return err
	}

	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := core.Client.Do(req)
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
