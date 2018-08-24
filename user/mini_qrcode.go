package user

import (
	"fmt"

	"bytes"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/util"
)

type MiniQrCode struct {
	*context.Context
}

const (
	miniForeverRoundCodeURL  = "https://api.weixin.qq.com/wxa/getwxacode?access_token=%s"
	miniForeverSquareCodeURL = "https://api.weixin.qq.com/cgi-bin/wxaapp/createwxaqrcode?access_token=%s"
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

func (qr *MiniQrCode) CreateMiniForeverRoundCode(params *MiniCodeParams) (reader *bytes.Reader, err error) {
	var accessToken string
	accessToken, err = qr.GetMiniAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(miniForeverRoundCodeURL, accessToken)
	response, err := util.PostJSON(uri, params)
	if err != nil {
		fmt.Println(err)
		return
	}

	return bytes.NewReader(response), nil
}

func (qr *MiniQrCode) CreateMiniForeverSquareCode(params *MiniCodeParams) (reader *bytes.Reader, err error) {
	var accessToken string
	accessToken, err = qr.GetMiniAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(miniForeverSquareCodeURL, accessToken)
	response, err := util.PostJSON(uri, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	return bytes.NewReader(response), nil
}
