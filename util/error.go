package util

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/define"
)

// CommonError 微信返回的通用错误json
type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func DecodeWithCommonError(context *context.Context, response []byte, apiName string) (err error) {
	var commError CommonError
	err = json.Unmarshal(response, &commError)
	if err != nil {
		return
	}
	if commError.ErrCode != 0 {
		CommonErrorHandle(commError, context)
		return fmt.Errorf("%s Error , errcode=%d , errmsg=%s", apiName, commError.ErrCode, commError.ErrMsg)
	}
	return nil
}

func CommonErrorHandle(commError CommonError, context *context.Context) {
	if commError.ErrCode == 40001 {
		accessTokenCacheKey := fmt.Sprintf(define.AccessTokenCacheKey, context.AppID)
		context.Cache.Delete(accessTokenCacheKey)
	}
}
