package shopping_guide

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/define"

	"github.com/dcsunny/wechat/common_error"
	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/util"
)

const (
	GuideBuyerRelationAdd        = "https://api.weixin.qq.com/cgi-bin/guide/addguidebuyerrelation"
	GuideBuyerRelationDelete     = "https://api.weixin.qq.com/cgi-bin/guide/delguidebuyerrelation"
	GuideBuyerRelationList       = "https://api.weixin.qq.com/cgi-bin/guide/getguidebuyerrelationlist"
	GuideBuyerRelationRebind     = "https://api.weixin.qq.com/cgi-bin/guide/rebindguideacctforbuyer"
	GuideBuyerUpdateNickname     = "https://api.weixin.qq.com/cgi-bin/guide/updateguidebuyerrelation"
	GetGuideBuyerRelationByBuyer = "https://api.weixin.qq.com/cgi-bin/guide/getguidebuyerrelationbybuyer"
	GetGuideBuyerRelation        = "https://api.weixin.qq.com/cgi-bin/guide/getguidebuyerrelation"
)

type GuideBuyer struct {
	*context.Context
}

func NewGuideBuyer(ctx *context.Context) *GuideBuyer {
	return &GuideBuyer{ctx}
}

type BuyerInfo struct {
	OpenID        string `json:"openid"`
	BuyerNickname string `json:"buyer_nickname"`
	CreateTime    int64  `json:"create_time"`
}

type GuideBuyerRelationAddReq struct {
	GuideAccount  string      `json:"guide_account"`
	GuideOpenID   string      `json:"guide_openid"`
	OpenID        string      `json:"openid"`
	BuyerNickname string      `json:"buyer_nickname"`
	BuyerList     []BuyerInfo `json:"buyer_list"` //客户列表（不超过200，openid和buyer_list二选一）
}

func (g *GuideBuyer) Add(req GuideBuyerRelationAddReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideBuyerRelationAdd, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "AddGuideBuyerRelation")
}

type GuideBuyerRelationDeleteReq struct {
	GuideAccount  string   `json:"guide_account"`
	GuideOpenID   string   `json:"guide_openid"`
	OpenID        string   `json:"openid"`
	BuyerNickname string   `json:"buyer_nickname"`
	BuyerList     []string `json:"buyer_list"` //客户列表（不超过200，openid和buyer_list二选一）
}

func (g *GuideBuyer) Delete(req GuideBuyerRelationDeleteReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideBuyerRelationDelete, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "DeleteGuideBuyerRelation")
}

type GuideBuyerRelationListReq struct {
	GuideAccount string `json:"guide_account"`
	GuideOpenID  string `json:"guide_openid"`
	Page         int    `json:"page"` //分页页数，从0开始
	Num          int    `json:"num"`
}

type GuideBuyerRelationListResp struct {
	TotalNum int         `json:"total_num"`
	List     []BuyerInfo `json:"list"`
	define.CommonError
}

func (g *GuideBuyer) List(req GuideBuyerRelationListReq) (*GuideBuyerRelationListResp, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideBuyerRelationList, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return nil, err
	}
	var result GuideBuyerRelationListResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "GuideBuyerRelationList")
	}
	return &result, nil
}

type GuideBuyerRelationRebindReq struct {
	NewGuideAccount string   `json:"new_guide_account"`
	NewGuideOpenid  string   `json:"new_guide_openid"`
	OldGuideAccount string   `json:"old_guide_account"`
	OldGuideOpenid  string   `json:"old_guide_openid"`
	Openid          string   `json:"openid"`
	OpenidList      []string `json:"openid_list"`
}

func (g *GuideBuyer) Rebind(req GuideBuyerRelationRebindReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideBuyerRelationRebind, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "GuideBuyerRelationRebind")
}

type BuyerUpdateNicknameReq struct {
	BuyerNickname string `json:"buyer_nickname"`
	GuideAccount  string `json:"guide_account"`
	GuideOpenid   string `json:"guide_openid"`
	Openid        string `json:"openid"`
}

func (g *GuideBuyer) UpdateNickname(req BuyerUpdateNicknameReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GuideBuyerUpdateNickname, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "BuyerUpdateNickname")
}

type GetGuideByBuyerResp struct {
	BuyerNickname string `json:"buyer_nickname"`
	CreateTime    int64  `json:"create_time"`
	GuideAccount  string `json:"guide_account"`
	GuideOpenid   string `json:"guide_openid"`
	Openid        string `json:"openid"`
	define.CommonError
}

func (g *GuideBuyer) GetGuideByBuyer(openID string) (GetGuideByBuyerResp, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return GetGuideByBuyerResp{}, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GetGuideBuyerRelationByBuyer, accessToken)
	response, err := util.PostJSON(uri, map[string]interface{}{
		"openid": openID,
	})
	if err != nil {
		return GetGuideByBuyerResp{}, err
	}
	var result GetGuideByBuyerResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return GetGuideByBuyerResp{}, err
	}
	if result.ErrCode != 0 {
		return GetGuideByBuyerResp{}, common_error.CommonErrorHandle(result.CommonError, g.Context, "GetGuideByBuyer")
	}
	return result, nil
}

type GetGuideBuyerRelationReq struct {
	GuideAccount string `json:"guide_account"`
	GuideOpenid  string `json:"guide_openid"`
	Openid       string `json:"openid"`
}

type GetGuideBuyerRelationResp struct {
	BuyerInfo
	define.CommonError
}

func (g *GuideBuyer) GetGuideBuyerRelation(req GetGuideBuyerRelationReq) (BuyerInfo, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return BuyerInfo{}, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", GetGuideBuyerRelation, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return BuyerInfo{}, err
	}
	var result GetGuideBuyerRelationResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return BuyerInfo{}, err
	}
	if result.ErrCode != 0 {
		return BuyerInfo{}, common_error.CommonErrorHandle(result.CommonError, g.Context, "GetGuideBuyerRelation")
	}
	return result.BuyerInfo, nil
}
