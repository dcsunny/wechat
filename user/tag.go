package user

import (
	"fmt"

	"github.com/dcsunny/wechat/common_error"

	"encoding/json"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/define"
	"github.com/dcsunny/wechat/util"
)

type Tag struct {
	*context.Context
}

const (
	createTagURL     = "https://api.weixin.qq.com/cgi-bin/tags/create?access_token=%s"
	getTagURL        = "https://api.weixin.qq.com/cgi-bin/tags/get?access_token=%s"
	updateTagURL     = "https://api.weixin.qq.com/cgi-bin/tags/update?access_token=%s"
	updateUserTagURL = "https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token=%s"
	cancelUserTagURL = "https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token=%s"
	getUserTagURL    = "https://api.weixin.qq.com/cgi-bin/tags/getidlist?access_token=%s"
)

func NewTag(context *context.Context) *Tag {
	tag := new(Tag)
	tag.Context = context
	return tag
}

type CreateTagResp struct {
	Tag TagInfo `json:"tag"`
	define.CommonError
}

type TagInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (tag *Tag) CreateTag(name string) (result CreateTagResp, err error) {
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
	err = json.Unmarshal(response, &result)
	if err != nil {
		fmt.Println(fmt.Sprintf("tag info:%s", string(response)))
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("CreateTag Error , errcode=%d , errmsg=%s", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

type GetTagsResp struct {
	Tags []TagInfo `json:"tags"`
	define.CommonError
}

func (tag *Tag) GetTags() (tags []TagInfo, err error) {
	var accessToken string
	accessToken, err = tag.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(getTagURL, accessToken)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}
	var result GetTagsResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		fmt.Println(fmt.Sprintf("tag info:%s", string(response)))
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("CreateTag Error , errcode=%d , errmsg=%s", result.ErrCode, result.ErrMsg)
		return
	}
	tags = result.Tags
	return
}

type UpdateTagsReq struct {
	Tag TagInfo `json:"tag"`
}

func (tag *Tag) UpdateTags(req TagInfo) (err error) {
	var accessToken string
	accessToken, err = tag.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(updateTagURL, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, UpdateTagsReq{Tag: req})
	if err != nil {
		err = fmt.Errorf("UpdateTags Error , err=%s", err)
		return
	}
	return common_error.DecodeWithCommonError(tag.Context, response, "UpdateTags")
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
	return common_error.DecodeWithCommonError(tag.Context, response, "UpdateUserTag")
}

func (tag *Tag) CancelUserTag(openIDs []string, tagID int) (err error) {
	var accessToken string
	accessToken, err = tag.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(cancelUserTagURL, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, map[string]interface{}{"openid_list": openIDs, "tagid": tagID})
	if err != nil {
		err = fmt.Errorf("CancelUserTag Error , err=%s", err)
		return
	}
	return common_error.DecodeWithCommonError(tag.Context, response, "CancelUserTag")
}

type GetUserTagResp struct {
	TagIDList []int `json:"tagid_list"`
	define.CommonError
}

func (tag *Tag) GetUserTag(openID string) (tags []int, err error) {
	var accessToken string
	accessToken, err = tag.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf(getUserTagURL, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, map[string]interface{}{"openid": openID})
	if err != nil {
		err = fmt.Errorf("GetUserTag Error , err=%s", err)
		return
	}
	var result GetUserTagResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserTag Error , errcode=%d , errmsg=%s", result.ErrCode, result.ErrMsg)
		return
	}
	tags = result.TagIDList
	return
}
