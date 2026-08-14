package config

type System struct {
	Name               string   `mapstructure:"name" json:"name" yaml:"name"`                                                 // 系统名称
	Env                string   `mapstructure:"env" json:"env" yaml:"env"`                                                    // 环境值
	Addr               int      `mapstructure:"addr" json:"addr" yaml:"addr"`                                                 // 端口值
	DbType             string   `mapstructure:"db-type" json:"db-type" yaml:"db-type"`                                        // 数据库类型:mysql(默认)|sqlite|postgresql
	OssType            string   `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"`                                     // Oss类型
	UseRedis           bool     `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                                  // 使用redis
	UseMultipoint      bool     `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"`                   // 多点登录限制, 禁止用户多点在线
	RouterPrefix       string   `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`                      // 路由前缀
	SystemRouterPrefix string   `mapstructure:"system-router-prefix" json:"system-router-prefix" yaml:"system-router-prefix"` // 系统路由前缀
	SystemRouters      []string `mapstructure:"system-router" json:"system-router" yaml:"system-router"`                      // 默认系统路由
	Middleware         []string `mapstructure:"middleware" json:"middleware" yaml:"middleware"`                               // 系统中间件
	Plugins            []string `mapstructure:"plugins" json:"plugins" yaml:"plugins"`                                        // 系统插件
	UseStrictAuth      bool     `mapstructure:"use-strict-auth" json:"use-strict-auth" yaml:"use-strict-auth"`                // 使用树形角色分配模式
	DisableAutoMigrate bool     `mapstructure:"disable-auto-migrate" json:"disable-auto-migrate" yaml:"disable-auto-migrate"` // 自动迁移数据库表结构，生产环境建议设为false，手动迁移
	// SystemTables       []string `mapstructure:"system-tables" json:"system-tables" yaml:"system-tables"`                      // 系统表
}
