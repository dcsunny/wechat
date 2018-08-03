package user

import (
	"fmt"

	"encoding/json"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/util"
	"github.com/qiniu/x/url.v7"
)

const (
	qrCodeURL          = `https://api.weixin.qq.com/cgi-bin/qrcode/create`
	qrCodeTicketURL    = `https://mp.weixin.qq.com/cgi-bin/showqrcode`
	QR_SCENE           = "QR_SCENE"
	QR_STR_SCENE       = "QR_STR_SCENE"
	QR_LIMIT_SCENE     = "QR_LIMIT_SCENE"
	QR_LIMIT_STR_SCENE = "QR_LIMIT_STR_SCENE"
)

type QrCode struct {
	*context.Context
}

type QrCodeParams struct {
	ExpireSeconds int64            `json:"expire_seconds"`
	ActionName    string           `json:"action_name"`
	ActionInfo    QrCodeActionInfo `json:"action_info"`
}

type QrCodeActionInfo struct {
	Scene QrCodeScene `json:"scene"`
}

type QrCodeScene struct {
	SceneID  int64  `json:"scene_id"`
	SceneStr string `json:"scene_str"`
}
type QrCodeRet struct {
	util.CommonError
	Ticket        string `json:"ticket"`
	ExpireSeconds int64  `json:"expire_seconds"`
	Url           string `json:"url"`
}

func NewQrCode(context *context.Context) *QrCode {
	qr := new(QrCode)
	qr.Context = context
	return qr
}

func (qr *QrCode) Create(params *QrCodeParams) (result QrCodeRet, err error) {
	var accessToken string
	accessToken, err = qr.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", qrCodeURL, accessToken)
	response, err := util.PostJSON(uri, params)
	err = json.Unmarshal(response, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("qrcode create error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

func (qr *QrCode) CreateTmp(code int64, expire int64) (result QrCodeRet, err error) {
	return qr.Create(&QrCodeParams{
		ExpireSeconds: expire,
		ActionName:    QR_LIMIT_SCENE,
		ActionInfo: QrCodeActionInfo{
			Scene: QrCodeScene{
				SceneID: code,
			},
		}})
}

func (qr *QrCode) CreateTmpStr(code string, expire int64) (result QrCodeRet, err error) {
	return qr.Create(&QrCodeParams{
		ExpireSeconds: expire,
		ActionName:    QR_LIMIT_STR_SCENE,
		ActionInfo: QrCodeActionInfo{
			Scene: QrCodeScene{
				SceneStr: code,
			},
		}})
}
func (qr *QrCode) CreateForever(code int64) (result QrCodeRet, err error) {
	return qr.Create(&QrCodeParams{
		ActionName: QR_SCENE,
		ActionInfo: QrCodeActionInfo{
			Scene: QrCodeScene{
				SceneID: code,
			},
		}})
}
func (qr *QrCode) CreateForeverStr(code string) (result QrCodeRet, err error) {
	return qr.Create(&QrCodeParams{
		ActionName: QR_SCENE,
		ActionInfo: QrCodeActionInfo{
			Scene: QrCodeScene{
				SceneStr: code,
			},
		}})
}

func (qr *QrCode) QrCodePictureUrl(ticket string) string {
	showUrl := "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s"
	resultUrl := fmt.Sprintf(showUrl, url.Escape(ticket))
	return resultUrl
}
