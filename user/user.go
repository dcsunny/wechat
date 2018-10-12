package user

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/util"
)

const (
	userInfoURL     = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
	updateRemarkURL = "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token=%s"
	userListURL     = "https://api.weixin.qq.com/cgi-bin/user/get"
)

//User 用户管理
type User struct {
	*context.Context
}

//NewUser 实例化
func NewUser(context *context.Context) *User {
	user := new(User)
	user.Context = context
	return user
}

//Info 用户基本信息
type Info struct {
	util.CommonError

	Subscribe     int32   `json:"subscribe"`
	OpenID        string  `json:"openid"`
	Nickname      string  `json:"nickname"`
	Sex           int32   `json:"sex"`
	City          string  `json:"city"`
	Country       string  `json:"country"`
	Province      string  `json:"province"`
	Language      string  `json:"language"`
	Headimgurl    string  `json:"headimgurl"`
	SubscribeTime int32   `json:"subscribe_time"`
	UnionID       string  `json:"unionid"`
	Remark        string  `json:"remark"`
	GroupID       int32   `json:"groupid"`
	TagidList     []int32 `json:"tagid_list"`
}

//GetUserInfo 获取用户基本信息
func (user *User) GetUserInfo(openID string) (userInfo Info, err error) {
	var accessToken string
	accessToken, err = user.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf(userInfoURL, accessToken, openID)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}
	userInfo = Info{}
	err = json.Unmarshal(response, &userInfo)
	if err != nil {
		fmt.Println(fmt.Sprintf("get user info:%s", string(response)))
		return
	}
	if userInfo.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo Error , errcode=%d , errmsg=%s", userInfo.ErrCode, userInfo.ErrMsg)
		return
	}
	return
}

type ListResult struct {
	util.CommonError

	Total int64 `json:"total"`
	Count int64 `json:"count"`
	Data  struct {
		OpenID []string `json:"openid"`
	} `json:"data"`
	NextOpenID string `json:"next_openid"`
}

func (user *User) List(nexOpenID string) (users ListResult, err error) {
	var accessToken string
	accessToken, err = user.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s&next_openid=%s", userListURL, accessToken, nexOpenID)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}

	users = ListResult{}
	err = json.Unmarshal(response, &users)
	if err != nil {
		fmt.Println(fmt.Sprintf("get user info:%s", string(response)))
		return
	}
	if users.ErrCode != 0 {
		err = fmt.Errorf("get user list Error , errcode=%d , errmsg=%s", users.ErrCode, users.ErrMsg)
		return
	}
	return
}

func (user *User) UpdateRemark(openID, remark string) (err error) {
	var accessToken string
	accessToken, err = user.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(updateRemarkURL, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, map[string]string{"openid": openID, "remark": remark})
	if err != nil {
		return
	}
	return util.DecodeWithCommonError(response, "UpdateRemark")
}
