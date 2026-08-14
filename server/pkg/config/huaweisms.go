package config

type HuaweiSMS struct {
	AK             string `mapstructure:"ak" json:"ak" yaml:"ak"`                                     //App Key
	SK             string `mapstructure:"sk" json:"sk" yaml:"sk"`                                     //App Secret
	ApiAddress     string `mapstructure:"apiaddress" json:"apiaddress" yaml:"apiaddress"`             // APP接入地址(在控制台"应用管理"页面获取)+接口访问URI
	Sender         string `mapstructure:"sender" json:"sender" yaml:"sender"`                         // 国内短信签名通道号
	Signature      string `mapstructure:"signature" json:"signature" yaml:"signature"`                // 签名名称
	TemplateID     string `mapstructure:"templateId" json:"templateId" yaml:"templateId"`             // 模板ID
	StatusCallBack string `mapstructure:"statusCallBack" json:"statusCallBack" yaml:"statusCallBack"` // 选填,短信状态报告接收地址,推荐使用域名,为空或者不填表示不接收状态报告
}
