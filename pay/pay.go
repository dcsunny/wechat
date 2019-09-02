package pay

import (
	"encoding/xml"
	"errors"
	"fmt"

	"time"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/util"
)

const (
	payGateway  = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	mchTransUri = "https://api.qpay.qq.com/cgi-bin/epay/qpay_epay_b2c.cgi"
	sendRedUri  = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack"
	refundUri   = "https://api.mch.weixin.qq.com/secapi/pay/refund"
)

// Pay struct extends context
type Pay struct {
	*context.Context
}

// 传入的参数，用于生成 prepay_id 的必需参数
// PayParams was NEEDED when request unifiedorder
type Params struct {
	TotalFee   int
	CreateIP   string
	Body       string
	OutTradeNo string
	OpenID     string
	TradeType  string
	//以下红包使用
	Wishing  string
	SendName string
	ActName  string
	Remark   string
	SceneID  string
}

// PayConfig 是传出用于 jsdk 用的参数
type PayConfig struct {
	AppID     string `xml:"appId" json:"appId"`
	TimeStamp string `xml:"timeStamp" json:"timeStamp"`
	NonceStr  string `xml:"nonceStr" json:"nonceStr"`
	Package   string `xml:"package" json:"package"`
	SignType  string `xml:"signType" json:"signType"`
	PaySign   string `xml:"paySign" json:"paySign"`
}

