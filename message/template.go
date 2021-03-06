package message

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/common_error"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/define"
	"github.com/dcsunny/wechat/util"
)

const (
	templateSendURL              = "https://api.weixin.qq.com/cgi-bin/message/template/send"                //公众号模板消息
	templateSubscribeSendURL     = "https://api.weixin.qq.com/cgi-bin/message/template/subscribe"           //公众号一次性订阅消息
	templateMiniSendURL          = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send"         //微信小程序模板消息发送
	templateMiniOrMpSendURL      = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/uniform_send" //下发小程序和公众号统一的服务消息
	templateMiniSubscribeSendURL = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send"               //微信小程序一次性订阅消息
)

//Template 模板消息
type Template struct {
	*context.Context
}

//NewTemplate 实例化
func NewTemplate(context *context.Context) *Template {
	tpl := new(Template)
	tpl.Context = context
	return tpl
}

//Message 发送的模板消息内容
type Message struct {
	ToUser      string               `json:"touser"`          // 必须, 接受者OpenID
	TemplateID  string               `json:"template_id"`     // 必须, 模版ID
	URL         string               `json:"url,omitempty"`   // 可选, 用户点击后跳转的URL, 该URL必须处于开发者在公众平台网站中设置的域中
	Color       string               `json:"color,omitempty"` // 可选, 整个消息的颜色, 可以不设置
	Data        map[string]*DataItem `json:"data"`            // 必须, 模板数据
	MiniProgram struct {
		AppID    string `json:"appid"`    //所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系）
		PagePath string `json:"pagepath"` //所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar）
	} `json:"miniprogram"` //可选,跳转至小程序地址
}

//DataItem 模版内某个 .DATA 的值
type DataItem struct {
	Value interface{} `json:"value"`
	Color string      `json:"color,omitempty"`
}

type resTemplateSend struct {
	define.CommonError

	MsgID int64 `json:"msgid"`
}
type resTemplateMiniSend struct {
	define.CommonError
	TemplateID string `json:"template_id"`
}

//Send 发送模板消息
func (tpl *Template) Send(msg *Message) (msgID int64, err error) {
	var accessToken string
	accessToken, err = tpl.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", templateSendURL, accessToken)
	response, err := util.PostJSON(uri, msg)

	var result resTemplateSend
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("template msg send err,result:%s", string(response))
		return
	}
	if result.ErrCode != 0 {
		err = common_error.CommonErrorHandle(result.CommonError, tpl.Context, "TemplateSend")
		return
	}
	msgID = result.MsgID
	return
}

type MiniMpMessage struct {
	ToUser        string        `json:"touser"`
	MpTemplateMsg MpTemplateMsg `json:"mp_template_msg"`
}

type MpTemplateMsg struct {
	AppID       string `json:"appid"`
	TemplateID  string `json:"template_id"`
	Url         string `json:"url"`
	Miniprogram struct {
		Appid    string `json:"appid"`
		Pagepath string `json:"pagepath"`
	} `json:"miniprogram"`
	Data map[string]*DataItem `json:"data"`
}

func (tpl *Template) SendMiniOrMp(msg *MiniMpMessage) (err error) {
	var accessToken string
	accessToken, err = tpl.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", templateMiniOrMpSendURL, accessToken)
	response, err := util.PostJSON(uri, msg)

	var result resTemplateMiniSend
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("body:%s", string(response))
		return
	}
	if result.ErrCode != 0 {
		err = common_error.CommonErrorHandle(result.CommonError, tpl.Context, "TemplateSendMiniOrMp")
		return
	}
	return
}

type SubscribeMessage struct {
	Message
	Scene string
	Title string
}

func (tpl *Template) SendSubscribeMessage(msg *SubscribeMessage) (err error) {
	var accessToken string
	accessToken, err = tpl.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", templateSubscribeSendURL, accessToken)
	response, err := util.PostJSON(uri, msg)
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = common_error.CommonErrorHandle(result, tpl.Context, "TemplateSendSubscribeMessage")
		return
	}
	return
}

type MiniSubscribeMessage struct {
	ToUser           string               `json:"touser"`      // 必须, 接受者OpenID
	TemplateID       string               `json:"template_id"` // 必须, 模版ID
	Data             map[string]*DataItem `json:"data"`        // 必须, 模板数据
	Page             string               `json:"page"`
	MiniprogramState string               `json:"miniprogram_state"` //跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	Lang             string               `json:"lang"`              //进入小程序查看”的语言类型，支持zh_CN(简体中文)、en_US(英文)、zh_HK(繁体中文)、zh_TW(繁体中文)，默认为zh_CN
}

func (tpl *Template) SendMiniSubscribeMessage(msg *MiniSubscribeMessage) (err error) {
	var accessToken string
	accessToken, err = tpl.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", templateMiniSubscribeSendURL, accessToken)
	response, err := util.PostJSON(uri, msg)
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = common_error.CommonErrorHandle(result, tpl.Context, "TemplateMiniSendSubscribeMessage")
		return
	}
	return
}
