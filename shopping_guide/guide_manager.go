package shopping_guide

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/define"

	"github.com/dcsunny/wechat/common_error"
	"github.com/dcsunny/wechat/util"

	"github.com/dcsunny/wechat/context"
)

const (
	GuideManagerAcctAdd          = "https://api.weixin.qq.com/cgi-bin/guide/addguideacct"
	GuideManagerAcctGet          = "https://api.weixin.qq.com/cgi-bin/guide/getguideacct"
	GuideManagerAcctUpdate       = "https://api.weixin.qq.com/cgi-bin/guide/updateguideacct"
	GuideManagerAcctDelete       = "https://api.weixin.qq.com/cgi-bin/guide/delguideacct"
	GuideManagerAcctList         = "https://api.weixin.qq.com/cgi-bin/guide/getguideacctlist"
	GuideManagerAcctCreateQrCode = "https://api.weixin.qq.com/cgi-bin/guide/guidecreateqrcode"
)

type GuideManager struct {
	*context.Context
}

func NewGuideManager(ctx *context.Context) *GuideManager {
	return &GuideManager{ctx}
}

type GuideManagerAddReq struct {
	Account    string `json:"guide_account"`    //顾问微信号（guide_account和guide_openid二选一）
	OpenID     string `json:"guide_openid"`     //顾问openid或者unionid（guide_account和guide_openid二选一）
	HeadImgURL string `json:"guide_headimgurl"` //顾问头像，头像url只能用《上传图文消息内的图片获取URL》
	Nickname   string `json:"guide_nickname"`
}

type GuideManagerAcctReq struct {
	Account string `json:"guide_account"` //顾问微信号（guide_account和guide_openid二选一）
	OpenID  string `json:"guide_openid"`  //顾问openid或者unionid（guide_account和guide_openid二选一）
}

type GuideAcctInfo struct {
	Account    string `json:"guide_account"`    //顾问微信号（guide_account和guide_openid二选一）
	OpenID     string `json:"guide_openid"`     //顾问openid或者unionid（guide_account和guide_openid二选一）
	HeadImgURL string `json:"guide_headimgurl"` //顾问头像，头像url只能用《上传图文消息内的图片获取URL》
	Nickname   string `json:"guide_nickname"`
	Status     int    `json:"status"` //顾问状态（1:确认中；2已确认；3已拒绝；4已过期）
	define.CommonError
}

//为服务号添加顾问
func (g *GuideManager) AddGuideAcct(req GuideManagerAddReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideManagerAcctAdd, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "AddGuideAcct")
}

//获取顾问信息
func (g *GuideManager) GetGuideAcct(req GuideManagerAcctReq) (*GuideAcctInfo, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideManagerAcctGet, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return nil, err
	}
	var result GuideAcctInfo
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "GetGuideAcct")
	}
	return &result, nil
}

//修改顾问的昵称或头像
func (g *GuideManager) UpdateGuideAcct(req GuideManagerAddReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideManagerAcctUpdate, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "UpdateGuideAcct")
}

//删除顾问
func (g *GuideManager) DeleteGuideAcct(req GuideManagerAcctReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideManagerAcctDelete, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "UpdateGuideAcct")
}

type GuideAcctList struct {
	TotalNum int             `json:"total_num"`
	List     []GuideAcctInfo `json:"list"`
	define.CommonError
}

//获取服务号顾问列表 page从0开始
func (g *GuideManager) GetGuideAcctList(page int, num int) (GuideAcctList, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return GuideAcctList{}, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideManagerAcctList, accessToken)
	response, err := util.PostJSON(uri, map[string]interface{}{
		"page": page,
		"num":  num,
	})
	if err != nil {
		return GuideAcctList{}, err
	}
	var result GuideAcctList
	err = json.Unmarshal(response, &result)
	if err != nil {
		return GuideAcctList{}, err
	}
	if result.ErrCode != 0 {
		return GuideAcctList{}, common_error.CommonErrorHandle(result.CommonError, g.Context, "GetGuideAcctList")
	}
	return result, nil
}

type GuideCreateQrCodeReq struct {
	GuideManagerAcctReq
	QrcodeInfo string `json:"qrcode_info"`
}

type GuideCreateQrCodeResp struct {
	QrcodeURL string `json:"qrcode_url"`
	define.CommonError
}

//生成顾问二维码
func (g *GuideManager) CreateQrCode(req GuideCreateQrCodeReq) (string, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return "", err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideManagerAcctCreateQrCode, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return "", err
	}
	var result GuideCreateQrCodeResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", common_error.CommonErrorHandle(result.CommonError, g.Context, "CreateQrCode")
	}
	return result.QrcodeURL, nil
}
