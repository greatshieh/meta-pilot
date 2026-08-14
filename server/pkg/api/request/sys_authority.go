package request

import "server/pkg/model/system"

type SetMenuAuthority struct {
	Menus []system.SysMenu `json:"menus" form:"menus"`
}
