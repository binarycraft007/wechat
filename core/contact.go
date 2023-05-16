package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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
		var contacts []Contact
		for _, contact := range result.MemberList {
			core.ContactMap[contact.UserName] = contact
			if strings.HasPrefix(contact.UserName, "@@") &&
				contact.MemberCount == 0 {
				contacts = append(contacts, contact)
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
