// Package system 提供系统服务层功能
// 包含用户、权限、菜单、API、JWT和Casbin等相关的业务逻辑处理
package system

import (
	"server/pkg/dao"
	"server/pkg/db"
	"server/pkg/service"
)

// CasbinService Casbin服务结构体
type CasbinService struct {
	tx           db.Transaction
	apiDao       dao.ApiDaoInterface
	authorityDao dao.AuthorityDaoInterface
	casbinDao    dao.CasbinDaoInterface
}

// NewCasbinService 创建Casbin服务实例
func NewCasbinService(apiDao dao.ApiDaoInterface, casbinDao dao.CasbinDaoInterface, authorityDao dao.AuthorityDaoInterface, tx db.Transaction) service.CasbinServiceInterface {
	return &CasbinService{
		tx:           tx,
		apiDao:       apiDao,
		casbinDao:    casbinDao,
		authorityDao: authorityDao,
	}
}

// AuthorityService 权限服务结构体
type AuthorityService struct {
	tx           db.Transaction
	authorityDao dao.AuthorityDaoInterface
	menuDao      dao.MenuDaoInterface
	casbinDao    dao.CasbinDaoInterface
}

// NewAuthorityService 创建权限服务实例
func NewAuthorityService(authorityDao dao.AuthorityDaoInterface, menuDao dao.MenuDaoInterface, casbinDao dao.CasbinDaoInterface, tx db.Transaction) service.AuthorityServiceInterface {
	return &AuthorityService{
		tx:           tx,
		authorityDao: authorityDao,
		menuDao:      menuDao,
		casbinDao:    casbinDao,
	}
}

// MenuService 菜单服务结构体
type MenuService struct {
	tx           db.Transaction
	menuDao      dao.MenuDaoInterface
	authorityDao dao.AuthorityDaoInterface
}

// NewMenuService 创建菜单服务实例
func NewMenuService(menuDao dao.MenuDaoInterface, authorityDao dao.AuthorityDaoInterface, tx db.Transaction) service.MenuServiceInterface {
	return &MenuService{
		tx:           tx,
		menuDao:      menuDao,
		authorityDao: authorityDao,
	}
}

// ApiService API服务结构体
type ApiService struct {
	tx           db.Transaction
	apiDao       dao.ApiDaoInterface
	casbinDao    dao.CasbinDaoInterface
	authorityDao dao.AuthorityDaoInterface
}

// NewApiService 创建API服务实例
func NewApiService(apiDao dao.ApiDaoInterface, casbinDao dao.CasbinDaoInterface, authorityDao dao.AuthorityDaoInterface, tx db.Transaction) service.ApiServiceInterface {
	return &ApiService{
		tx:           tx,
		apiDao:       apiDao,
		casbinDao:    casbinDao,
		authorityDao: authorityDao,
	}
}

// UserService 用户服务结构体
// 提供用户相关的业务逻辑处理
type UserService struct {
	tx         db.Transaction
	jwtService service.JwtServiceInterface
	userDao    dao.UserDaoInterface
	jwtDao     dao.JwtDaoInterface
}

// NewUserService 创建用户服务实例
func NewUserService(jwtService service.JwtServiceInterface, userDao dao.UserDaoInterface, jwtDao dao.JwtDaoInterface, tx db.Transaction) service.UserServiceInterface {
	return &UserService{
		tx:         tx,
		jwtService: jwtService,
		userDao:    userDao,
		jwtDao:     jwtDao,
	}
}

// JwtService JWT 服务结构体
// 提供 JWT 相关的业务逻辑处理，如 JWT 黑名单管理和令牌生成
type JwtService struct {
	tx      db.Transaction
	jwtDao  dao.JwtDaoInterface
	userDao dao.UserDaoInterface
}

// NewJwtService 创建 JWT 服务实例
func NewJwtService(tx db.Transaction, jwtDao dao.JwtDaoInterface, userDao dao.UserDaoInterface) service.JwtServiceInterface {
	return &JwtService{
		tx:      tx,
		jwtDao:  jwtDao,
		userDao: userDao,
	}
}

// 可以根据需要定义其他版本的Dao
