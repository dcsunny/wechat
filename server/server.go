package server

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime/debug"
	"strconv"

	"strings"

	"time"

	"net/http"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/message"
	"github.com/dcsunny/wechat/util"
)

//Server struct
type Server struct {
	*context.Context

	openID string

	messageHandler func(message.MixMessage) *message.Reply

	requestRawXMLMsg    []byte
	requestMsg          message.MixMessage
	responseRawXMLMsg   []byte
	responseNeedForward bool
	responseMsg         interface{}

	isSafeMode bool
	random     []byte
	nonce      string
	timestamp  int64

	mssageForwardUrl    string
	messageForwardToken string
}

//NewServer init
func NewServer(context *context.Context) *Server {
	srv := new(Server)
	srv.Context = context
	return srv
}

//Serve 处理微信的请求消息
func (srv *Server) Serve() error {
	if !srv.Validate() {
		return fmt.Errorf("请求校验失败")
	}

	echostr, exists := srv.GetQuery("echostr")
	if exists {
		srv.String(echostr)
		return nil
	}

	response, err := srv.handleRequest()
	if err != nil {
		return err
	}

	return srv.buildResponse(response)
}

//Validate 校验请求是否合法
func (srv *Server) Validate() bool {
	timestamp := srv.Query("timestamp")
	nonce := srv.Query("nonce")
	signature := srv.Query("signature")
	return signature == util.Signature(srv.Token, timestamp, nonce)
}

//HandleRequest 处理微信的请求
func (srv *Server) handleRequest() (reply *message.Reply, err error) {
	//set isSafeMode
	srv.isSafeMode = false
	encryptType := srv.Query("encrypt_type")
	if encryptType == "aes" {
		srv.isSafeMode = true
	}

	//set openID
	srv.openID = srv.Query("openid")

	var msg interface{}
	msg, err = srv.getMessage()
	if err != nil {
		return
	}
	mixMessage, success := msg.(message.MixMessage)
	if !success {
		err = errors.New("消息类型转换失败")
	}
	srv.requestMsg = mixMessage
	reply = srv.messageHandler(mixMessage)
	return
}

//GetOpenID return openID
func (srv *Server) GetOpenID() string {
	return srv.openID
}

//getMessage 解析微信返回的消息
func (srv *Server) getMessage() (interface{}, error) {
	var rawXMLMsgBytes []byte
	var err error
	if srv.isSafeMode {
		var encryptedXMLMsg message.EncryptedXMLMsg
		if err := xml.NewDecoder(srv.Request.Body).Decode(&encryptedXMLMsg); err != nil {
			return nil, fmt.Errorf("从body中解析xml失败,err=%v", err)
		}

		//验证消息签名
		timestamp := srv.Query("timestamp")
		srv.timestamp, err = strconv.ParseInt(timestamp, 10, 32)
		if err != nil {
			return nil, err
		}
		nonce := srv.Query("nonce")
		srv.nonce = nonce
		msgSignature := srv.Query("msg_signature")
		msgSignatureGen := util.Signature(srv.Token, timestamp, nonce, encryptedXMLMsg.EncryptedMsg)
		if msgSignature != msgSignatureGen {
			return nil, fmt.Errorf("消息不合法，验证签名失败")
		}

		//解密
		srv.random, rawXMLMsgBytes, err = util.DecryptMsg(srv.AppID, encryptedXMLMsg.EncryptedMsg, srv.EncodingAESKey)
		if err != nil {
			return nil, fmt.Errorf("消息解密失败, err=%v", err)
		}
	} else {
		rawXMLMsgBytes, err = ioutil.ReadAll(srv.Request.Body)
		if err != nil {
			return nil, fmt.Errorf("从body中解析xml失败, err=%v", err)
		}
	}

	srv.requestRawXMLMsg = rawXMLMsgBytes

	return srv.parseRequestMessage(rawXMLMsgBytes)
}

func (srv *Server) parseRequestMessage(rawXMLMsgBytes []byte) (msg message.MixMessage, err error) {
	msg = message.MixMessage{}
	err = xml.Unmarshal(rawXMLMsgBytes, &msg)
	return
}

//SetMessageHandler 设置用户自定义的回调方法
func (srv *Server) SetMessageHandler(handler func(message.MixMessage) *message.Reply) {
	srv.messageHandler = handler
}

