package context

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/dcsunny/wechat/define"
	"github.com/dcsunny/wechat/util"
)

const (
	//AccessTokenURL 获取access_token的接口
	AccessTokenURL = "/cgi-bin/token"
)

//ResAccessToken struct
type ResAccessToken struct {
	define.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

//GetAccessTokenFunc 获取 access token 的函数签名
type GetAccessTokenFunc func(ctx *Context) (accessToken string, err error)

//SetAccessTokenLock 设置读写锁（一个appID一个读写锁）
func (ctx *Context) SetAccessTokenLock(l *sync.RWMutex) {
	ctx.accessTokenLock = l
}

//SetGetAccessTokenFunc 设置自定义获取accessToken的方式, 需要自己实现缓存
func (ctx *Context) SetGetAccessTokenFunc(f GetAccessTokenFunc) {
	ctx.accessTokenFunc = f
}

//GetAccessToken 获取access_token
func (ctx *Context) GetAccessToken() (accessToken string, err error) {
	ctx.accessTokenLock.Lock()
	defer ctx.accessTokenLock.Unlock()
	if ctx.accessTokenFunc != nil {
		return ctx.accessTokenFunc(ctx)
	}
	accessTokenCacheKey := fmt.Sprintf(define.AccessTokenCacheKey, ctx.AppID)
	accessToken = ctx.Cache.GetString(accessTokenCacheKey)
	if accessToken != "" {
		return
	}
	//从微信服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = ctx.GetAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = resAccessToken.AccessToken
	return
}

//GetAccessTokenFromServer 强制从微信服务器获取token
func (ctx *Context) GetAccessTokenFromServer() (resAccessToken ResAccessToken, err error) {
	accessTokenUrl := ctx.ApiBaseUrl + AccessTokenURL
	if ctx.AccessTokenURL != "" {
		accessTokenUrl = ctx.AccessTokenURL
	}
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", accessTokenUrl, ctx.AppID, ctx.AppSecret)
	var body []byte
	body, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrMsg != "" {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	accessTokenCacheKey := fmt.Sprintf(define.AccessTokenCacheKey, ctx.AppID)
	expires := resAccessToken.ExpiresIn - 1500
	err = ctx.Cache.SetString(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)
	return
}
