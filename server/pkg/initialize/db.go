package initialize

import (
	"server/pkg/db"
	"server/pkg/global"
	"server/pkg/model/system"

	adapter "github.com/casbin/gorm-adapter/v3"
)

func registerTables(customizedTables ...any) {
	if global.MPA_DB == nil {
		global.MPA_LOG.Error("db is nil")
		return
	}

	var tables []any

	// 注册自定义表
	tables = append(tables, customizedTables...)

	defaultTables := defaultTable()

	// 注册系统表
	tables = append(tables, defaultTables...)

	db.RegisterTables(tables...)
}

func defaultTable() []any {
	return []any{
		system.SysUser{},
		system.SysRoleAuthority{},
		system.SysUserAuthorityRelation{},
		system.SysMenu{},
		system.SysMenuAuthorityRelation{},
		system.SysApi{},
		system.SysIgnoreApi{},
		system.InitFlag{},
		adapter.CasbinRule{},
	}
}