func (srv *Server) buildResponse(reply *message.Reply) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic error: %v\n%s", e, debug.Stack())
		}
	}()
	srv.responseNeedForward = true
	if reply == nil {
		//do nothing
		return nil
	}
	msgType := reply.MsgType
	switch msgType {
	case message.MsgTypeText:
		text := reply.MsgData.(*message.Text)
		if text.Content == "" {
			srv.responseNeedForward = false
			return nil
		}
	case message.MsgTypeImage:
	case message.MsgTypeVoice:
	case message.MsgTypeVideo:
	case message.MsgTypeMusic:
	case message.MsgTypeNews:
	case message.MsgTypeTransfer:
	default:
		err = message.ErrUnsupportReply
		return
	}

	msgData := reply.MsgData
	value := reflect.ValueOf(msgData)
	//msgData must be a ptr
	kind := value.Kind().String()
	if "ptr" != kind {
		return message.ErrUnsupportReply
	}

	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(srv.requestMsg.FromUserName)
	value.MethodByName("SetToUserName").Call(params)

	params[0] = reflect.ValueOf(srv.requestMsg.ToUserName)
	value.MethodByName("SetFromUserName").Call(params)

	params[0] = reflect.ValueOf(msgType)
	value.MethodByName("SetMsgType").Call(params)

	params[0] = reflect.ValueOf(util.GetCurrTs())
	value.MethodByName("SetCreateTime").Call(params)

	srv.responseMsg = msgData
	srv.responseRawXMLMsg, err = xml.Marshal(msgData)
	return
}

//Send 将自定义的消息发送
func (srv *Server) Send() (err error) {
	replyMsg, err := srv.sendBuildMsg(srv.responseMsg)
	if err != nil {
		return
	}
	if replyMsg != nil {
		srv.XML(replyMsg)
	} else {
		if srv.responseNeedForward {
			if srv.mssageForwardUrl != "" {
				srv.MessageForward()
			}
		}
	}
	return
}
func (srv *Server) SetMessageForward(url string, token string) {
	srv.mssageForwardUrl = url
	srv.messageForwardToken = token
}
func (srv *Server) sendBuildMsg(replyMsg interface{}) (interface{}, error) {
	if srv.isSafeMode {
		//安全模式下对消息进行加密
		var encryptedMsg []byte
		var err error
		encryptedMsg, err = util.EncryptMsg(srv.random, srv.responseRawXMLMsg, srv.AppID, srv.EncodingAESKey)
		if err != nil {
			return nil, err
		}
		//TODO 如果获取不到timestamp nonce 则自己生成
		timestamp := srv.timestamp
		timestampStr := strconv.FormatInt(timestamp, 10)
		msgSignature := util.Signature(srv.Token, timestampStr, srv.nonce, string(encryptedMsg))
		replyMsg = message.ResponseEncryptedXMLMsg{
			EncryptedMsg: string(encryptedMsg),
			MsgSignature: msgSignature,
			Timestamp:    timestamp,
			Nonce:        srv.nonce,
		}
	}
	return replyMsg, nil
}
func (srv *Server) MessageForward() {
	signature := util.Signature(srv.messageForwardToken, fmt.Sprint(srv.timestamp), srv.nonce)
	postUrl := srv.mssageForwardUrl + fmt.Sprintf("&timestamp=%d&nonce=%s&signature=%s", srv.timestamp, srv.nonce, signature)
	retryNum := 0
	srv.MessageForwardSend(postUrl, &retryNum)
}

func (srv *Server) MessageForwardSend(postUrl string, retryNum *int) {

	timeout := 4500 * time.Millisecond
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := util.PostXML(postUrl, srv.requestMsg, client)
	if err != nil {
		if strings.Contains(err.Error(), "request canceled (Client.Timeout exceeded while awaiting headers)") {
			msg := &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("系统异常,请稍后再试")}
			err = srv.buildResponse(msg)
			if err != nil {
				fmt.Println("error:", err.Error())
				return
			}
			srv.XML(srv.responseMsg)
			return
		} else if strings.Contains(err.Error(), "request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)") {
			_retryNum := 0
			if retryNum != nil {
				_retryNum = *retryNum
			}
			if _retryNum > 1 {
				return
			}
			_retryNum++
			retryNum = &_retryNum
			srv.MessageForwardSend(postUrl, retryNum)
			return
		}
		fmt.Println("http error,err:", err.Error())
		return
	}
	srv.Render(resp)
}
