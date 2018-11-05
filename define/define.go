package define

const (
	AccessTokenCacheKey = "access_token:%s"
)

// CommonError 微信返回的通用错误json
type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
