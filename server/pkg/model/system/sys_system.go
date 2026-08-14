package system

import (
	"server/pkg/config"
	"server/pkg/global"
)

// 配置文件结构体
type System struct {
	Config config.Server `json:"config"`
}

// 系统初始化标记表（核心：控制只执行一次）
type InitFlag struct {
	global.MPA_MODEL
	IsInitialized bool `gorm:"default:false"` // 是否已初始化
}
