package oauth

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/define"
	"github.com/dcsunny/wechat/util"
)

type MiniSession struct {
	define.CommonError
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid"`
	SessionKey string `json:"session_key"`
}

const (
	Jscode2SessionURL = "/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func (oauth *Oauth) Jscode2Session(code string) (session MiniSession, err error) {
	urlStr := fmt.Sprintf(oauth.ApiBaseUrl+Jscode2SessionURL, oauth.AppID, oauth.AppSecret, code)
	var response []byte
	response, err = util.HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &session)
	if err != nil {
		return
	}
	if session.ErrCode != 0 {
		err = fmt.Errorf("get user session key error : errcode=%v , errmsg=%v", session.ErrCode, session.ErrMsg)
		return
	}
	return
}
