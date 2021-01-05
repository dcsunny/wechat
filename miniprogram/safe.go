package miniprogram

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/common_error"
	"github.com/dcsunny/wechat/define"
	"github.com/dcsunny/wechat/util"
)

const (
	getUserRiskRank = "https://api.weixin.qq.com/wxa/getuserriskrank"
)

type UserRiskRankReq struct {
	AppID        string `json:"appid"`
	BankCardNo   string `json:"bank_card_no"`
	CertNo       string `json:"cert_no"`
	ClientIP     string `json:"client_ip"`
	EmailAddress string `json:"email_address"`
	ExtendedInfo string `json:"extended_info"`
	MobileNo     string `json:"mobile_no"`
	OpenID       string `json:"openid"`
	Scene        int    `json:"scene"`
}

type UserRiskRankResp struct {
	RiskRank int   `json:"risk_rank"`
	UnoinID  int64 `json:"unoin_id"` //唯一请求标识，标记单次请求
	define.CommonError
}

//根据提交的用户信息数据获取用户的安全等级 risk_rank，无需用户授权。
func (wxa *MiniProgram) GetUserRiskRank(req UserRiskRankReq) (UserRiskRankResp, error) {
	accessToken, err := wxa.GetAccessToken()
	if err != nil {
		return UserRiskRankResp{}, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", getUserRiskRank, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return UserRiskRankResp{}, err
	}
	var result UserRiskRankResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return UserRiskRankResp{}, err
	}
	if result.ErrCode != 0 {
		return UserRiskRankResp{}, common_error.CommonErrorHandle(result.CommonError, wxa.Context, "getUserRiskRank")
	}
	return result, nil
}
