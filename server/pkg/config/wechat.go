package config

type WeChat struct {
	AppID     string `mapstructure:"appid" json:"appid" yaml:"appid"`
	AppSecret string `mapstructure:"appsecret" json:"appsecret" yaml:"appsecret"`
	Token     string `mapstructure:"token" json:"token" yaml:"token"`
}
