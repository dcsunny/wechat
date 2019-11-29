package error

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/context"

	"github.com/dcsunny/wechat/define"
)

func DecodeWithCommonError(context *context.Context, response []byte, apiName string) (err error) {
	var commError define.CommonError
	err = json.Unmarshal(response, &commError)
	if err != nil {
		return
	}
	if commError.ErrCode != 0 {
		return CommonErrorHandle(commError, context, apiName)
	}
	return nil
}

func CommonErrorHandle(commError define.CommonError, context *context.Context, apiName string) error {
	if commError.ErrCode == 40001 {
		accessTokenCacheKey := fmt.Sprintf(define.AccessTokenCacheKey, context.AppID)
		context.Cache.Delete(accessTokenCacheKey)
	}
	return fmt.Errorf("%s Error , errcode=%d , errmsg=%s", apiName, commError.ErrCode, commError.ErrMsg)
}
