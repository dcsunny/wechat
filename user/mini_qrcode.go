package user

import (
	"fmt"

	"bytes"

	"encoding/json"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/define"
	error2 "github.com/dcsunny/wechat/error"
	"github.com/dcsunny/wechat/util"
)

type MiniQrCode struct {
	*context.Context
}

const (
	miniForeverRoundCodeURL   = "https://api.weixin.qq.com/wxa/getwxacode?access_token=%s"
	miniForeverSquareCodeURL  = "https://api.weixin.qq.com/cgi-bin/wxaapp/createwxaqrcode?access_token=%s"
	miniForeverUnLimitCodeURL = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
)

type MiniCodeParams struct {
	Scene     *string     `json:"scene"`
	Path      string      `json:"path"`
	Width     int         `json:"width"`
	AutoColor bool        `json:"auto_color"`
	LineColor LineCodeRGB `json:"line_color"`
	IsHyaline bool        `json:"is_hyaline"`
}

type LineCodeRGB struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

func NewMiniQrCode(context *context.Context) *MiniQrCode {
	qr := new(MiniQrCode)
	qr.Context = context
	return qr
}

//小程序码 不带scene参数
func (qr *MiniQrCode) CreateMiniForeverRoundCode(params *MiniCodeParams) (reader *bytes.Reader, err error) {
	var accessToken string
	accessToken, err = qr.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(miniForeverRoundCodeURL, accessToken)
	response, err := util.PostJSON(uri, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return bytes.NewReader(response), nil
	}
	if result.ErrCode != 0 {
		error2.CommonErrorHandle(result, qr.Context)
		err = fmt.Errorf("qrcode create error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return bytes.NewReader(response), nil
}

//小程序二维码码 //底部有 "微信扫一扫,使用小程序"字样
func (qr *MiniQrCode) CreateMiniForeverSquareCode(params *MiniCodeParams) (reader *bytes.Reader, err error) {
	var accessToken string
	accessToken, err = qr.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(miniForeverSquareCodeURL, accessToken)
	response, err := util.PostJSON(uri, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return bytes.NewReader(response), nil
	}
	if result.ErrCode != 0 {
		error2.CommonErrorHandle(result, qr.Context)
		err = fmt.Errorf("qrcode create error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return bytes.NewReader(response), nil
}

//无限制的小程序码,必须携带scene
func (qr *MiniQrCode) CreateMiniForeverUnLimitCode(params *MiniCodeParams) (reader *bytes.Reader, err error) {
	var accessToken string
	accessToken, err = qr.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(miniForeverUnLimitCodeURL, accessToken)
	response, err := util.PostJSON(uri, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return bytes.NewReader(response), nil
	}
	if result.ErrCode != 0 {
		error2.CommonErrorHandle(result, qr.Context)
		err = fmt.Errorf("qrcode create error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return bytes.NewReader(response), nil
}
