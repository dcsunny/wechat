package miniprogram

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/define"

	"github.com/dcsunny/wechat/common_error"
	"github.com/dcsunny/wechat/util"
)

const (
	urlschemeGenerateUrl = "https://api.weixin.qq.com/wxa/generatescheme"
)

type UrlschemeGenerateReq struct {
	JumpWxa    UrlschemeGenerateJumWxa `json:"jump_wxa"`
	IsExpire   bool                    `json:"is_expire"`
	ExpireTime int64                   `json:"expire_time"`
}

type UrlschemeGenerateJumWxa struct {
	Path  string `json:"path"`
	Query string `json:"query"`
}

type UrlschemeGenerateResp struct {
	define.CommonError
	Openlink string `json:"openlink"`
}

func (wxa *MiniProgram) UrlschemeGenerate(req UrlschemeGenerateReq) (UrlschemeGenerateResp, error) {
	accessToken, err := wxa.GetAccessToken()
	if err != nil {
		return UrlschemeGenerateResp{}, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", urlschemeGenerateUrl, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return UrlschemeGenerateResp{}, err
	}
	var result UrlschemeGenerateResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return UrlschemeGenerateResp{}, err
	}
	if result.ErrCode != 0 {
		return UrlschemeGenerateResp{}, common_error.CommonErrorHandle(result.CommonError, wxa.Context, "urlschemeGenerate")
	}
	return result, nil
}
