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
	setGuideCardMaterial  = "https://api.weixin.qq.com/cgi-bin/guide/setguidecardmaterial"
	getGuideCardMaterial  = "https://api.weixin.qq.com/cgi-bin/guide/getguidecardmaterial"
	delGuideCardMaterial  = "https://api.weixin.qq.com/cgi-bin/guide/delguidecardmaterial"
	setGuideImageMaterial = "https://api.weixin.qq.com/cgi-bin/guide/setguideimagematerial"
	getGuideImageMaterial = "https://api.weixin.qq.com/cgi-bin/guide/getguideimagematerial"
	delGuideImageMaterial = "https://api.weixin.qq.com/cgi-bin/guide/delguideimagematerial"
	setGuideWordMaterial  = "https://api.weixin.qq.com/cgi-bin/guide/setguidewordmaterial"
	getGuideWordMaterial  = "https://api.weixin.qq.com/cgi-bin/guide/getguidewordmaterial"
	delGuideWordMaterial  = "https://api.weixin.qq.com/cgi-bin/guide/delguidewordmaterial"
)

type GuideMaterial struct {
	*context.Context
}

func NewGuideMaterial(ctx *context.Context) *GuideMaterial {
	return &GuideMaterial{ctx}
}

type SetGuideCardMaterialReq struct {
	AppID   string `json:"appid"`
	MediaID string `json:"media_id"`
	Path    string `json:"path"`
	Title   string `json:"title"`
	Type    int    `json:"type"`
}

func (g *GuideMaterial) SetGuideCardMaterial(req SetGuideCardMaterialReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", setGuideCardMaterial, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "setGuideCardMaterial")
}

type GetGuideCardMaterialResp struct {
	CardList []GuideCardInfo `json:"card_list"`
	define.CommonError
}

type GuideCardInfo struct {
	AppID    string `json:"appid"`
	MasterID int64  `json:"master_id"`
	Path     string `json:"path"`
	PicURL   string `json:"picurl"`
	SlaveID  int64  `json:"slave_id"`
	Title    string `json:"title"`
}

func (g *GuideMaterial) GetGuideCardMaterial(_type int) ([]GuideCardInfo, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", getGuideCardMaterial, accessToken)
	response, err := util.PostJSON(uri, map[string]interface{}{
		"type": _type,
	})
	if err != nil {
		return nil, err
	}
	var result GetGuideCardMaterialResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "getGuideCardMaterial")
	}
	return result.CardList, nil
}

type DelGuideCardMaterialReq struct {
	AppID string `json:"appid"`
	Path  string `json:"path"`
	Title string `json:"title"`
	Type  int    `json:"type"`
}

func (g *GuideMaterial) DelGuideCardMaterial(req DelGuideCardMaterialReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", delGuideCardMaterial, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "delGuideCardMaterial")
}

type SetGuideImageMaterialReq struct {
	AppID   string `json:"appid"`
	MediaID string `json:"media_id"`
}

func (g *GuideMaterial) SetGuideImageMaterial(req SetGuideImageMaterialReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", setGuideImageMaterial, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "setGuideImageMaterial")
}

type GetGuideImageMaterialResp struct {
	ModelList []struct {
		Picurl string `json:"picurl"`
	} `json:"model_list"`
	TotalNum int64 `json:"total_num"`
	define.CommonError
}

func (g *GuideMaterial) GetGuideImageMaterial(_type int, start, num int) (*GetGuideImageMaterialResp, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", getGuideImageMaterial, accessToken)
	response, err := util.PostJSON(uri, map[string]interface{}{
		"type":  _type,
		"start": start,
		"num":   num,
	})
	if err != nil {
		return nil, err
	}
	var result GetGuideImageMaterialResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "getGuideImageMaterial")
	}
	return &result, nil
}

type DelGuideImageMaterialReq struct {
	Type   int    `json:"type"`
	PicURL string `json:"picurl"`
}

func (g *GuideMaterial) DelGuideImageMaterial(req DelGuideBuyerTagReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", delGuideImageMaterial, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "delGuideImageMaterial")
}

type SetOrDelGuideWordMaterialReq struct {
	Type int    `json:"type"`
	Word string `json:"word"`
}

//文字素材最多支持 300 字
func (g *GuideMaterial) SetGuideWordMaterial(req SetOrDelGuideWordMaterialReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", setGuideWordMaterial, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "setGuideWordMaterial")
}

type GetGuideWordMaterialResp struct {
	TotalNum int64                   `json:"total_num"`
	WordList []GuideWordMaterialInfo `json:"word_list"`
	define.CommonError
}

type GuideWordMaterialInfo struct {
	CreateTime int64  `json:"create_time"`
	Word       string `json:"word"`
}

func (g *GuideMaterial) GetGuideWordMaterial(_type int, start, num int) (*GetGuideWordMaterialResp, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", getGuideWordMaterial, accessToken)
	response, err := util.PostJSON(uri, map[string]interface{}{
		"type":  _type,
		"start": start,
		"num":   num,
	})
	if err != nil {
		return nil, err
	}
	var result GetGuideWordMaterialResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "getGuideWordMaterial")
	}
	return &result, nil
}

func (g *GuideMaterial) DelGuideWordMaterial(req SetOrDelGuideWordMaterialReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", delGuideWordMaterial, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "delGuideWordMaterial")
}
