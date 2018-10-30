package message_mass

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/util"
)

const (
	MessageMassSendByOpenIdURL = "https://api.weixin.qq.com/cgi-bin/message/mass/send?access_token=%s"
	MessageMassSendByTagURL    = "https://api.weixin.qq.com/cgi-bin/message/mass/sendall?access_token=%s"
)

type MessageMass struct {
	*context.Context
}

func NewMessageMass(context *context.Context) *MessageMass {
	service := new(MessageMass)
	service.Context = context
	return service
}

type MessageByOpen struct {
	Touser  []string `json:"touser"`
	Msgtype string   `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Clientmsgid string `json:"clientmsgid"`
}

type MessageByTag struct {
	Filter struct {
		IsToAll bool `json:"is_to_all"`
		TagID   int  `json:"tag_id"`
	} `json:"filter"`
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Clientmsgid string `json:"clientmsgid"`
}

type MessageSassResult struct {
	util.CommonError
	Type      string `json:"type"`
	MsgId     int64  `json:"msg_id"`
	MsgDataId int64  `json:"msg_data_id"`
}

func NewTextMessage(openids []string, content string, tag string) MessageByOpen {
	message := MessageByOpen{
		Touser:      openids,
		Msgtype:     "text",
		Clientmsgid: tag,
	}
	message.Text.Content = content
	return message
}

func (service *MessageMass) Send(msg *MessageByOpen) (result MessageSassResult, err error) {
	var accessToken string
	accessToken, err = service.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(MessageMassSendByOpenIdURL, accessToken)
	response, err := util.PostJSON(uri, msg)

	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("message mass send error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

func NewTextMessageByTag(tag string, wxTag int, content string) MessageByTag {
	message := MessageByTag{
		Msgtype:     "text",
		Clientmsgid: fmt.Sprint(wxTag) + tag,
	}
	message.Text.Content = content
	message.Filter.TagID = wxTag
	return message
}

func (service *MessageMass) SendByTag(msg *MessageByTag) (result MessageSassResult, err error) {
	var accessToken string
	accessToken, err = service.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(MessageMassSendByTagURL, accessToken)
	response, err := util.PostJSON(uri, msg)

	err = json.Unmarshal(response, &result)
	if err != nil {
		fmt.Println(fmt.Sprintf("message mass send error,json:%s", string(response)))
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("message mass send error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}
