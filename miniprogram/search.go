package miniprogram

import (
	"encoding/json"
	"fmt"

	"github.com/dcsunny/wechat/define"

	"github.com/dcsunny/wechat/common_error"
	"github.com/dcsunny/wechat/util"
)

const (
	searchSubmitPages = "https://api.weixin.qq.com/wxa/search/wxaapi_submitpages"
	searchSiteSearch  = "https://api.weixin.qq.com/wxa/sitesearch"
)

type SearchSubmitPagesReq struct {
	Pages []SearchSubmitPagesInfo `json:"pages"`
}

type SearchSubmitPagesInfo struct {
	Path  string `json:"path"`
	Query string `json:"query"`
}

func (wxa *MiniProgram) SearchSubmitPages(req SearchSubmitPagesReq) error {
	accessToken, err := wxa.GetAccessToken()
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s?access_token=%s", searchSubmitPages, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return err
	}
	return common_error.DecodeWithCommonError(wxa.Context, response, "searchSubmitPages")
}

type SearchSiteSearchReq struct {
	Keyword      string `json:"keyword"`
	NextPageInfo string `json:"next_page_info"`
}

type SearchSiteSearchResp struct {
	define.CommonError
	HasNextPage  int                    `json:"has_next_page"`
	HitCount     int                    `json:"hit_count"`
	Items        []SearchSiteSearchItem `json:"items"`
	NextPageInfo string                 `json:"next_page_info"`
}

type SearchSiteSearchItem struct {
	Description string `json:"description"`
	Image       string `json:"image"`
	Path        string `json:"path"`
	Title       string `json:"title"`
}

func (wxa *MiniProgram) SearchSiteSearch(req SearchSiteSearchReq) (SearchSiteSearchResp, error) {
	accessToken, err := wxa.GetAccessToken()
	if err != nil {
		return SearchSiteSearchResp{}, err
	}

	uri := fmt.Sprintf("%s?access_token=%s", searchSiteSearch, accessToken)
	response, err := util.PostJSON(uri, req)
	if err != nil {
		return SearchSiteSearchResp{}, err
	}
	var result SearchSiteSearchResp
	err = json.Unmarshal(response, &result)
	if err != nil {
		return SearchSiteSearchResp{}, err
	}
	if result.ErrCode != 0 {
		return SearchSiteSearchResp{}, common_error.CommonErrorHandle(result.CommonError, wxa.Context, "searchSiteSearch")
	}
	return result, nil
}
