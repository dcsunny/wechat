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
	miniAccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

type MiniAccessToken struct {
	define.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (ctx *Context) SetMiniAccessTokenLock(l *sync.RWMutex) {
	ctx.accessTokenLock = l
}

//GetQyAccessToken 获取access_token
func (ctx *Context) GetMiniAccessToken() (accessToken string, err error) {
	ctx.accessTokenLock.Lock()
	defer ctx.accessTokenLock.Unlock()

	accessTokenCacheKey := fmt.Sprintf("mini_access_token_%s", ctx.AppID)
	accessToken = ctx.Cache.GetString(accessTokenCacheKey)
	if accessToken != "" {
		return
	}

	//从微信服务器获取
	var resQyAccessToken ResQyAccessToken
	resQyAccessToken, err = ctx.GetMiniAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = resQyAccessToken.AccessToken
	return
}

//GetQyAccessTokenFromServer 强制从微信服务器获取token
func (ctx *Context) GetMiniAccessTokenFromServer() (resQyAccessToken ResQyAccessToken, err error) {
	url := fmt.Sprintf(miniAccessTokenURL, ctx.AppID, ctx.AppSecret)
	var body []byte
	body, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resQyAccessToken)
	if err != nil {
		return
	}
	if resQyAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get mini_access_token error : errcode=%v , errormsg=%v", resQyAccessToken.ErrCode, resQyAccessToken.ErrMsg)
		return
	}

	AccessTokenCacheKey := fmt.Sprintf("mini_access_token_%s", ctx.AppID)
	expires := resQyAccessToken.ExpiresIn - 1500
	err = ctx.Cache.SetString(AccessTokenCacheKey, resQyAccessToken.AccessToken, time.Duration(expires)*time.Second)
	return
}
