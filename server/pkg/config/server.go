package config

type Server struct {
	JWT         JWT           `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap         Zap           `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis       Redis         `mapstructure:"redis" json:"redis" yaml:"redis"`
	Email       Email         `mapstructure:"email" json:"email" yaml:"email"`
	System      System        `mapstructure:"system" json:"system" yaml:"system"`
	RateLimiter []RateLimiter `mapstructure:"ratelimiter" json:"ratelimiter" yaml:"ratelimiter"`

	// gorm
	Mysql  Mysql           `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Pgsql  Pgsql           `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	Sqlite Sqlite          `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`
	DBList []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`

	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`

	// wechat配置
	WeChat WeChat `mapstructure:"wechat" json:"wechat" yaml:"wechat"`

	// 华为短息服务
	HuaweiSMS HuaweiSMS `mapstructure:"huaweiSMS" json:"huaweiSMS" yaml:"huaweiSMS"`

	// 七牛Kodo
	QiniuKodo Qiniu `mapstructure:"qiniukodo" json:"qiniukodo" yaml:"qiniukodo"`

	// 华为OBS
	HuaweiOBS Huawei `mapstructure:"huaweiobs" json:"huaweiobs" yaml:"huaweiobs"`
}
