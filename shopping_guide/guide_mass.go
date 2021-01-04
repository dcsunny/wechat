package shopping_guide

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/common_error"
	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/define"
	"github.com/dcsunny/wechat/util"
)

const (
	addGuideMasssendJob    = "https://api.weixin.qq.com/cgi-bin/guide/addguidemassendjob"
	getGuideMassendJobList = "https://api.weixin.qq.com/cgi-bin/guide/getguidemassendjoblist"
	getGuideMassendJob     = "https://api.weixin.qq.com/cgi-bin/guide/getguidemassendjob"
	updateGuideMasssendJob = "https://api.weixin.qq.com/cgi-bin/guide/updateguidemassendjob"
	cancelGuideMassendJob  = "https://api.weixin.qq.com/cgi-bin/guide/cancelguidemassendjob"
)

const (
	GuideMassTaskStatusNotExecuted = iota + 1
	GuideMassTaskStatusHadExecuted
	GuideMassTaskStatusHadComplete
	GuideMassTaskStatusCancel
)

const (
	GuideMassTaskSendStatusNotSend = iota + 1
	GuideMassTaskSendStatusSendSuccess
	GuideMassTaskSendStatusSendNotSubscribe
	GuideMassTaskSendStatusSendNotQuota
	GuideMassTaskSendStatusSendSystemError
)

type GuideMass struct {
	*context.Context
}

func NewGuideMass(ctx *context.Context) *GuideMass {
	return &GuideMass{ctx}
}

type AddGuideMasssendJobReq struct {
	GuideAccount string                `json:"guide_account"`
	GuideOpenID  string                `json:"guide_openid"`
	Material     []MassSendJobMaterial `json:"material"`
	OpenID       []string              `json:"openid"`
	PushTime     int64                 `json:"push_time"`
	TaskName     string                `json:"task_name"`
	TaskRemark   string                `json:"task_remark"`
}

type MassSendJobMaterial struct {
	AppID   string `json:"appid"`
	MediaID string `json:"media_id"`
	Path    string `json:"path"`
	Title   string `json:"title"`
	Type    int64  `json:"type"`
	Word    string `json:"word"`
}

type AddGuideMasssendJobResp struct {
	TaskResult []AddGuideMasssendJobTaskInfo `json:"task_result"`
	define.CommonError
}

type AddGuideMasssendJobTaskInfo struct {
	OpenID []string `json:"openid"`
	TaskID int64    `json:"task_id"`
}

func (g *GuideMass) AddGuideMasssendJob(req AddGuideMasssendJobReq) ([]AddGuideMasssendJobTaskInfo, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", addGuideMasssendJob, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return nil, err
	}
	var result AddGuideMasssendJobResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "addGuideMasssendJob")
	}
	return result.TaskResult, nil
}

type GetGuideMassendJobListReq struct {
	GuideAccount string `json:"guide_account"`
	GuideOpenID  string `json:"guide_openid"`
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
	TaskStatus   []int  `json:"task_status"`
}

type GetGuideMassendJobListResp struct {
	define.CommonError
	List       []GuideMassJob `json:"list"`
	TotalCount int64          `json:"total_count"`
}

type GuideMassJob struct {
	BuyerInfo  []GuideMassBuyerInfo  `json:"buyer_info"`
	CreateTime int64                 `json:"create_time"`
	FinishTime int64                 `json:"finish_time"`
	Material   []MassSendJobMaterial `json:"material"`
	PushTime   int64                 `json:"push_time"`
	TaskID     int64                 `json:"task_id"`
	TaskName   string                `json:"task_name"`
	TaskRemark string                `json:"task_remark"`
	TaskStatus int                   `json:"task_status"`
	UpdateTime int64                 `json:"update_time"`
}

type GuideMassBuyerInfo struct {
	OpenID     string `json:"openid"`
	SendStatus int    `json:"send_status"`
}

func (g *GuideMass) GetGuideMassendJobList(req GetGuideMassendJobListReq) (*GetGuideMassendJobListResp, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", getGuideMassendJobList, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return nil, err
	}
	var result GetGuideMassendJobListResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, common_error.CommonErrorHandle(result.CommonError, g.Context, "getGuideMassendJobList")
	}
	return &result, nil
}

type GetGuideMassendJobResp struct {
	Job GuideMassJob `json:"job"`
	define.CommonError
}

func (g *GuideMass) GetGuideMassendJob(taskID int64) (GuideMassJob, error) {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return GuideMassJob{}, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", getGuideMassendJob, accessToken)
	response, err := util.PostJSON(uri, map[string]interface{}{
		"task_id": taskID,
	})
	if err != nil {
		return GuideMassJob{}, err
	}
	var result GetGuideMassendJobResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return GuideMassJob{}, err
	}
	if result.ErrCode != 0 {
		return GuideMassJob{}, common_error.CommonErrorHandle(result.CommonError, g.Context, "getGuideMassendJob")
	}
	return result.Job, nil
}

type UpdateGuideMasssendJobReq struct {
	Material   []MassSendJobMaterial `json:"material"`
	OpenID     []string              `json:"openid"`
	PushTime   int64                 `json:"push_time"`
	TaskID     int64                 `json:"task_id"`
	TaskName   string                `json:"task_name"`
	TaskRemark string                `json:"task_remark"`
}

//无法修改已经执行的任务，返回参数错误。
func (g *GuideMass) UpdateGuideMasssendJob(req UpdateGuideMasssendJobReq) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", updateGuideMasssendJob, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "updateGuideMasssendJob")
}

func (g *GuideMass) CancelGuideMassendJob(taskID int64) error {
	accessToken, err := g.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", cancelGuideMassendJob, accessToken)
	response, err := util.PostJSON(uri, map[string]interface{}{
		"task_id": taskID,
	})
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(g.Context, response, "cancelGuideMassendJob")
}
