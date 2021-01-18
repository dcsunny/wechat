package miniprogram

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/define"
	"github.com/dcsunny/wechat/util"
)

const (
	imgSecCheckUrl     = "https://api.weixin.qq.com/wxa/img_sec_check"
	msgSecCheckUrl     = "https://api.weixin.qq.com/wxa/msg_sec_check"
	mediaCheckAsyncUrl = "https://api.weixin.qq.com/wxa/media_check_async"
)

func (s *MiniProgram) ImgSecCheck(filename string, fileBytes []byte) (err error) {
	var accessToken string
	accessToken, err = s.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s", imgSecCheckUrl, accessToken)

	fields := []util.MultipartFormField{
		{
			IsFile:    true,
			Fieldname: "media",
			Filename:  filename,
			Value:     fileBytes,
		},
	}

	var response []byte
	response, err = util.PostMultipartForm(fields, uri)
	if err != nil {
		return
	}
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("ImgSecCheck error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

func (s *MiniProgram) MsgSecCheck(content string) (err error) {
	var accessToken string
	accessToken, err = s.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", msgSecCheckUrl, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, map[string]interface{}{
		"content": content,
	})
	if err != nil {
		return
	}
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("MsgSecCheck error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return err
}

func (s *MiniProgram) MediaCheckAsync(mediaUrl string, mediaType int) (err error) {
	var accessToken string
	accessToken, err = s.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", mediaCheckAsyncUrl, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, map[string]interface{}{
		"media_url":  mediaUrl,
		"media_type": mediaType,
	})
	if err != nil {
		return
	}
	var result define.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("MediaCheckAsync error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return err
}
