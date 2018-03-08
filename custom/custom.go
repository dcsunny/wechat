package custom

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/util"
)

const (
	customUrl       = "https://api.weixin.qq.com/cgi-bin/message/custom/send"
	CustomTypeText  = "text"
	CustomTypeImage = "image"
)

//客服消息
type Custom struct {
	*context.Context
}

func NewCustom(context *context.Context) *Custom {
	tpl := new(Custom)
	tpl.Context = context
	return tpl
}

type Message struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text,omitempty"`
	Image struct {
		MediaId string `json:"media_id"`
	} `json:"image,omitempty"`
}

//Send 发送客服消息

func (tpl *Custom) Send(msg *Message) (result util.CommonError, err error) {
	var accessToken string
	accessToken, err = tpl.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", customUrl, accessToken)
	response, err := util.PostJSON(uri, msg)

	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("custom msg send error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}
