package template

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/define"
	error2 "github.com/dcsunny/wechat/error"
	"github.com/dcsunny/wechat/util"
)

const (
	templateSendURL          = "https://api.weixin.qq.com/cgi-bin/message/template/send"
	templateSubscribeSendURL = "https://api.weixin.qq.com/cgi-bin/message/template/subscribe"
	templateMiniSendURL      = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send"
	templateMiniOrMpSendURL  = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/uniform_send"
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

type MiniMessage struct {
	ToUser          string               `json:"touser"`      // 必须, 接受者OpenID
	TemplateID      string               `json:"template_id"` // 必须, 模版ID
	Data            map[string]*DataItem `json:"data"`        // 必须, 模板数据
	Page            string               `json:"page"`
	FormID          string               `json:"form_id"`
	EmphasisKeyword string               `json:"emphasis_keyword"`
}

//DataItem 模版内某个 .DATA 的值
type DataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
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
		err = error2.CommonErrorHandle(result.CommonError, tpl.Context, "TemplateSend")
		return
	}
	msgID = result.MsgID
	return
}

func (tpl *Template) MiniSend(msg *MiniMessage) (templateID string, err error) {
	var accessToken string
	accessToken, err = tpl.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", templateMiniSendURL, accessToken)
	response, err := util.PostJSON(uri, msg)

	var result resTemplateMiniSend
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("body:%s", string(response))
		return
	}
	if result.ErrCode != 0 {
		err = error2.CommonErrorHandle(result.CommonError, tpl.Context, "TemplateMiniSend")
		return
	}
	templateID = result.TemplateID
	return
}

type MiniMpMessage struct {
	ToUser           string           `json:"touser"`
	WeappTemplateMsg WeappTemplateMsg `json:"weapp_template_msg"`
	MpTemplateMsg    MpTemplateMsg    `json:"mp_template_msg"`
}

type WeappTemplateMsg struct {
	TemplateID      string               `json:"template_id"`
	Page            string               `json:"page"`
	FormID          string               `json:"form_id"`
	Data            map[string]*DataItem `json:"data"`
	EmphasisKeyword string               `json:"emphasis_keyword"`
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
		err = error2.CommonErrorHandle(result.CommonError, tpl.Context, "TemplateSendMiniOrMp")
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
		err = error2.CommonErrorHandle(result, tpl.Context, "TemplateSendSubscribeMessage")
		return
	}
	return
}
