package sms

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"server/pkg/global"
)

type HuaweiSMS struct{}

func (*HuaweiSMS) Sender(receiver string, templateParas string) (SmsResponse, error) {
	app := NewSMSApp()

	templateParas = fmt.Sprintf("[\"%s\"]", templateParas)
	body := buildRequestBody(receiver, templateParas)
	return post([]byte(body), app)
}

func NewSMSApp() Signer {
	return Signer{Key: global.MPA_CONFIG.HuaweiSMS.AK, Secret: global.MPA_CONFIG.HuaweiSMS.SK}
}

/**
 * sender,receiver,templateId不能为空
 */
func buildRequestBody(receiver, templateParas string) string {
	param := "from=" + url.QueryEscape(global.MPA_CONFIG.HuaweiSMS.Sender) + "&to=" + url.QueryEscape(receiver) + "&templateId=" + url.QueryEscape(global.MPA_CONFIG.HuaweiSMS.TemplateID)
	if templateParas != "" {
		param += "&templateParas=" + url.QueryEscape(templateParas)
	}
	if global.MPA_CONFIG.HuaweiSMS.StatusCallBack != "" {
		param += "&statusCallback=" + url.QueryEscape(global.MPA_CONFIG.HuaweiSMS.StatusCallBack)
	}
	if global.MPA_CONFIG.HuaweiSMS.Signature != "" {
		param += "&signature=" + url.QueryEscape(global.MPA_CONFIG.HuaweiSMS.Signature)
	}
	return param
}

type SmsResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func post(param []byte, appInfo Signer) (SmsResponse, error) {
	var smsResp SmsResponse

	if param == nil || (appInfo == Signer{}) {
		return smsResp, nil
	}

	// 代码样例为了简便，设置了不进行证书校验，请在商用环境自行开启证书校验。
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", global.MPA_CONFIG.HuaweiSMS.ApiAddress, bytes.NewBuffer(param))
	if err != nil {
		return smsResp, err
	}

	// 对请求增加内容格式，固定头域
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// 对请求进行HMAC算法签名，并将签名结果设置到Authorization头域。
	appInfo.Sign(req)

	// 发送短信请求
	resp, err := client.Do(req)
	if err != nil {
		return smsResp, err
	}

	// 获取短信响应
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return smsResp, err
	}

	json.Unmarshal(body, &smsResp)
	return smsResp, nil
}
