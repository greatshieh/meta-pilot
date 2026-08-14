package system

import "server/pkg/global"

type SysBaseMenuBtn struct {
	global.MPA_MODEL
	Name      string `json:"name" gorm:"comment:按钮关键key"`
	Desc      string `json:"desc" gorm:"按钮备注"`
	SysMenuID uint   `json:"sys_menu_id" gorm:"comment:菜单ID"`
}
