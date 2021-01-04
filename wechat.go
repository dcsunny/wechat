package wechat

import (
	"net/http"
	"sync"

	"github.com/dcsunny/wechat/shopping_guide"

	"github.com/dcsunny/wechat/message_mass"

	"github.com/dcsunny/wechat/message"

	"github.com/dcsunny/wechat/safe"

	"github.com/dcsunny/wechat/cache"
	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/js"
	"github.com/dcsunny/wechat/material"
	"github.com/dcsunny/wechat/menu"
	"github.com/dcsunny/wechat/miniprogram"
	"github.com/dcsunny/wechat/oauth"
	"github.com/dcsunny/wechat/pay"
	"github.com/dcsunny/wechat/qr"
	"github.com/dcsunny/wechat/server"
	"github.com/dcsunny/wechat/user"
)

// Wechat struct
type Wechat struct {
	Context *context.Context
}

// Config for user
type Config struct {
	AppID           string
	AppSecret       string
	Token           string
	EncodingAESKey  string
	PayMchID        string //支付 - 商户 ID
	PayNotifyURL    string //支付 - 接受微信支付结果通知的接口地址
	PayKey          string //支付 - 商户后台设置的支付 key
	AccessTokenURL  string
	PayCertPEMBlock string
	PayKeyPEMBlock  string
	Cache           cache.Cache
}

// NewWechat init
func NewWechat(cfg *Config) *Wechat {
	context := new(context.Context)
	copyConfigToContext(cfg, context)
	return &Wechat{context}
}

func copyConfigToContext(cfg *Config, context *context.Context) {
	context.AppID = cfg.AppID
	context.AppSecret = cfg.AppSecret
	context.Token = cfg.Token
	context.EncodingAESKey = cfg.EncodingAESKey
	context.PayMchID = cfg.PayMchID
	context.PayKey = cfg.PayKey
	context.PayNotifyURL = cfg.PayNotifyURL
	context.PayCertPEMBlock = cfg.PayCertPEMBlock
	context.PayKeyPEMBlock = cfg.PayKeyPEMBlock
	context.Cache = cfg.Cache
	context.AccessTokenURL = cfg.AccessTokenURL
	context.SetAccessTokenLock(new(sync.RWMutex))
	context.SetJsAPITicketLock(new(sync.RWMutex))
}

// GetServer 消息管理
func (wc *Wechat) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	wc.Context.Request = req
	wc.Context.Writer = writer
	return server.NewServer(wc.Context)
}

//GetAccessToken 获取access_token
func (wc *Wechat) GetAccessToken() (string, error) {
	return wc.Context.GetAccessToken()
}

// GetOauth oauth2网页授权
func (wc *Wechat) GetOauth() *oauth.Oauth {
	return oauth.NewOauth(wc.Context)
}

// GetMaterial 素材管理
func (wc *Wechat) GetMaterial() *material.Material {
	return material.NewMaterial(wc.Context)
}

// GetJs js-sdk配置
func (wc *Wechat) GetJs() *js.Js {
	return js.NewJs(wc.Context)
}

// GetMenu 菜单管理接口
func (wc *Wechat) GetMenu() *menu.Menu {
	return menu.NewMenu(wc.Context)
}

// GetUser 用户管理接口
func (wc *Wechat) GetUser() *user.User {
	return user.NewUser(wc.Context)
}

// GetTemplate 模板消息接口
func (wc *Wechat) GetTemplate() *message.Template {
	return message.NewTemplate(wc.Context)
}

// GetMessageMass 群发消息接口
func (wc *Wechat) GetMessageMass() *message_mass.MessageMass {
	return message_mass.NewMessageMass(wc.Context)
}

// GetPay 返回支付消息的实例
func (wc *Wechat) GetPay() *pay.Pay {
	return pay.NewPay(wc.Context)
}

// 带参数二维码接口
func (wc *Wechat) GetQrCode() *qr.QR {
	return qr.NewQR(wc.Context)
}

//小程序二维码
func (wc *Wechat) GetMiniQrCode() *miniprogram.MiniProgram {
	return miniprogram.NewMiniProgram(wc.Context)
}

//用户标签接口
func (wc *Wechat) GetTag() *user.Tag {
	return user.NewTag(wc.Context)
}

//客服消息接口
func (wc *Wechat) GetCustom() *message.Manager {
	return message.NewMessageManager(wc.Context)
}

func (wc *Wechat) GetSafe() *safe.WxSafe {
	return safe.NewWxSafe(wc.Context)
}

//对话能力接口
func (wc *Wechat) GetGuide() *shopping_guide.Guide {
	return shopping_guide.NewGuide(wc.Context)
}
