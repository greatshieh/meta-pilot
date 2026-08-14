package system

import "server/pkg/global"

type SysMenu struct {
	global.MPA_MODEL
	MenuLevel     uint                                       `json:"-"`
	ParentID      uint                                       `json:"parentID" gorm:"comment:父菜单ID"`     // 父菜单ID
	Path          string                                     `json:"path" gorm:"comment:路由path"`        // 路由path
	Name          string                                     `json:"name" gorm:"comment:路由name"`        // 路由name
	Component     string                                     `json:"component" gorm:"comment:对应前端文件路径"` // 对应前端文件路径
	Redirect      string                                     `json:"redirect" gorm:"comment:路由跳转地址"`
	Sort          int                                        `json:"sort" gorm:"comment:排序标记"` // 排序标记
	Meta          `json:"meta" gorm:"embedded;comment:附加属性"` // 附加属性
	SysAuthoritys []SysRoleAuthority                         `json:"authoritys" gorm:"many2many:sys_menu_authority_relation;joinForeignKey:menu_id;joinReferences:authority_id;"` // 授权的子角色
	Children      []SysMenu                                  `json:"children" gorm:"-"`
	// Parameters    []SysBaseMenuParameter                     `json:"parameters"`
	// MenuBtn       []SysBaseMenuBtn                           `json:"menuBtn"`
}

type Meta struct {
	ActiveName     string `json:"activeName" gorm:"comment:高亮菜单"`
	Affix          bool   `json:"affix" gorm:"default:false; comment:是否固定在标签页"`
	AlwaysShow     bool   `json:"alwaysShow" gorm:"default:false; comment:是否总是显示"`
	Breadcrumb     bool   `json:"breadcrumb" gorm:"default:true; comment:是否显示面包屑"`
	CloseTab       bool   `json:"closeTab" gorm:"default:false;comment:自动关闭tab"`
	Hidden         bool   `json:"hidden" gorm:"default:false; comment:是否隐藏"`
	Href           string `json:"href" gorm:"comment:菜单外链"`
	Icon           string `json:"icon" gorm:"comment:菜单图标"`
	KeepAlive      bool   `json:"keepAlive" gorm:"default:true;comment:是否缓存"`
	Title          string `json:"title" gorm:"comment:菜单名"`
	TransitionType string `json:"transitionType" gorm:"comment:路由切换动画"`
}
type SysBaseMenuParameter struct {
	global.MPA_MODEL
	SysBaseMenuID uint
	Type          string `json:"type" gorm:"comment:地址栏携带参数为params还是query"` // 地址栏携带参数为params还是query
	Key           string `json:"key" gorm:"comment:地址栏携带参数的key"`            // 地址栏携带参数的key
	Value         string `json:"value" gorm:"comment:地址栏携带参数的值"`            // 地址栏携带参数的值
}

func (SysMenu) TableName() string {
	return "sys_menus"
}