type PreOrder struct {
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

type AppPayConfig struct {
	AppID     string `xml:"appid" json:"appid"`
	PartnerID string `xml:"partnerid" json:"partnerid"`
	PrePayID  string `xml:"prepayid" json:"prepayid"`
	Package   string `xml:"package" json:"package"`
	NonceStr  string `xml:"noncestr" json:"noncestr"`
	Timestamp string `xml:"timestamp" json:"timestamp"`
	Sign      string `xml:"sign" json:"sign"`
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
	TotalFee       int    `xml:"total_fee"`             //标价金额
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

func (pcf *Pay) PrePayIdByJs(p *Params) (prePayID string, err error) {
	p.TradeType = "JSAPI"
	return pcf.PrePayId(p)
}

func (pcf *Pay) PrePayIdByApp(p *Params) (prePayID string, err error) {
	p.TradeType = "APP"
	return pcf.PrePayId(p)
}

func (pcf *Pay) PrePayOrderByJs(p *Params) (payOrder PreOrder, err error) {
	p.TradeType = "JSAPI"
	return pcf.PrePayOrder(p)
}

func (pcf *Pay) PrePayOrderByApp(p *Params) (payOrder PreOrder, err error) {
	p.TradeType = "APP"
	return pcf.PrePayOrder(p)
}

func (pcf *Pay) PrePayOrder(p *Params) (payOrder PreOrder, err error) {
	nonceStr := util.RandomStr(32)

	request := payRequest{
		AppID:          pcf.AppID,
		MchID:          pcf.PayMchID,
		NonceStr:       nonceStr,
		Body:           p.Body,
		OutTradeNo:     p.OutTradeNo,
		TotalFee:       p.TotalFee,
		SpbillCreateIp: p.CreateIP,
		NotifyUrl:      pcf.PayNotifyURL,
		TradeType:      p.TradeType,
		OpenID:         p.OpenID,
	}
	sign, err := pcf.Sign(&request, pcf.PayKey)
	if err != nil {
		fmt.Println(err)
		return payOrder, err
	}
	request.Sign = sign
	rawRet, err := util.PostXML(payGateway, request, "payRequest", nil)
	if err != nil {
		return PreOrder{}, errors.New(err.Error())
	}
	err = xml.Unmarshal(rawRet, &payOrder)
	if err != nil {
		return payOrder, errors.New(err.Error())
	}
	if payOrder.ReturnCode == "SUCCESS" {
		//pay success
		if payOrder.ResultCode == "SUCCESS" {
			return payOrder, nil
		}
		return payOrder, errors.New(payOrder.ErrCode + payOrder.ErrCodeDes)
	} else {
		return payOrder, errors.New("[msg : xmlUnmarshalError] [rawReturn : " + string(rawRet) + "] [sign : " + sign + "]")
	}
}

// PrePayId will request wechat merchant api and request for a pre payment order id
func (pcf *Pay) PrePayId(p *Params) (prePayID string, err error) {
	order, err := pcf.PrePayOrder(p)
	if err != nil {
		return
	}
	if order.PrePayID == "" {
		err = errors.New("empty prepayid")
	}
	prePayID = order.PrePayID
	return
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
	payConf.PaySign = util.MD5Sum(str)
	return payConf
}

func (pcf *Pay) AppPayParams(prePayID string) AppPayConfig {
	payConf := AppPayConfig{
		AppID:     pcf.AppID,
		PartnerID: pcf.PayMchID,
		PrePayID:  prePayID,
		Package:   "Sign=WXPay",
		NonceStr:  util.RandomStr(32),
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
	}
	sign, err := pcf.Sign(&payConf, pcf.PayKey)
	if err != nil {
		fmt.Println(err)
		return payConf
	}
	payConf.Sign = sign
	return payConf
}

func (pcf *Pay) Sign(variable interface{}, key string) (sign string, err error) {
	ss := &SignStruct{
		ToLower: false,
		Tag:     "xml",
	}
	sign, err = ss.Sign(variable, nil, key)
	return
}

type MchTransfersParams struct {
	InputCharset   string `xml:"input_charset"`
	MchID          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	OutTradeNo     string `xml:"out_trade_no"`
	TotalFee       int    `xml:"total_fee"`
	Memo           string `xml:"memo"`
	AppID          string `xml:"appid"`
	OpenID         string `xml:"openid"`
	OpUserID       string `xml:"op_user_id"`
	OpUserPasswd   string `xml:"op_user_passwd"`
	SpbillCreateIp string `xml:"spbill_create_ip"`
}

func (pcf *Pay) MchPay(p *Params) error {
	nonceStr := util.RandomStr(32)
	params := &MchTransfersParams{
		InputCharset:   "UTF-8",
		AppID:          pcf.AppID,
		OpenID:         p.OpenID,
		MchID:          pcf.PayMchID,
		OutTradeNo:     p.OutTradeNo,
		NonceStr:       nonceStr,
		TotalFee:       p.TotalFee,
		Memo:           p.Body,
		SpbillCreateIp: p.CreateIP,
		OpUserID:       pcf.PayOpUserID,
		OpUserPasswd:   pcf.PayOpUserPwd,
	}
	sign, err := pcf.Sign(params, pcf.PayKey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	params.Sign = sign
	client, err := util.NewTLSHttpClient([]byte(pcf.PayCertPEMBlock), []byte(pcf.PayKeyPEMBlock))
	if err != nil {
		return err
	}
	rawRet, err := util.PostXML(mchTransUri, params, "MchTransfersParams", client)
	if err != nil {
		fmt.Println(err)
		return err
	}
	payRet := payResult{}
	err = xml.Unmarshal(rawRet, &payRet)
	if err != nil {
		fmt.Println("xmlUnmarshalError,res:" + string(rawRet))
		return err
	}
	if payRet.ReturnCode == "SUCCESS" {
		if payRet.ResultCode == "SUCCESS" {
			return nil
		}
		return errors.New(payRet.ErrCodeDes)
	} else {
		return errors.New("[msg : xmlUnmarshalError] [rawReturn : " + string(rawRet) + "]")
	}
	return nil
}

type RedParams struct {
	NonceStr     string `xml:"nonce_str"`
	Sign         string `xml:"sign"`
	MchBillno    string `xml:"mch_billno"`
	MchID        string `xml:"mch_id"`
	WxAppID      string `xml:"wxappid"`
	SendName     string `xml:"send_name"`
	ReOpenID     string `xml:"re_openid"`
	TotalAmount  int    `xml:"total_amount"`
	TotalNum     int    `xml:"total_num"`
	Wishing      string `xml:"wishing"`
	ClientIP     string `xml:"client_ip"`
	ActName      string `xml:"act_name"`
	Remark       string `xml:"remark"`
	SceneID      string `xml:"scene_id"`
	RiskInfo     string `xml:"risk_info"`
	ConsumeMchID string `xml:"consume_mch_id"`
}

type RedResult struct {
	ReturnCode  string `xml:"return_code"`
	ReturnMsg   string `xml:"return_msg"`
	ResultCode  string `xml:"result_code"`
	ErrCode     string `xml:"err_code"`
	ErrCodeDes  string `xml:"err_code_des"`
	MchBillno   string `xml:"mch_billno"`
	MchID       string `xml:"mch_id"`
	WxAppID     string `xml:"wxappid"`
	ReOpenID    string `xml:"re_openid"`
	TotalAmount int    `xml:"total_amount"`
	SendListid  string `xml:"send_listid"`
}

func (pcf *Pay) SendRed(p *Params) error {
	nonceStr := util.RandomStr(32)
	params := &RedParams{
		NonceStr:    nonceStr,
		MchBillno:   p.OutTradeNo,
		MchID:       pcf.PayMchID,
		WxAppID:     pcf.AppID,
		SendName:    p.SendName,
		ReOpenID:    p.OpenID,
		TotalAmount: p.TotalFee,
		TotalNum:    1,
		Wishing:     p.Wishing,
		ClientIP:    p.CreateIP,
		ActName:     p.ActName,
		Remark:      p.Remark,
		SceneID:     p.SceneID,
	}
	sign, err := pcf.Sign(params, pcf.PayKey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	params.Sign = sign
	client, err := util.NewTLSHttpClient([]byte(pcf.PayCertPEMBlock), []byte(pcf.PayKeyPEMBlock))
	if err != nil {
		return err
	}
	rawRet, err := util.PostXML(sendRedUri, params, "RedParams", client)
	if err != nil {
		fmt.Println(err)
		return err
	}
	payRet := RedResult{}
	err = xml.Unmarshal(rawRet, &payRet)
	if err != nil {
		fmt.Println("xmlUnmarshalError,res:" + string(rawRet))
		return err
	}
	if payRet.ReturnCode == "SUCCESS" {
		if payRet.ResultCode == "SUCCESS" {
			return nil
		}
		return errors.New(payRet.ErrCodeDes)
	} else {
		return errors.New("[msg : xmlUnmarshalError] [rawReturn : " + string(rawRet) + "]")
	}
	return nil
}

type WxRefundParams struct {
	AppID         string  `xml:"appid"`
	MchID         string  `xml:"mch_id"`
	NonceStr      string  `xml:"nonce_str"`
	Sign          string  `xml:"sign"`
	SignType      *string `xml:"sign_type"`
	TransactionID *string `xml:"transaction_id"`
	OutTradeNo    *string `xml:"out_trade_no"`
	OutRefundNo   string  `xml:"out_refund_no"`
	TotalFee      int     `xml:"total_fee"`
	RefundFee     int     `xml:"refund_fee"`
	RefundFeeType *string `xml:"refund_fee_type"`
	RefundDesc    *string `xml:"refund_desc"`
	RefundAccount *string `xml:"refund_account"`
	NotifyUrl     *string `xml:"notify_url"`
}
type RefundParams struct {
	TotalFee      int
	RefundFee     int
	CreateIP      string
	OutRefundNo   string
	TransactionID *string
	OutTradeNo    *string
}
type RefundResult struct {
	ReturnCode          string `xml:"return_code"`
	ReturnMsg           string `xml:"return_msg"`
	ResultCode          string `xml:"result_code"`
	ErrCode             string `xml:"err_code"`
	ErrCodeDes          string `xml:"err_code_des"`
	AppID               string `xml:"appid"`
	MchID               string `xml:"mch_id"`
	NonceStr            string `xml:"nonce_str"`
	Sign                string `xml:"sign"`
	TransactionID       string `xml:"transaction_id"`
	OutTradeNo          string `xml:"out_trade_no"`
	OutRefundNo         string `xml:"out_refund_no"`
	RefundID            string `xml:"refund_id"`
	RefundFee           int    `xml:"refund_fee"`
	SettlementRefundFee int    `xml:"settlement_refund_fee"`
	TotalFee            int    `xml:"total_fee"`
	SettlementTotalFee  int    `xml:"settlement_total_fee"`
	FeeType             string `xml:"fee_type"`
	CashFee             int    `xml:"cash_fee"`
	CashFeeType         string `xml:"cash_fee_type"`
	CashRefundFee       int    `xml:"cash_refund_fee"`
}

func (pcf *Pay) Refund(p *RefundParams) error {
	nonceStr := util.RandomStr(32)
	params := &WxRefundParams{
		AppID:         pcf.AppID,
		MchID:         pcf.PayMchID,
		NonceStr:      nonceStr,
		OutRefundNo:   p.OutRefundNo,
		TotalFee:      p.TotalFee,
		RefundFee:     p.RefundFee,
		TransactionID: p.TransactionID,
		OutTradeNo:    p.OutTradeNo,
	}
	sign, err := pcf.Sign(params, pcf.PayKey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	params.Sign = sign
	client, err := util.NewTLSHttpClient([]byte(pcf.PayCertPEMBlock), []byte(pcf.PayKeyPEMBlock))
	if err != nil {
		return err
	}
	rawRet, err := util.PostXML(refundUri, params, "WxRefundParams", client)
	if err != nil {
		fmt.Println(err)
		return err
	}
	payRet := RefundResult{}
	err = xml.Unmarshal(rawRet, &payRet)
	if err != nil {
		fmt.Println("xmlUnmarshalError,res:" + string(rawRet))
		return err
	}
	if payRet.ReturnCode == "SUCCESS" {
		if payRet.ResultCode == "SUCCESS" {
			return nil
		}
		return errors.New(payRet.ErrCodeDes)
	} else {
		return errors.New("[msg : xmlUnmarshalError] [rawReturn : " + string(rawRet) + "]")
	}
	return nil
}
