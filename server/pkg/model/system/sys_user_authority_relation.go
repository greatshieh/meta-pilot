package system

// SysUserAuthorityRelation 是 sysUser 和 sysRoleAuthority 的连接表
type SysUserAuthorityRelation struct {
	UserID      uint `gorm:"column:user_id"`
	AuthorityID uint `gorm:"column:authority_id"`
}

func (s *SysUserAuthorityRelation) TableName() string {
	return "sys_user_authority_relation"
}
