package pay

import (
	"encoding/xml"
	"errors"
	"fmt"

	"time"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/util"
)

var payGateway = "https://api.mch.weixin.qq.com/pay/unifiedorder"

// Pay struct extends context
type Pay struct {
	*context.Context
}

// 传入的参数，用于生成 prepay_id 的必需参数
// PayParams was NEEDED when request unifiedorder
type PayParams struct {
	TotalFee   string
	CreateIP   string
	Body       string
	OutTradeNo string
	OpenID     string
}

// PayConfig 是传出用于 jsdk 用的参数
type PayConfig struct {
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

// payResult 是 unifie order 接口的返回
type payResult struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppID      string `xml:"appid,omitempty"`
	MchID      string `xml:"mch_id,omitempty"`
	NonceStr   string `xml:"nonce_str,omitempty"`
	Sign       string `xml:"sign,omitempty"`
	ResultCode string `xml:"result_code,omitempty"`
	TradeType  string `xml:"trade_type,omitempty"`
	PrePayID   string `xml:"prepay_id,omitempty"`
	CodeURL    string `xml:"code_url,omitempty"`
	ErrCode    string `xml:"err_code,omitempty"`
	ErrCodeDes string `xml:"err_code_des,omitempty"`
}

//payRequest 接口请求参数
type payRequest struct {
	AppID          string `xml:"appid"`
	MchID          string `xml:"mch_id"`
	DeviceInfo     string `xml:"device_info,omitempty"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	SignType       string `xml:"sign_type,omitempty"`
	Body           string `xml:"body"`
	Detail         string `xml:"detail,omitempty"`
	Attach         string `xml:"attach,omitempty"`      //附加数据
	OutTradeNo     string `xml:"out_trade_no"`          //商户订单号
	FeeType        string `xml:"fee_type,omitempty"`    //标价币种
	TotalFee       string `xml:"total_fee"`             //标价金额
	SpbillCreateIp string `xml:"spbill_create_ip"`      //终端IP
	TimeStart      string `xml:"time_start,omitempty"`  //交易起始时间
	TimeExpire     string `xml:"time_expire,omitempty"` //交易结束时间
	GoodsTag       string `xml:"goods_tag,omitempty"`   //订单优惠标记
	NotifyUrl      string `xml:"notify_url"`            //通知地址
	TradeType      string `xml:"trade_type"`            //交易类型
	ProductId      string `xml:"product_id,omitempty"`  //商品ID
	LimitPay       string `xml:"limit_pay,omitempty"`   //
	OpenID         string `xml:"openid,omitempty"`      //用户标识
	SceneInfo      string `xml:"scene_info,omitempty"`  //场景信息
}

type NotifyResult struct {
	ReturnCode         string `xml:"return_code"`
	ReturnMsg          string `xml:"return_msg"`
	Appid              string `xml:"appid"`
	MchID              string `xml:"mch_id"`
	DeviceInfo         string `xml:"device_info"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	SignType           string `xml:"sign_type"`
	ResultCode         string `xml:"result_code"`
	ErrCode            string `xml:"err_code"`
	ErrCodeDes         string `xml:"err_code_des"`
	Openid             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	BankType           string `xml:"bank_type"`
	TotalFee           string `xml:"total_fee"`
	SettlementTotalFee string `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            string `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          string `xml:"coupon_fee"`
	CouponCount        string `xml:"coupon_count"`
	TransactionId      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	Attach             string `xml:"attach"`
	TimeEnd            string `xml:"time_end"`
}

type NotifyReturn struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// NewPay return an instance of Pay package
func NewPay(ctx *context.Context) *Pay {
	pay := Pay{Context: ctx}
	return &pay
}

// PrePayId will request wechat merchant api and request for a pre payment order id
func (pcf *Pay) PrePayId(p *PayParams) (prePayID string, err error) {
	nonceStr := util.RandomStr(32)
	tradeType := "JSAPI"
	template := "appid=%s&body=%s&mch_id=%s&nonce_str=%s"
	if pcf.PayNotifyURL != "" {
		template = template + fmt.Sprintf("&notify_url=%s", pcf.PayNotifyURL)
	}
	template = template + "&openid=%s&out_trade_no=%s&spbill_create_ip=%s&total_fee=%s&trade_type=%s&key=%s"
	str := fmt.Sprintf(template, pcf.AppID, p.Body, pcf.PayMchID, nonceStr, p.OpenID, p.OutTradeNo, p.CreateIP, p.TotalFee, tradeType, pcf.PayKey)
	sign := util.MD5Sum(str)
	request := payRequest{
		AppID:          pcf.AppID,
		MchID:          pcf.PayMchID,
		NonceStr:       nonceStr,
		Sign:           sign,
		Body:           p.Body,
		OutTradeNo:     p.OutTradeNo,
		TotalFee:       p.TotalFee,
		SpbillCreateIp: p.CreateIP,
		NotifyUrl:      pcf.PayNotifyURL,
		TradeType:      tradeType,
		OpenID:         p.OpenID,
	}
	rawRet, err := util.PostXML(payGateway, request)
	if err != nil {
		return "", errors.New(err.Error() + " parameters : " + str)
	}
	payRet := payResult{}
	err = xml.Unmarshal(rawRet, &payRet)
	if err != nil {
		return "", errors.New(err.Error())
	}
	if payRet.ReturnCode == "SUCCESS" {
		//pay success
		if payRet.ResultCode == "SUCCESS" {
			return payRet.PrePayID, nil
		}
		return "", errors.New(payRet.ErrCode + payRet.ErrCodeDes)
	} else {
		return "", errors.New("[msg : xmlUnmarshalError] [rawReturn : " + string(rawRet) + "] [params : " + str + "] [sign : " + sign + "]")
	}
}

func (pcf *Pay) JSPayParams(prePayID string) PayConfig {
	payConf := PayConfig{
		AppID:     pcf.AppID,
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
		NonceStr:  util.RandomStr(32),
		Package:   fmt.Sprintf("prepay_id=%s", prePayID),
		SignType:  "MD5",
	}
	str := fmt.Sprintf("appId=%s&nonceStr=%s&package=%s&signType=%s&timeStamp=%s&key=%s", payConf.AppID, payConf.NonceStr, payConf.Package, payConf.SignType, payConf.TimeStamp, pcf.PayKey)
	payConf.Sign = util.MD5Sum(str)
	return payConf
}
