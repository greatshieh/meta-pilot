package api

type UserWechatToken struct {
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	OpenID         string `json:"openid"`
	Scope          string `json:"scope"`
	IsSnapShotUser int    `json:"is_snapshotuser"`
	Unionid        string `json:"unionid"`
	ErrCode        int    `json:"errcode"`
	Errmsg         string `json:"errmsg"`
}

type UserWechatInfo struct {
	OpenID     string `json:"openid"`
	Headimgurl string `json:"Headimgurl"`
}
