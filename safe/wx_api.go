package safe

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/define"

	"github.com/dcsunny/wechat/context"
	error2 "github.com/dcsunny/wechat/error"
	"github.com/dcsunny/wechat/util"
)

const (
	getWxIpURL = "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token=%s"
)

type WxSafe struct {
	*context.Context
}

type IpListResult struct {
	IPList []string `json:"ip_list"`
	define.CommonError
}

func NewWxSafe(context *context.Context) *WxSafe {
	tpl := new(WxSafe)
	tpl.Context = context
	return tpl
}

func (s *WxSafe) GetWxIp() (result IpListResult, err error) {
	var accessToken string
	accessToken, err = s.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", getWxIpURL, accessToken)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	err = error2.CommonErrorHandle(result.CommonError, s.Context, "getWxIp")
	return result, err
}
