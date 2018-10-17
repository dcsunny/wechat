package message_mass

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/util"
)

const (
	MessageMassSendByOpenIdURL = "https://api.weixin.qq.com/cgi-bin/message/mass/send?access_token=%s"
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
}

type MessageSassResult struct {
	util.CommonError
	Type      string `json:"type"`
	MsgId     int64  `json:"msg_id"`
	MsgDataId int64  `json:"msg_data_id"`
}

func NewTextMessage(openids []string, content string) MessageByOpen {
	message := MessageByOpen{
		Touser:  openids,
		Msgtype: "text",
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
