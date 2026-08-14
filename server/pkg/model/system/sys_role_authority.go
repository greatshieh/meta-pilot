package system

import (
	"time"
)

// SysRoleAuthority 系统角色模型
// 该结构体定义了系统角色的基本信息，包括角色名称、权限等属性
type SysRoleAuthority struct {
	CreatedAt          time.Time           `json:"createdAt"` // 创建时间
	UpdatedAt          time.Time           `json:"updatedAt"` // 更新时间
	DeletedAt          *time.Time          `sql:"index"`
	AuthorityID        uint                `json:"authorityID" gorm:"not null;unique;primaryKey;comment:角色ID;size:90"`      // 角色ID
	AuthorityName      string              `json:"authorityName" gorm:"unique;comment:角色名"`                                 // 角色名
	ParentID           *uint               `json:"parentID" gorm:"comment:父角色ID"`                                           // 父角色ID
	AuthorizedSubRoles []*SysRoleAuthority `json:"authorized_sub_roles" gorm:"many2many:sys_role_data_authority_relation;"` //授权的子角色
	Children           []SysRoleAuthority  `json:"children" gorm:"-"`
	SysMenu            []SysMenu           `json:"menus" gorm:"many2many:sys_menu_authority_relation;joinForeignKey:authority_id;joinReferences:menu_id;"` // 授权的菜单
	Users              []SysUser           `json:"users" gorm:"many2many:sys_user_authority_relation;joinForeignKey:authority_id;joinReferences:user_id;"`                                                         // 授权的用户
	DefaultRouter      string              `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"`                                                    // 默认菜单(默认dashboard)
}

func (SysRoleAuthority) TableName() string {
	return "sys_role_authorities"
}
