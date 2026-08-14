package api

// 微信服务器验证request model
type ServerValidateReqModel struct {
	Signature string `json:"signature" form:"signature"`
	TimeStamp string `json:"timestamp" form:"timestamp"`
	Nonce     string `json:"nonce" form:"nonce"`
	Echostr   string `json:"echostr" form:"echostr"`
}
