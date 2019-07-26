package menu

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/context"
	"github.com/dcsunny/wechat/define"
	error2 "github.com/dcsunny/wechat/error"
	"github.com/dcsunny/wechat/util"
)

const (
	menuCreateURL            = "/cgi-bin/menu/create"
	menuGetURL               = "/cgi-bin/menu/get"
	menuDeleteURL            = "/cgi-bin/menu/delete"
	menuAddConditionalURL    = "/cgi-bin/menu/addconditional"
	menuDeleteConditionalURL = "/cgi-bin/menu/delconditional"
	menuTryMatchURL          = "/cgi-bin/menu/trymatch"
	menuSelfMenuInfoURL      = "/cgi-bin/get_current_selfmenu_info"
)

//Menu struct
type Menu struct {
	*context.Context
}

//reqMenu 设置菜单请求数据
type reqMenu struct {
	Button    []*Button  `json:"button,omitempty"`
	MatchRule *MatchRule `json:"matchrule,omitempty"`
}

//reqDeleteConditional 删除个性化菜单请求数据
type reqDeleteConditional struct {
	MenuID int64 `json:"menuid"`
}

//reqMenuTryMatch 菜单匹配请求
type reqMenuTryMatch struct {
	UserID string `json:"user_id"`
}

//resConditionalMenu 个性化菜单返回结果
type resConditionalMenu struct {
	Button    []Button  `json:"button"`
	MatchRule MatchRule `json:"matchrule"`
	MenuID    int64     `json:"menuid"`
}

//resMenuTryMatch 菜单匹配请求结果
type resMenuTryMatch struct {
	define.CommonError

	Button []Button `json:"button"`
}

//ResMenu 查询菜单的返回数据
type ResMenu struct {
	define.CommonError

	Menu struct {
		Button []Button `json:"button"`
		MenuID int64    `json:"menuid"`
	} `json:"menu"`
	Conditionalmenu []resConditionalMenu `json:"conditionalmenu"`
}

//ResSelfMenuInfo 自定义菜单配置返回结果
type ResSelfMenuInfo struct {
	define.CommonError

	IsMenuOpen   int32 `json:"is_menu_open"`
	SelfMenuInfo struct {
		Button []SelfMenuButton `json:"button"`
	} `json:"selfmenu_info"`
}

//SelfMenuButton 自定义菜单配置详情
type SelfMenuButton struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	URL       string `json:"url,omitempty"`
	Value     string `json:"value,omitempty"`
	SubButton struct {
		List []SelfMenuButton `json:"list"`
	} `json:"sub_button,omitempty"`
	NewsInfo struct {
		List []ButtonNew `json:"list"`
	} `json:"news_info,omitempty"`
}

//ButtonNew 图文消息菜单
type ButtonNew struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	Digest     string `json:"digest"`
	ShowCover  int32  `json:"show_cover"`
	CoverURL   string `json:"cover_url"`
	ContentURL string `json:"content_url"`
	SourceURL  string `json:"source_url"`
}

//MatchRule 个性化菜单规则
type MatchRule struct {
	GroupID            int32  `json:"group_id,omitempty"`
	Sex                int32  `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType int32  `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
}

//NewMenu 实例
func NewMenu(context *context.Context) *Menu {
	menu := new(Menu)
	menu.Context = context
	return menu
}

//SetMenu 设置按钮
func (menu *Menu) SetMenu(buttons []*Button) error {
	accessToken, err := menu.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", menu.ApiBaseUrl+menuCreateURL, accessToken)
	reqMenu := &reqMenu{
		Button: buttons,
	}

	response, err := util.PostJSON(uri, reqMenu)
	if err != nil {
		return err
	}
	return error2.DecodeWithCommonError(menu.Context, response, "SetMenu")
}

//GetMenu 获取菜单配置
func (menu *Menu) GetMenu() (resMenu ResMenu, err error) {
	var accessToken string
	accessToken, err = menu.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", menu.ApiBaseUrl+menuGetURL, accessToken)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &resMenu)
	if err != nil {
		return
	}
	if resMenu.ErrCode != 0 {
		err = fmt.Errorf("GetMenu Error , errcode=%d , errmsg=%s", resMenu.ErrCode, resMenu.ErrMsg)
		return
	}
	return
}

//DeleteMenu 删除菜单
func (menu *Menu) DeleteMenu() error {
	accessToken, err := menu.GetAccessToken()
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("%s?access_token=%s", menu.ApiBaseUrl+menuDeleteURL, accessToken)
	response, err := util.HTTPGet(uri)
	if err != nil {
		return err
	}
	return error2.DecodeWithCommonError(menu.Context, response, "DeleteMenu")
}

//AddConditional 添加个性化菜单
func (menu *Menu) AddConditional(buttons []*Button, matchRule *MatchRule) error {
	accessToken, err := menu.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", menu.ApiBaseUrl+menuAddConditionalURL, accessToken)
	reqMenu := &reqMenu{
		Button:    buttons,
		MatchRule: matchRule,
	}

	response, err := util.PostJSON(uri, reqMenu)
	if err != nil {
		return err
	}
	return error2.DecodeWithCommonError(menu.Context, response, "AddConditional")
}

//DeleteConditional 删除个性化菜单
func (menu *Menu) DeleteConditional(menuID int64) error {
	accessToken, err := menu.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", menu.ApiBaseUrl+menuDeleteConditionalURL, accessToken)
	reqDeleteConditional := &reqDeleteConditional{
		MenuID: menuID,
	}

	response, err := util.PostJSON(uri, reqDeleteConditional)
	if err != nil {
		return err
	}
	return error2.DecodeWithCommonError(menu.Context, response, "DeleteConditional")
}

//MenuTryMatch 菜单匹配
func (menu *Menu) MenuTryMatch(userID string) (buttons []Button, err error) {
	var accessToken string
	accessToken, err = menu.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", menu.ApiBaseUrl+menuTryMatchURL, accessToken)
	reqMenuTryMatch := &reqMenuTryMatch{userID}
	var response []byte
	response, err = util.PostJSON(uri, reqMenuTryMatch)
	if err != nil {
		return
	}
	var resMenuTryMatch resMenuTryMatch
	err = json.Unmarshal(response, &resMenuTryMatch)
	if err != nil {
		return
	}
	if resMenuTryMatch.ErrCode != 0 {
		err = fmt.Errorf("MenuTryMatch Error , errcode=%d , errmsg=%s", resMenuTryMatch.ErrCode, resMenuTryMatch.ErrMsg)
		return
	}
	buttons = resMenuTryMatch.Button
	return
}

//GetCurrentSelfMenuInfo 获取自定义菜单配置接口
func (menu *Menu) GetCurrentSelfMenuInfo() (resSelfMenuInfo ResSelfMenuInfo, err error) {
	var accessToken string
	accessToken, err = menu.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", menu.ApiBaseUrl+menuSelfMenuInfoURL, accessToken)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}
	fmt.Println(string(response))
	err = json.Unmarshal(response, &resSelfMenuInfo)
	if err != nil {
		return
	}
	if resSelfMenuInfo.ErrCode != 0 {
		err = fmt.Errorf("GetCurrentSelfMenuInfo Error , errcode=%d , errmsg=%s", resSelfMenuInfo.ErrCode, resSelfMenuInfo.ErrMsg)
		return
	}
	return
}
