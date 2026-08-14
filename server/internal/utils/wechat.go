package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/internal/api"
	"server/pkg/global"

	"github.com/marmotedu/errors"
)

/*
从微信服务器获得用户信息

	params:
	------------
	code: string
	  从前端获取的code

	return
	------------
	user_info: dict
	  从微信服务器获取的用户信息，包括：
	  {
	    "openid":" OPENID",
	    "nickname": NICKNAME,
	    "sex":"1",
	    "province":"PROVINCE"
	    "city":"CITY",
	    "country":"COUNTRY",
	    "headimgurl":"http://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/46",
	    "privilege":[ "PRIVILEGE1" "PRIVILEGE2"     ],
	    "unionid": "o6_bmasdasdsad6_2sgVt7hMZOPfL"
	   }
*/

// 通过code换取微信授权access_token
func GetUserAccssToken(code string) (api.UserWechatToken, error) {
	var accessTokenResp api.UserWechatToken
	client := http.DefaultClient

	uri := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", global.MPA_CONFIG.WeChat.AppID, global.MPA_CONFIG.WeChat.AppSecret, code)

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return accessTokenResp, errors.WithCode(50001, "%s", err.Error())
	}

	resp, err := client.Do(request)
	if err != nil {
		return accessTokenResp, errors.WithCode(50002, "%s", err.Error())
	}

	defer resp.Body.Close()

	content, _ := io.ReadAll(resp.Body)

	json.Unmarshal(content, &accessTokenResp)

	return accessTokenResp, nil

}

// 从微信服务器获取用户信息
func GetWechatUserInfo(accessToken string, openID string) (api.UserWechatInfo, error) {
	var userInfo api.UserWechatInfo
	client := http.DefaultClient

	infoUri := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", accessToken, openID)

	request, err := http.NewRequest("GET", infoUri, nil)
	if err != nil {
		return userInfo, errors.WithCode(50001, "%s", err.Error())
	}

	resp, err := client.Do(request)
	if err != nil {
		return userInfo, errors.WithCode(50003, "%s", err.Error())
	}

	defer resp.Body.Close()

	content, _ := io.ReadAll(resp.Body)

	json.Unmarshal(content, &userInfo)

	return userInfo, nil
}

func GetUserInfo(code string) (api.UserWechatInfo, error) {
	var userInfo api.UserWechatInfo

	accessTokenResp, err := GetUserAccssToken(code)

	if err != nil {
		return userInfo, err
	}

	if accessTokenResp.ErrCode != 0 {
		return userInfo, errors.WrapC(errors.New("未知错误"), 50004, "错误的用户授权码")
	}

	return GetWechatUserInfo(accessTokenResp.AccessToken, accessTokenResp.OpenID)
}
