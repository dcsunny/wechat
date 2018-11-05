package user

import (
	"fmt"

	"encoding/json"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/util"
)

type Tag struct {
	*context.Context
}

const (
	createTagURL     = "https://api.weixin.qq.com/cgi-bin/tags/create?access_token=%s"
	updateUserTagURL = "https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token=%s"
)

func NewTag(context *context.Context) *Tag {
	tag := new(Tag)
	tag.Context = context
	return tag
}

type TagInfo struct {
	Tag struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"tag"`
	util.CommonError
}

func (tag *Tag) CreateTag(name string) (tagInfo TagInfo, err error) {
	var accessToken string
	accessToken, err = tag.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(createTagURL, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, map[string]map[string]string{"tag": map[string]string{"name": name}})
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &tagInfo)
	if err != nil {
		fmt.Println(fmt.Sprintf("tag info:%s", string(response)))
		return
	}
	if tagInfo.ErrCode != 0 {
		err = fmt.Errorf("CreateTag Error , errcode=%d , errmsg=%s", tagInfo.ErrCode, tagInfo.ErrMsg)
		return
	}
	return
}

func (tag *Tag) UpdateUserTag(openIDs []string, tagID int) (err error) {
	var accessToken string
	accessToken, err = tag.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(updateUserTagURL, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, map[string]interface{}{"openid_list": openIDs, "tagid": tagID})
	if err != nil {
		err = fmt.Errorf("UpdateUserTag Error , err=%s", err)
		return
	}
	fmt.Println(string(response))
	return util.DecodeWithCommonError(tag.Context, response, "UpdateUserTag")
}
