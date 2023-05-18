package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/binarycraft007/wechat/utils"
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
		Message: MessageRequest{
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

func (core *Core) UploadMedia(file *os.File, to string) error {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

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
		// TODO handle more file types
		return ErrUnknownFileType
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	baseRequest, err := core.GetBaseRequest()
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("f", "json")

	u, err := url.ParseRequestURI(core.Config.Api.UploadMedia)
	if err != nil {
		return err
	}
	u.RawQuery = params.Encode()

	clientMsgId := utils.GetClientMsgId()

	data := UploadMediaRequest{
		BaseRequest:   *baseRequest,
		ClientMediaId: clientMsgId,
		TotalLen:      fileInfo.Size(),
		StartPos:      0,
		DataLen:       fileInfo.Size(),
		MediaType:     4,
		UploadType:    2,
		FromUserName:  core.User.UserName,
		ToUserName:    to,
	}

	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	gmt := time.Now().UTC().Format(http.TimeFormat)

	formData := &bytes.Buffer{}
	writer := multipart.NewWriter(formData)

	// Add the form fields to the form.
	writer.WriteField("name", fileInfo.Name())
	writer.WriteField("type", mimeType)
	writer.WriteField("lastModifiedDate", gmt)
	writer.WriteField("size", fmt.Sprintf("%d", fileInfo.Size()))
	writer.WriteField("mediatype", mediaType)
	writer.WriteField("uploadmediarequest", string(marshalled))
	writer.WriteField("webwx_data_ticket", core.SessionData.DataTicket)
	writer.WriteField("pass_ticket", core.SessionData.PassTicket)

	// Create a new form field for the file.
	part, err := writer.CreateFormFile("filename", fileInfo.Name())
	if err != nil {
		return err
	}

	part.Write(fileBytes)

	// Close writer before use it in post request
	writer.Close()

	req, err := http.NewRequest("POST", u.String(), formData)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
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

	var result UploadMediaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.BaseResponse.Ret != 0 {
		return errors.New(result.BaseResponse.ErrMsg)
	}

	return nil
}
