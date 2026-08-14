package system

// 用于构建特定角色的菜单树，包含权限信息
// 用于返回指定角色的菜单数
type SysAuthorizedMenu struct {
	SysMenu
	MenuID      uint                   `json:"menuID" gorm:"comment:菜单ID"`
	AuthorityID uint                   `json:"-" gorm:"comment:角色ID"`
	Children    []SysAuthorizedMenu    `json:"children" gorm:"-"`
	Parameters  []SysBaseMenuParameter `json:"parameters" gorm:"foreignKey:SysBaseMenuID;references:MenuID"`
	Btns        map[string]uint        `json:"btns" gorm:"-"`
}

// SysMenu 和 SysRoleAuthority 之间的多对多关联表
// 作为角色与菜单的关联表模型，实现多对多关系
type SysMenuAuthorityRelation struct {
	MenuID      string `json:"menuID" gorm:"comment:菜单ID;column:menu_id"`
	AuthorityID string `json:"-" gorm:"comment:角色ID;column:authority_id"`
}

func (SysMenuAuthorityRelation) TableName() string {
	return "sys_menu_authority_relation"
}
