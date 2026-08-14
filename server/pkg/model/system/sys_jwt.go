package system

import "server/pkg/global"

type JwtBlacklist struct {
	global.MPA_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
