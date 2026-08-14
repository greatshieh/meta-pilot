package v1

import (
	"server/pkg/dao"

	"gorm.io/gorm"
)

// UserDao 用户DAO结构体
// 实现UserDaoInterface接口
type UserDao struct {
	*dao.BaseDao
}

// NewUserDao 创建用户DAO实例
func NewUserDao(db *gorm.DB) dao.UserDaoInterface {
	return &UserDao{
		BaseDao: dao.NewBaseDao(db),
	}
}

// JwtDao JWT DAO结构体
// 实现JwtDaoInterface接口
type JwtDao struct {
	*dao.BaseDao
}

// NewJwtDao 创建JWT DAO实例
func NewJwtDao(db *gorm.DB) dao.JwtDaoInterface {
	return &JwtDao{
		BaseDao: dao.NewBaseDao(db),
	}
}

// CasbinDao Casbin DAO结构体
// 实现CasbinDaoInterface接口
type CasbinDao struct {
	*dao.BaseDao
}

// NewCasbinDao 创建Casbin DAO实例
func NewCasbinDao(db *gorm.DB) dao.CasbinDaoInterface {
	return &CasbinDao{
		BaseDao: dao.NewBaseDao(db),
	}
}

// AuthorityDao 权限DAO结构体
// 实现AuthorityDaoInterface接口
type AuthorityDao struct {
	*dao.BaseDao
}

// NewAuthorityDao 创建权限DAO实例
func NewAuthorityDao(db *gorm.DB) dao.AuthorityDaoInterface {
	return &AuthorityDao{
		BaseDao: dao.NewBaseDao(db),
	}
}

// MenuDao 菜单DAO结构体
// 实现MenuDaoInterface接口
type MenuDao struct {
	*dao.BaseDao
}

// NewMenuDao 创建菜单DAO实例
func NewMenuDao(db *gorm.DB) dao.MenuDaoInterface {
	return &MenuDao{
		BaseDao: dao.NewBaseDao(db),
	}
}

// ApiDao API DAO结构体
// 实现ApiDaoInterface接口
type ApiDao struct {
	*dao.BaseDao
}

// NewApiDao 创建API DAO实例
func NewApiDao(db *gorm.DB) dao.ApiDaoInterface {
	return &ApiDao{
		BaseDao: dao.NewBaseDao(db),
	}
}
