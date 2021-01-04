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
	newGuideTagOption       = "https://api.weixin.qq.com/cgi-bin/guide/newguidetagoption"
	delguidetagoption       = "https://api.weixin.qq.com/cgi-bin/guide/delguidetagoption"
	addGuideTagOption       = "https://api.weixin.qq.com/cgi-bin/guide/addguidetagoption"
	getGuideTagOption       = "https://api.weixin.qq.com/cgi-bin/guide/getguidetagoption"
	addGuideBuyerTag        = "https://api.weixin.qq.com/cgi-bin/guide/addguidebuyertag"
	getGuideBuyerTag        = "https://api.weixin.qq.com/cgi-bin/guide/getguidebuyertag"
	queryGuideBuyerByTag    = "https://api.weixin.qq.com/cgi-bin/guide/queryguidebuyerbytag"
	delGuideBuyerTag        = "https://api.weixin.qq.com/cgi-bin/guide/delguidebuyertag"
	addGuideBuyerDisplayTag = "https://api.weixin.qq.com/cgi-bin/guide/addguidebuyerdisplaytag"
	getGuideBuyerDisplayTag = "https://api.weixin.qq.com/cgi-bin/guide/getguidebuyerdisplaytag"
)

type GuideTag struct {
	*context.Context
}

func NewGuideTag(ctx *context.Context) *GuideTag {
	return &GuideTag{ctx}
}

type GuideTagInfo struct {
	TagName   string   `json:"tag_name"`
	TagValues []string `json:"tag_values"`
}

func (g *GuideTag) Add(req GuideTagInfo) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", newGuideTagOption, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "newGuideTagOption")
}

func (g *GuideTag) Delete(tagName string) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", delguidetagoption, accessToken)
	response, err := util.PostJSON(uri, GuideTagInfo{TagName: tagName})
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "delguidetagoption")
}

func (g *GuideTag) AddGuideTagOption(req GuideTagInfo) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", addGuideTagOption, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "addGuideTagOption")
}

type GetGuideTagOptionResp struct {
	Options []GuideTagInfo `json:"options"`
	define.CommonError
}

func (g *GuideTag) GetGuideTagOption() ([]GuideTagInfo, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", getGuideTagOption, accessToken)
	response, err := util.PostJSON(uri, nil)
	if err != nil {
		return nil, err
	}
	var result GetGuideTagOptionResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "getGuideTagOption")
	}
	return result.Options, nil
}

type AddGuideBuyerTagReq struct {
	GuideAccount string   `json:"guide_account"`
	GuideOpenid  string   `json:"guide_openid"`
	Openid       string   `json:"openid"`
	OpenidList   []string `json:"openid_list"`
	TagValue     string   `json:"tag_value"`
}

func (g *GuideTag) AddGuideBuyerTag(req AddGuideBuyerTagReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", addGuideBuyerTag, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "addGuideBuyerTag")
}

type GetGuideBuyerTagReq struct {
	GuideAccount string `json:"guide_account"`
	GuideOpenid  string `json:"guide_openid"`
	Openid       string `json:"openid"`
}

type GetGuideBuyerTagResp struct {
	TagValues []string `json:"tag_values"`
	define.CommonError
}

func (g *GuideTag) GetGuideBuyerTag(req GetGuideBuyerTagReq) ([]string, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", getGuideBuyerTag, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return nil, err
	}
	var result GetGuideBuyerTagResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "getGuideBuyerTag")
	}
	return result.TagValues, nil
}

type QueryGuideBuyerByTagReq struct {
	GuideAccount string   `json:"guide_account"`
	GuideOpenid  string   `json:"guide_openid"`
	PushCount    int64    `json:"push_count"`
	TagValues    []string `json:"tag_values"`
}

type QueryGuideBuyerByTagResp struct {
	OpenidList []string `json:"openid_list"`
	define.CommonError
}

func (g *GuideTag) QueryGuideBuyerByTag(req QueryGuideBuyerByTagReq) ([]string, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", queryGuideBuyerByTag, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return nil, err
	}
	var result QueryGuideBuyerByTagResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "queryGuideBuyerByTag")
	}
	return result.OpenidList, nil
}

type DelGuideBuyerTagReq struct {
	GuideAccount string   `json:"guide_account"`
	GuideOpenid  string   `json:"guide_openid"`
	Openid       string   `json:"openid"`
	OpenidList   []string `json:"openid_list"`
	TagValue     string   `json:"tag_value"`
}

func (g *GuideTag) DelGuideBuyerTag(req DelGuideBuyerTagReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", delGuideBuyerTag, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "delGuideBuyerTag")
}

type AddGuideBuyerDisplayTagReq struct {
	GuideAccount   string   `json:"guide_account"`
	GuideOpenid    string   `json:"guide_openid"`
	Openid         string   `json:"openid"`
	DisplayTagList []string `json:"display_tag_list"`
}

func (g *GuideTag) AddGuideBuyerDisplayTag(req AddGuideBuyerDisplayTagReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", addGuideBuyerDisplayTag, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "addGuideBuyerDisplayTag")
}

type GetGuideBuyerDisplayTagReq struct {
	GuideAccount string `json:"guide_account"`
	GuideOpenid  string `json:"guide_openid"`
	Openid       string `json:"openid"`
}

type GetGuideBuyerDisplayTagResp struct {
	DisplayTagList []string `json:"display_tag_list"`
	define.CommonError
}

func (g *GuideTag) GetGuideBuyerDisplayTag(req GetGuideBuyerDisplayTagReq) ([]string, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", getGuideBuyerDisplayTag, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return nil, err
	}
	var result GetGuideBuyerDisplayTagResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "getGuideBuyerDisplayTag")
	}
	return result.DisplayTagList, nil
}
