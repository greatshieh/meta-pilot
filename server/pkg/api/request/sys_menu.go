package request

import (
	"server/pkg/global"
	"server/pkg/model/system"
)

func DefaultMenu() []system.SysMenu {
	return []system.SysMenu{{
		MPA_MODEL: global.MPA_MODEL{ID: 1},
		ParentID:  0,
		Path:      "dashboard",
		Name:      "dashboard",
		Component: "views/dashboard/index.vue",
		Sort:      1,
		Meta: system.Meta{
			Title: "首页",
			Icon:  "dashboard",
			Affix: true,
		},
	}}
}
